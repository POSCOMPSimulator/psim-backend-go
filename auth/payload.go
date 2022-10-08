package auth

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has experied")
)

type Payload struct {
	ID         uuid.UUID `json:"id"`
	Username   string    `json:"username"`
	IssuedAt   time.Time `json:"issued_at"`
	ExperiesAt time.Time `json:"experies_at"`
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:         tokenId,
		Username:   username,
		IssuedAt:   time.Now(),
		ExperiesAt: time.Now().Add(duration),
	}

	return payload, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExperiesAt) {
		return ErrExpiredToken
	}
	return nil
}
