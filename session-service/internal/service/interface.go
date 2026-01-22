package service

// SessionService defines the interface for session operations
type SessionServiceInterface interface {
	StartSession(req StartSessionRequest, userID string) (*SessionResponse, error)
	StopSession(req StopSessionRequest, userID string) (*SessionResponse, error)
}
