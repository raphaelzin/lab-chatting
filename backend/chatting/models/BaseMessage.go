package models

import (
	"encoding/json"
	"errors"

	"github.com/go-playground/validator"
)

type MessageType string

const (
	Text MessageType = "text"
	Info MessageType = "info"
)

func (t MessageType) IsValid() bool {
	switch t {
	case Text:
		return true
	}
	return false
}

type BaseMessage struct {
	Type MessageType `json:"type" validate:"required"`
}

func (m BaseMessage) validate() error {
	if !m.Type.IsValid() {
		return errors.New("invalid message type")
	}
	return nil
}

type Message interface {
	validate() error
}

func Parse[T Message](data []byte) (message T, err error) {
	err = json.Unmarshal(data, &message)
	validate := validator.New()

	if error := validate.Struct(message); error != nil {
		return message, error
	}

	if error := message.validate(); error != nil {
		return message, error
	}

	return message, err
}
