package auth

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type PasetoMaker struct {
	paseto        *paseto.V2
	symmectricKey []byte
}

func NewPasetoMaker(symmectricKey string) (Maker, error) {
	if len(symmectricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must have at least %d characters", chacha20poly1305.KeySize)
	}

	maker := &PasetoMaker{
		paseto:        paseto.NewV2(),
		symmectricKey: []byte(symmectricKey),
	}
	return maker, nil
}

func (maker *PasetoMaker) CreateToken(username string, level int16, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, level, duration)
	if err != nil {
		return "", err
	}

	return maker.paseto.Encrypt(maker.symmectricKey, payload, nil)
}

func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := maker.paseto.Decrypt(token, maker.symmectricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}
