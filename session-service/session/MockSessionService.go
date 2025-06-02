package session

import (
	"productivity-planner/task-service/models"
	"time"
)

type MockSessionService struct {
	Repo models.Repository
}

func (s *MockSessionService) StartSession(sessionDto SessionRequest, userID string) (*SessionResponse, error) {
	return &SessionResponse{
		Status: STARTED,
		Session: Session{
			SessionId:   "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			SessionType: string(sessionDto.SessionType),
			StartTime:   time.Now().Add(-1 * time.Hour).UTC().String(),
		},
	}, nil
}

func (s *MockSessionService) StopSession(sessionDto SessionRequest, userID string) (*SessionResponse, error) {

	return &SessionResponse{
		Status: ENDED,
		Session: Session{
			SessionId:   "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			SessionType: string(sessionDto.SessionType),
			StartTime:   time.Now().Add(-1 * time.Hour).UTC().String(),
			EndTime:     time.Now().UTC().String(),
		},
	}, nil
}
