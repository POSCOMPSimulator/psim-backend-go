package auth

import "time"

type Maker interface {
	CreateToken(username string, level int16, duration time.Duration) (string, error)

	VerifyToken(token string) (*Payload, error)
}
