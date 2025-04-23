package services

import (
	"testing-project/domain"
	"testing-project/utils/error_utils"
)

var (
	MessagesService messageServiceInterface = &messagesService{}
)

type messagesService struct{}

type messageServiceInterface interface {
	GetMessage(int64) (*domain.Message, error_utils.MessageErr)
	GetAllMessages() ([]domain.Message, error_utils.MessageErr)
}

func (m *messagesService) GetMessage(msgId int64) (*domain.Message, error_utils.MessageErr) {
	message, err := domain.MessageRepo.Get(msgId)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (m *messagesService) GetAllMessages() ([]domain.Message, error_utils.MessageErr) {
	messages, err := domain.MessageRepo.GetAll()
	if err != nil {
		return nil, err
	}
	return messages, nil
}
