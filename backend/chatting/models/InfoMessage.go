package models

import (
	"errors"
	main_models "main/models"
)

type InfoMessageType string

const (
	Login  InfoMessageType = "login"
	Logout InfoMessageType = "logout"
)

func (t InfoMessageType) IsValid() bool {
	switch t {
	case Login, Logout:
		return true
	}
	return false
}

type InfoMessage struct {
	Content  string          `json:"content" validate:"required"`
	InfoType InfoMessageType `json:"info-type" validate:"required"`
	User     main_models.User
	BaseMessage
}

func (m InfoMessage) validate() error {
	if !m.InfoType.IsValid() {
		return errors.New("invalid info message type")
	}
	return nil
}
