package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"vstats/internal/cloud/config"

	"github.com/redis/go-redis/v9"
)

var client *redis.Client

// Key prefixes
const (
	PrefixSession    = "vstats:session:"
	PrefixUserCache  = "vstats:user:"
	PrefixServerLive = "vstats:server:"
	PrefixRateLimit  = "vstats:ratelimit:"
	PrefixOAuthState = "vstats:oauth:"
	PrefixWSConn     = "vstats:ws:"
)

// Connect establishes connection to Redis
func Connect(cfg *config.Config) (*redis.Client, error) {
	opt, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	client = redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return client, nil
}

// GetClient returns the Redis client
func GetClient() *redis.Client {
	return client
}

// Close closes the Redis connection
func Close() error {
	if client != nil {
		return client.Close()
	}
	return nil
}

// HealthCheck checks Redis health
func HealthCheck(ctx context.Context) error {
	if client == nil {
		return fmt.Errorf("redis client not initialized")
	}
	return client.Ping(ctx).Err()
}

// ============================================================================
// Session Management
// ============================================================================

type SessionData struct {
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	Email     string    `json:"email,omitempty"`
	Plan      string    `json:"plan"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// SetSession stores session in Redis
func SetSession(ctx context.Context, sessionID string, data *SessionData, expiry time.Duration) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return client.Set(ctx, PrefixSession+sessionID, jsonData, expiry).Err()
}

// GetSession retrieves session from Redis
func GetSession(ctx context.Context, sessionID string) (*SessionData, error) {
	val, err := client.Get(ctx, PrefixSession+sessionID).Result()
	if err != nil {
		return nil, err
	}

	var data SessionData
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return nil, err
	}

	return &data, nil
}

// DeleteSession removes session from Redis
func DeleteSession(ctx context.Context, sessionID string) error {
	return client.Del(ctx, PrefixSession+sessionID).Err()
}

// ============================================================================
// OAuth State Management
// ============================================================================

type OAuthStateData struct {
	Provider    string `json:"provider"`
	RedirectURL string `json:"redirect_url"`
	CreatedAt   int64  `json:"created_at"`
}

// SetOAuthState stores OAuth state
func SetOAuthState(ctx context.Context, state string, data *OAuthStateData) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	// OAuth state expires in 10 minutes
	return client.Set(ctx, PrefixOAuthState+state, jsonData, 10*time.Minute).Err()
}

// GetOAuthState retrieves OAuth state
func GetOAuthState(ctx context.Context, state string) (*OAuthStateData, error) {
	val, err := client.Get(ctx, PrefixOAuthState+state).Result()
	if err != nil {
		return nil, err
	}

	var data OAuthStateData
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return nil, err
	}

	// Delete after retrieval (one-time use)
	client.Del(ctx, PrefixOAuthState+state)

	return &data, nil
}

// ============================================================================
// Rate Limiting
// ============================================================================

// CheckRateLimit checks if request is within rate limit
// Returns remaining requests and error
func CheckRateLimit(ctx context.Context, key string, limit int, window time.Duration) (int, error) {
	fullKey := PrefixRateLimit + key

	pipe := client.Pipeline()
	incr := pipe.Incr(ctx, fullKey)
	pipe.Expire(ctx, fullKey, window)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return 0, err
	}

	count := int(incr.Val())
	remaining := limit - count
	if remaining < 0 {
		remaining = 0
	}

	return remaining, nil
}

// ============================================================================
// Server Live Status (for real-time dashboard)
// ============================================================================

type ServerLiveData struct {
	ServerID    string          `json:"server_id"`
	Status      string          `json:"status"`
	LastSeenAt  time.Time       `json:"last_seen_at"`
	Metrics     json.RawMessage `json:"metrics,omitempty"`
	AgentConnID string          `json:"agent_conn_id,omitempty"`
}

// SetServerLive updates server live status
func SetServerLive(ctx context.Context, serverID string, data *ServerLiveData) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	// Server live status expires in 60 seconds (should be refreshed by agent)
	return client.Set(ctx, PrefixServerLive+serverID+":live", jsonData, 60*time.Second).Err()
}

// GetServerLive retrieves server live status
func GetServerLive(ctx context.Context, serverID string) (*ServerLiveData, error) {
	val, err := client.Get(ctx, PrefixServerLive+serverID+":live").Result()
	if err != nil {
		return nil, err
	}

	var data ServerLiveData
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return nil, err
	}

	return &data, nil
}

// DeleteServerLive removes server live status
func DeleteServerLive(ctx context.Context, serverID string) error {
	return client.Del(ctx, PrefixServerLive+serverID+":live").Err()
}

// GetAllLiveServers retrieves all live server statuses
func GetAllLiveServers(ctx context.Context) (map[string]*ServerLiveData, error) {
	keys, err := client.Keys(ctx, PrefixServerLive+"*:live").Result()
	if err != nil {
		return nil, err
	}

	result := make(map[string]*ServerLiveData)
	for _, key := range keys {
		val, err := client.Get(ctx, key).Result()
		if err != nil {
			continue
		}

		var data ServerLiveData
		if err := json.Unmarshal([]byte(val), &data); err != nil {
			continue
		}
		result[data.ServerID] = &data
	}

	return result, nil
}

// ============================================================================
// WebSocket Connection Tracking
// ============================================================================

// AddWSConnection registers a WebSocket connection
func AddWSConnection(ctx context.Context, connType, connID, userID string) error {
	key := PrefixWSConn + connType + ":" + connID
	return client.Set(ctx, key, userID, 24*time.Hour).Err()
}

// RemoveWSConnection removes a WebSocket connection
func RemoveWSConnection(ctx context.Context, connType, connID string) error {
	key := PrefixWSConn + connType + ":" + connID
	return client.Del(ctx, key).Err()
}

// CountWSConnections counts active WebSocket connections
func CountWSConnections(ctx context.Context, connType string) (int64, error) {
	keys, err := client.Keys(ctx, PrefixWSConn+connType+":*").Result()
	if err != nil {
		return 0, err
	}
	return int64(len(keys)), nil
}

// ============================================================================
// Release Cache (for GitHub releases)
// ============================================================================

const (
	PrefixRelease = "vstats:release:"
	PrefixBinary  = "vstats:binary:"
)

// ReleaseInfo represents cached release info
type ReleaseInfo struct {
	Version     string            `json:"version"`
	Name        string            `json:"name"`
	Body        string            `json:"body"`
	PublishedAt string            `json:"published_at"`
	Assets      map[string]string `json:"assets"`
	CachedAt    int64             `json:"cached_at"`
}

// SetReleaseInfo caches release info
func SetReleaseInfo(ctx context.Context, info *ReleaseInfo, expiry time.Duration) error {
	jsonData, err := json.Marshal(info)
	if err != nil {
		return err
	}
	return client.Set(ctx, PrefixRelease+"latest", jsonData, expiry).Err()
}

// GetReleaseInfo retrieves cached release info
func GetReleaseInfo(ctx context.Context) (*ReleaseInfo, error) {
	val, err := client.Get(ctx, PrefixRelease+"latest").Result()
	if err != nil {
		return nil, err
	}

	var info ReleaseInfo
	if err := json.Unmarshal([]byte(val), &info); err != nil {
		return nil, err
	}

	return &info, nil
}

// DeleteReleaseInfo removes cached release info
func DeleteReleaseInfo(ctx context.Context) error {
	return client.Del(ctx, PrefixRelease+"latest").Err()
}

// BinaryCache represents cached binary data
type BinaryCache struct {
	Data        []byte `json:"data"`
	ContentType string `json:"content_type"`
	Size        int64  `json:"size"`
	CachedAt    int64  `json:"cached_at"`
}

// SetBinaryCache caches binary data
func SetBinaryCache(ctx context.Context, version, name string, data []byte, contentType string, expiry time.Duration) error {
	// Don't cache files larger than 100MB
	if len(data) > 100*1024*1024 {
		return fmt.Errorf("binary too large to cache: %d bytes", len(data))
	}

	cache := &BinaryCache{
		Data:        data,
		ContentType: contentType,
		Size:        int64(len(data)),
		CachedAt:    time.Now().Unix(),
	}

	jsonData, err := json.Marshal(cache)
	if err != nil {
		return err
	}

	key := PrefixBinary + version + ":" + name
	return client.Set(ctx, key, jsonData, expiry).Err()
}

// GetBinaryCache retrieves cached binary data
func GetBinaryCache(ctx context.Context, version, name string) ([]byte, string, error) {
	key := PrefixBinary + version + ":" + name
	val, err := client.Get(ctx, key).Result()
	if err != nil {
		return nil, "", err
	}

	var cache BinaryCache
	if err := json.Unmarshal([]byte(val), &cache); err != nil {
		return nil, "", err
	}

	return cache.Data, cache.ContentType, nil
}

// DeleteBinaryCache removes cached binary
func DeleteBinaryCache(ctx context.Context, version, name string) error {
	key := PrefixBinary + version + ":" + name
	return client.Del(ctx, key).Err()
}

// ClearAllBinaryCache removes all cached binaries
func ClearAllBinaryCache(ctx context.Context) error {
	keys, err := client.Keys(ctx, PrefixBinary+"*").Result()
	if err != nil {
		return err
	}
	if len(keys) > 0 {
		return client.Del(ctx, keys...).Err()
	}
	return nil
}
