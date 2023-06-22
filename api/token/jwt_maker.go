package token

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const minKeySize = 32

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secreykey string) (Maker, error) {
	if len(secreykey) < minKeySize {
		return nil, errors.New(fmt.Sprintf("error: invalid key size. must be aqt least%v characters", minKeySize))
	}

	return &JWTMaker{secreykey}, nil
}

func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString([]byte(maker.secretKey))
}

func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	panic("implement me")
}
