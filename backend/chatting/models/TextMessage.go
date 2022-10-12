package models

import (
	main_models "main/models"
	"time"

	"github.com/google/uuid"
)

type TextMessage struct {
	Content string           `json:"content" validate:"required"`
	User    main_models.User `json:"user"`
	BaseMessage
}

func (m TextMessage) validate() error {
	return nil
}

func NewTextMessage(m RawMessage) *TextMessage {
	textMessage := new(TextMessage)
	textMessage.Type = m.Type
	textMessage.Content = m.Content
	textMessage.Id = uuid.NewString()
	textMessage.Timestamp = time.Now().UTC().Format(time.RFC3339)

	return textMessage
}
