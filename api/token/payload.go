package token

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

type Payload struct {
	ID        uuid.UUID `json:"id"`
	UserName  string    `json:"user_name"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func (p Payload) GetExpirationTime() (*jwt.NumericDate, error) {
	//TODO implement me
	panic("implement me")
}

func (p Payload) GetIssuedAt() (*jwt.NumericDate, error) {
	//TODO implement me
	panic("implement me")
}

func (p Payload) GetNotBefore() (*jwt.NumericDate, error) {
	//TODO implement me
	panic("implement me")
}

func (p Payload) GetIssuer() (string, error) {
	//TODO implement me
	panic("implement me")
}

func (p Payload) GetSubject() (string, error) {
	//TODO implement me
	panic("implement me")
}

func (p Payload) GetAudience() (jwt.ClaimStrings, error) {
	//TODO implement me
	panic("implement me")
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
