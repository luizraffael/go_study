package school

import (
	"errors"
	"regexp"
)

type Email struct {
	address string
}

func NewEmail(email string) (*Email, error) {
	if !regexp.MustCompile(`^[a-zA-Z0-9._]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).MatchString(email) {
		return nil, errors.New("Email invalido")
	}
	return &Email{address: email}, nil
}

func (e *Email) GetAddress() string {
	return e.address
}
