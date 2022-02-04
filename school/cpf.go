package school

import (
	"errors"
	"regexp"
)

type Cpf struct {
	number string
}

var ErrCpfNotValid = errors.New("Invalid CPF")

func NewCPF(number string) (*Cpf, error) {
	validCPF := regexp.MustCompile(`^\d{3}\.\d{3}\.\d{3}\-\d{2}$`)

	if number == "" || !validCPF.MatchString(number) {
		return nil, ErrCpfNotValid
	}
	return &Cpf{
		number: number,
	}, nil
}

func (c *Cpf) GetNumero() string {
	return c.number
}
