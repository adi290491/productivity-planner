package models

import (
	"errors"
	"fmt"
)

type TestDBRepo struct {
	ActiveSessions map[string]*Session
}

func (r *TestDBRepo) CreateSession(sessionDao *Session) (*Session, error) {
	// Simulate: if user already has an active session, return error
	if existing, ok := r.ActiveSessions[sessionDao.UserId.String()]; ok && existing.EndTime == nil {
		return nil, errors.New("user already has an active session — please end it before starting a new one")
	}

	// Simulate: create a new session
	r.ActiveSessions[sessionDao.UserId.String()] = sessionDao
	return sessionDao, nil
}

func (r *TestDBRepo) StopSession(sessionDao *Session) (*Session, error) {
	// Simulate: no active session found
	existing, ok := r.ActiveSessions[sessionDao.UserId.String()]
	if !ok || existing.EndTime != nil {
		return nil, fmt.Errorf("no active session found")
	}

	// Simulate: stop session
	existing.EndTime = sessionDao.EndTime
	return existing, nil
}
