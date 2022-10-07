package models

import (
	main_models "main/models"
)

type TextMessage struct {
	Content string `json:"content" validate:"required"`
	User    main_models.User
	BaseMessage
}

func (m TextMessage) validate() error {
	return nil
}
