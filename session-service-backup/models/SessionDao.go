package models

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type PostgresRepository struct {
	DB *gorm.DB
}

func (p *PostgresRepository) CreateSession(sessionDao *Session) (*Session, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var session Session

	// check if active session exists
	err := p.DB.WithContext(ctx).Where("user_id = ? AND end_time IS NULL", sessionDao.UserId).First(&session).Error

	// if yes, return error
	if err == nil {
		return nil, errors.New("user already has an active session — please end it before starting a new one")
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = p.DB.WithContext(ctx).Create(sessionDao).Error

		if err != nil {
			return nil, fmt.Errorf("session creation failed. %w", err)
		}

		return sessionDao, nil
	}

	return nil, fmt.Errorf("error checking for active session: %w", err)
}

func (p *PostgresRepository) StopSession(sessionDao *Session) (*Session, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var session Session

	err := p.DB.WithContext(ctx).Where("user_id = ? AND end_time IS NULL", sessionDao.UserId).First(&session).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("no active session found. %v", err)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to query session: %w", err)
	}

	session.EndTime = sessionDao.EndTime

	result := p.DB.WithContext(ctx).Model(&session).UpdateColumn("end_time", session.EndTime)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to update session: %w", result.Error)
	}
	if result.RowsAffected != 1 {
		return nil, fmt.Errorf("unexpected number of rows affected: %d", result.RowsAffected)
	}

	return &session, nil
}
