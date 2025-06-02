package session

import (
	"fmt"
	"productivity-planner/task-service/models"
	"time"

	"github.com/google/uuid"
)

type SessionService struct {
	Repo models.Repository
}

func (s *SessionService) StartSession(sessionDto SessionRequest, userID string) (*SessionResponse, error) {

	UUID, err := uuid.Parse(userID)

	if err != nil {
		return nil, fmt.Errorf("uuid parse error: %v", err)
	}

	session := &models.Session{
		ID:          uuid.New(),
		UserId:      UUID,
		SessionType: string(sessionDto.SessionType),
		StartTime:   time.Now().UTC(),
		EndTime:     nil,
	}

	resp, err := s.Repo.CreateSession(session)

	if err != nil {
		return nil, fmt.Errorf("session creation error: %w", err)
	}

	sessionResponse := &SessionResponse{
		Status: STARTED,
		Session: Session{
			SessionId:   resp.ID.String(),
			SessionType: resp.SessionType,
			StartTime:   session.StartTime.String(),
		},
	}

	return sessionResponse, nil
}

func (s *SessionService) StopSession(sessionDto SessionRequest, userID string) (*SessionResponse, error) {

	UUID, err := uuid.Parse(userID)

	if err != nil {
		return nil, fmt.Errorf("uuid parse error: %v", err)
	}

	endTime := time.Now().UTC()
	endTimePtr := &endTime

	session := &models.Session{
		UserId:      UUID,
		SessionType: string(sessionDto.SessionType),
		EndTime:     endTimePtr,
	}

	resp, err := s.Repo.StopSession(session)

	if err != nil {
		return nil, fmt.Errorf("error while ending session: %w", err)
	}

	sessionResponse := &SessionResponse{
		Status: ENDED,
		Session: Session{
			SessionId:   resp.ID.String(),
			SessionType: resp.SessionType,
			StartTime:   resp.StartTime.String(),
			EndTime:     resp.EndTime.String(),
		},
	}

	return sessionResponse, nil
}
