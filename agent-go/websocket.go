package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	InitialReconnectDelay = 5 * time.Second
	MaxReconnectDelay     = 60 * time.Second
	AuthTimeout           = 10 * time.Second
	PingInterval          = 30 * time.Second
)

type WebSocketClient struct {
	config    *AgentConfig
	collector *MetricsCollector
}

func NewWebSocketClient(config *AgentConfig) *WebSocketClient {
	return &WebSocketClient{
		config:    config,
		collector: NewMetricsCollector(),
	}
}

func (wsc *WebSocketClient) Run() {
	reconnectDelay := InitialReconnectDelay

	for {
		log.Printf("Connecting to %s...", wsc.config.WSUrl())

		if err := wsc.connectAndRun(); err != nil {
			log.Printf("Connection error: %v", err)
		} else {
			log.Println("Connection closed normally")
			reconnectDelay = InitialReconnectDelay
		}

		log.Printf("Reconnecting in %v...", reconnectDelay)
		time.Sleep(reconnectDelay)

		// Exponential backoff
		reconnectDelay *= 2
		if reconnectDelay > MaxReconnectDelay {
			reconnectDelay = MaxReconnectDelay
		}
	}
}

func (wsc *WebSocketClient) connectAndRun() error {
	wsURL := wsc.config.WSUrl()

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer conn.Close()

	log.Println("Connected to WebSocket server")

	// Send authentication message
	authMsg := AuthMessage{
		Type:     "auth",
		ServerID: wsc.config.ServerID,
		Token:    wsc.config.AgentToken,
		Version:  AgentVersion,
	}

	authData, err := json.Marshal(authMsg)
	if err != nil {
		return fmt.Errorf("failed to serialize auth message: %w", err)
	}

	if err := conn.WriteMessage(websocket.TextMessage, authData); err != nil {
		return fmt.Errorf("failed to send auth message: %w", err)
	}

	log.Println("Sent authentication message")

	// Wait for auth response
	conn.SetReadDeadline(time.Now().Add(AuthTimeout))
	_, message, err := conn.ReadMessage()
	if err != nil {
		return fmt.Errorf("failed to receive auth response: %w", err)
	}

	var response ServerResponse
	if err := json.Unmarshal(message, &response); err != nil {
		return fmt.Errorf("failed to parse auth response: %w", err)
	}

	if response.Status != "ok" {
		return fmt.Errorf("authentication failed: %s", response.Message)
	}

	// Update ping targets from server config if provided
	if len(response.PingTargets) > 0 {
		log.Printf("Received %d ping targets from server", len(response.PingTargets))
		wsc.collector.SetPingTargets(response.PingTargets)
	}

	log.Println("Authentication successful!")

	// Reset read deadline
	conn.SetReadDeadline(time.Time{})

	// Start metrics sending loop
	metricsTicker := time.NewTicker(time.Duration(wsc.config.IntervalSecs) * time.Second)
	defer metricsTicker.Stop()

	pingTicker := time.NewTicker(PingInterval)
	defer pingTicker.Stop()

	// Handle incoming messages
	done := make(chan error, 1)

	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				done <- err
				return
			}

			var response ServerResponse
			if err := json.Unmarshal(message, &response); err != nil {
				continue
			}

			switch response.Type {
			case "error":
				log.Printf("Server error: %s", response.Message)
			case "command":
				if response.Command == "update" {
					if response.Force {
						log.Println("Received FORCE update command from server")
					} else {
						log.Println("Received update command from server")
					}
					wsc.handleUpdateCommand(response.DownloadURL, response.Force)
				}
			case "config":
				// Handle runtime config update (e.g., ping targets)
				if len(response.PingTargets) > 0 {
					log.Printf("Received updated ping targets from server: %d targets", len(response.PingTargets))
					wsc.collector.SetPingTargets(response.PingTargets)
				} else {
					log.Println("Received config update: clearing ping targets")
					wsc.collector.SetPingTargets(nil)
				}
			}
		}
	}()

	for {
		select {
		case <-metricsTicker.C:
			metrics := wsc.collector.Collect()
			msg := MetricsMessage{
				Type:    "metrics",
				Metrics: metrics,
			}

			data, err := json.Marshal(msg)
			if err != nil {
				log.Printf("Failed to serialize metrics: %v", err)
				continue
			}

			if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
				return fmt.Errorf("failed to send metrics: %w", err)
			}

		case <-pingTicker.C:
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return fmt.Errorf("failed to send ping: %w", err)
			}

		case err := <-done:
			return err
		}
	}
}

