package session

type SessionServiceInterface interface {
	StartSession(sessionDto SessionRequest, userID string) (*SessionResponse, error)
	StopSession(sessionDto SessionRequest, userID string) (*SessionResponse, error)
}
