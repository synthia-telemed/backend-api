package id

import "github.com/jaevor/go-nanoid"

type Generator interface {
	GenerateRoomID() (string, error)
}

type NanoID struct {
}

func NewNanoID() Generator {
	return &NanoID{}
}

func (n NanoID) GenerateRoomID() (string, error) {
	generator, err := nanoid.Standard(21)
	if err != nil {
		return "", err
	}
	return generator(), nil
}
