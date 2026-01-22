package models

type Repository interface {
	CreateSession(session *Session) (*Session, error)
	StopSession(session *Session) (*Session, error)
}
