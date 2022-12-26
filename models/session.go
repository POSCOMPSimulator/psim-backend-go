package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createSession = `
INSERT INTO sessions (
  id,
  email,
  refresh_token,
  user_agent,
  client_ip,
  is_blocked,
  expires_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
`

type Session struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIp     string    `json:"client_ip"`
	IsBlocked    bool      `json:"is_blocked"`
	CreateAt     time.Time `json:"created_at"`
	ExpiresAt    time.Time `json:"expires_at"`
}

func (s *Session) CreateSession(db *sql.DB) error {

	if _, err := db.Exec(createSession, s.ID, s.Username, s.RefreshToken,
		s.UserAgent, s.ClientIp, s.IsBlocked, s.ExpiresAt); err != nil {
		return err
	}

	return nil

}

const getSession = `
SELECT id, email, refresh_token, user_agent, client_ip, is_blocked, expires_at, created_at FROM sessions
WHERE id = $1 LIMIT 1
`

func (s *Session) GetSession(db *sql.DB) error {

	if err := db.QueryRow(getSession, s.ID).Scan(&s.ID, &s.Username, &s.RefreshToken,
		&s.UserAgent, &s.ClientIp, &s.IsBlocked, &s.ExpiresAt, &s.CreateAt); err != nil {
		return err
	}

	return nil

}
