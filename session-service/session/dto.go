package session

type SessionRequest struct {
	SessionType SessionType `json:"session_type"`
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
	ENDED   SessionStatus = "ended"
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
/*
You are a professional devops guide. Given the github link of a project, scan through it and help me with
certain devops related issues I am facing. The github repo contains multiple project folders as follows

- productivity-planner (root folder) (contains all project folders)
	- gateway (Go)
	- session-service (Go)
	- summary-service (Go)
	- trend-analysis-worker (Go) (schedulers)
	- trend-service (Go)
	- user-service (Go)
	- productivity-frontend (this is in React + Tailwind)
	- docker-compose
	- init
		- init.sql

Each Go project folder contains the same folder structure
- project-folder
	- cmd
	- config
	- models
	- <dto>
	- utils
	- Dockerfile

I will specify the project and you help me out with them
*/