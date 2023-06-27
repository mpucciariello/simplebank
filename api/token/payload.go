package token

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

var ErrExpiredToken = errors.New("token is expired")
var ErrInvalidToken = errors.New("token is invalid")

type Payload struct {
	ID        uuid.UUID `json:"id"`
	UserName  string    `json:"user_name"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, errors.New("error generating token id")
	}

	payload := Payload{
		ID:        id,
		UserName:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return &payload, nil
}

func (p Payload) Valid() error {
	if time.Now().After(p.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
