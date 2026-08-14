package sid

import "github.com/google/uuid"

type Sid struct{}

func NewSid() *Sid {
	return &Sid{}
}

func (s Sid) GenString() (string, error) {
	return uuid.New().String(), nil
}
