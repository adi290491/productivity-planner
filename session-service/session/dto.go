package session

type SessionRequest struct {
	SessionType SessionType `json:"type"`
}

type SessionResponse struct {
	Status  SessionStatus `json:"status"`
	Session Session       `json:"session"`
}

type Session struct {
	SessionId   string `json:"sessionId"`
	SessionType string `json:"type"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
}

type SessionStatus string

const (
	STARTED SessionStatus = "started"
	ENDED  SessionStatus = "ended"
)

type SessionType string

const (
	FOCUS   SessionType = "focus"
	MEETING SessionType = "meeting"
	BREAK   SessionType = "break"
)

func (t SessionType) IsValid() bool {
	switch t {
	case FOCUS, MEETING, BREAK:
		return true
	default:
		return false
	}
}
