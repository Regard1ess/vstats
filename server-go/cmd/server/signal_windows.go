//go:build windows
// +build windows

package main

// SignalError represents different types of signal errors
type SignalError struct {
	Type    string // "not_found", "permission_denied", "other"
	Message string
	PID     int
}

func (e *SignalError) Error() string {
	return e.Message
}

// SetupSignalHandler is a no-op on Windows
// Windows doesn't support SIGHUP
func SetupSignalHandler(state *AppState) {
	// No-op on Windows
}

// findAndSignalServer is not supported on Windows
func findAndSignalServer() error {
	return &SignalError{
		Type:    "other",
		Message: "auto-reload is not supported on Windows",
	}
}
