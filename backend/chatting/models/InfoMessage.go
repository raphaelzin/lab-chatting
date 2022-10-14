package models

import (
	"encoding/json"
	"errors"
	"main/models"
	"time"

	"github.com/google/uuid"
)

type InfoMessageType string

const (
	Login   InfoMessageType = "login"
	Logout  InfoMessageType = "logout"
	Welcome InfoMessageType = "welcome"
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
	BaseMessage
}

func (m InfoMessage) validate() error {
	if !m.InfoType.IsValid() {
		return errors.New("invalid info message type")
	}
	return nil
}

func newBaseInfoMessage() InfoMessage {
	message := new(InfoMessage)
	message.Type = Info
	message.Id = uuid.NewString()
	message.Timestamp = time.Now().UTC().Format(time.RFC3339)

	return *message
}

func NewLoginMessage(user models.User) InfoMessage {
	message := newBaseInfoMessage()
	message.InfoType = Login
	message.Content = user.Username + " joined the party!"
	return message
}

func NewLogoutMessage(user models.User) InfoMessage {
	message := newBaseInfoMessage()
	message.InfoType = Login
	message.Content = user.Username + " is a quitter, he just left!"

	return message
}

func NewWelcomeMessage(token string) InfoMessage {
	message := newBaseInfoMessage()
	message.InfoType = Welcome
	message.Content = token
	return message
}

func (m InfoMessage) AsData() []byte {
	data, _ := json.Marshal(m)
	return data
}
