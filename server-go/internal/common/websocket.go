package common

// ============================================================================
// WebSocket Message Types
// ============================================================================

type AuthMessage struct {
	Type     string `json:"type"`
	ServerID string `json:"server_id"`
	Token    string `json:"token"`
	Version  string `json:"version"`
}

type MetricsMessage struct {
	Type    string        `json:"type"`
	Metrics SystemMetrics `json:"metrics"`
}

type ServerResponse struct {
	Type        string             `json:"type"`
	Status      string             `json:"status,omitempty"`
	Message     string             `json:"message,omitempty"`
	Command     string             `json:"command,omitempty"`
	DownloadURL string             `json:"download_url,omitempty"`
	Force       bool               `json:"force,omitempty"`
	PingTargets []PingTargetConfig `json:"ping_targets,omitempty"`
}

// ============================================================================
// Registration Types
// ============================================================================

type RegisterRequest struct {
	Name     string `json:"name"`
	Location string `json:"location"`
	Provider string `json:"provider"`
}

type RegisterResponse struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

