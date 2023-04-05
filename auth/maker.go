package auth

import "time"

type Maker interface {
	CreateToken(username string, level int16, verificado bool, duration time.Duration) (string, *Payload, error)

	VerifyToken(token string) (*Payload, error)
}