func (wsc *WebSocketClient) handleUpdateCommand(downloadURL string, force bool) {
	if force {
		log.Println("Starting FORCE self-update process (will update regardless of version)...")
	} else {
		log.Println("Starting self-update process...")
	}

	// Get the current executable path
	currentExe, err := os.Executable()
	if err != nil {
		log.Printf("Failed to get current executable path: %v", err)
		return
	}

	// Determine download URL and check version
	url := downloadURL
	var latestVersion string
	if url == "" {
		// Build GitHub Releases URL based on OS and architecture
		osName := runtime.GOOS
		arch := runtime.GOARCH
		
		// Map Go architecture names to release naming
		if arch == "amd64" {
			arch = "amd64"
		} else if arch == "arm64" {
			arch = "arm64"
		} else if arch == "386" {
			arch = "386"
		}
		
		// Determine binary name
		binaryName := fmt.Sprintf("vstats-agent-%s-%s", osName, arch)
		if osName == "windows" {
			binaryName += ".exe"
		}
		
		// Try to get latest version from GitHub API
		latestVersion = "latest"
		if latest, err := fetchLatestGitHubVersion("zsai001", "vstats"); err == nil && latest != nil {
			latestVersion = *latest
			
			// Skip update if already on latest version (unless force is true)
			// Compare versions without 'v' prefix
			latestVersionClean := strings.TrimPrefix(latestVersion, "v")
			currentVersionClean := strings.TrimPrefix(AgentVersion, "v")
			if !force && latestVersionClean == currentVersionClean {
				log.Printf("Already on latest version %s, skipping update", AgentVersion)
				return
			}
			log.Printf("Update available: current=%s, latest=%s", AgentVersion, latestVersion)
		}
		
		// Build GitHub Releases download URL
		url = fmt.Sprintf("https://github.com/zsai001/vstats/releases/download/%s/%s", latestVersion, binaryName)
		log.Printf("No download URL provided, using GitHub Releases: %s", url)
	} else {
		log.Printf("Using provided download URL: %s", url)
	}
	
	if force {
		log.Printf("Force update enabled, current version: %s", AgentVersion)
	}

	log.Printf("Downloading update from: %s", url)

	// Download to a temporary file
	tempPath := currentExe + ".new"

	if err := downloadFile(url, tempPath); err != nil {
		log.Printf("Failed to download update: %v", err)
		return
	}

	log.Println("Download complete, applying update...")

	// On Unix, set execute permissions
	if runtime.GOOS != "windows" {
		if err := os.Chmod(tempPath, 0755); err != nil {
			log.Printf("Failed to set permissions: %v", err)
			os.Remove(tempPath)
			return
		}
	}

	// Backup current executable
	backupPath := currentExe + ".backup"
	if err := os.Rename(currentExe, backupPath); err != nil {
		log.Printf("Failed to backup current executable: %v", err)
		os.Remove(tempPath)
		return
	}

	// Move new executable to current path
	if err := os.Rename(tempPath, currentExe); err != nil {
		log.Printf("Failed to install new executable: %v", err)
		// Try to restore backup
		os.Rename(backupPath, currentExe)
		return
	}

	// Remove backup
	os.Remove(backupPath)

	log.Println("Update installed successfully! Restarting...")

	// Restart the agent using systemd-run to avoid being killed by cgroup
	if runtime.GOOS == "linux" {
		// Use systemd-run --no-block to run restart in an independent transient unit
		// This prevents the restart command from being killed when vstats-agent stops
		cmd := exec.Command("systemd-run", "--no-block", "systemctl", "restart", "vstats-agent")
		if err := cmd.Start(); err != nil {
			log.Printf("Failed to schedule restart via systemd-run: %v", err)
			// Fallback to direct systemctl (may not work in all cases)
			exec.Command("systemctl", "restart", "vstats-agent").Start()
		} else {
			log.Println("Restart scheduled via systemd-run")
		}
	} else if runtime.GOOS == "windows" {
		// On Windows, use sc.exe to restart the service
		cmd := exec.Command("cmd", "/C", "sc", "stop", "vstats-agent", "&&", "timeout", "/t", "2", "&&", "sc", "start", "vstats-agent")
		cmd.Start()
	}

	// Give systemd-run a moment to register the restart command
	time.Sleep(500 * time.Millisecond)

	// Exit to allow restart
	os.Exit(0)
}

// downloadFile downloads a file from URL to path
func downloadFile(url, path string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status: %d", resp.StatusCode)
	}

	out, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		os.Remove(path)
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// fetchLatestGitHubVersion fetches the latest release version from GitHub
func fetchLatestGitHubVersion(owner, repo string) (*string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo)

	client := &http.Client{Timeout: 10 * time.Second}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "vstats-agent")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status: %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	tagName, ok := result["tag_name"].(string)
	if !ok {
		return nil, fmt.Errorf("no tag_name in response")
	}

	// Keep the original tag name (with 'v' prefix) for download URL
	return &tagName, nil
}
