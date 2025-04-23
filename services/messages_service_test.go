package services

import (
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"testing-project/domain"
	"testing-project/utils/error_utils"
	"time"
)

var (
	tm                   = time.Now()
	getMessageDomain     func(messageId int64) (*domain.Message, error_utils.MessageErr)
	getAllMessagesDomain func() ([]domain.Message, error_utils.MessageErr)
)

type getDBMock struct{}

func (m *getDBMock) Get(messageId int64) (*domain.Message, error_utils.MessageErr) {
	return getMessageDomain(messageId)
}
func (m *getDBMock) GetAll() ([]domain.Message, error_utils.MessageErr) {
	return getAllMessagesDomain()
}
func (m *getDBMock) Save(*domain.Message) error_utils.MessageErr {
	return nil
}
func (m *getDBMock) Update(*domain.Message) error_utils.MessageErr {
	return nil
}
func (m *getDBMock) Delete(int64) error_utils.MessageErr {
	return nil
}
func (m *getDBMock) Initialize(a, b, c string) *redis.Client {
	return nil
}

// "GetMessage" test cases

func TestMessagesService_GetMessage_Success(t *testing.T) {
	domain.MessageRepo = &getDBMock{}
	getMessageDomain = func(messageId int64) (*domain.Message, error_utils.MessageErr) {
		return &domain.Message{
			Id:        1,
			Title:     "the title",
			Body:      "the body",
			CreatedAt: tm,
		}, nil
	}
	msg, err := MessagesService.GetMessage(1)
	assert.NotNil(t, msg)
	assert.Nil(t, err)
	assert.EqualValues(t, 1, msg.Id)
	assert.EqualValues(t, "the title", msg.Title)
	assert.EqualValues(t, "the body", msg.Body)
	assert.EqualValues(t, tm, msg.CreatedAt)
}

func TestMessagesService_GetMessage_NotFoundID(t *testing.T) {
	domain.MessageRepo = &getDBMock{}
	getMessageDomain = func(messageId int64) (*domain.Message, error_utils.MessageErr) {
		return nil, error_utils.NewNotFoundError("the id is not found")
	}
	msg, err := MessagesService.GetMessage(1)
	assert.Nil(t, msg)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.Status())
	assert.EqualValues(t, "the id is not found", err.Message())
	assert.EqualValues(t, "not_found", err.Error())
}

// "GetAllMessages" test cases

func TestMessagesService_GetAllMessages(t *testing.T) {
	domain.MessageRepo = &getDBMock{}
	getAllMessagesDomain = func() ([]domain.Message, error_utils.MessageErr) {
		return []domain.Message{
			{Id: 1, Title: "first title", Body: "first body"},
			{Id: 2, Title: "second title", Body: "second body"},
		}, nil
	}
	messages, err := MessagesService.GetAllMessages()
	assert.Nil(t, err)
	assert.NotNil(t, messages)
	assert.EqualValues(t, 2, len(messages))
	assert.EqualValues(t, 1, messages[0].Id)
	assert.EqualValues(t, "first title", messages[0].Title)
	assert.EqualValues(t, "first body", messages[0].Body)
	assert.EqualValues(t, 2, messages[1].Id)
	assert.EqualValues(t, "second title", messages[1].Title)
	assert.EqualValues(t, "second body", messages[1].Body)
}

func TestMessagesService_GetAllMessages_Error_Getting_Messages(t *testing.T) {
	domain.MessageRepo = &getDBMock{}
	getAllMessagesDomain = func() ([]domain.Message, error_utils.MessageErr) {
		return nil, error_utils.NewInternalServerError("error getting messages")
	}
	messages, err := MessagesService.GetAllMessages()
	assert.NotNil(t, err)
	assert.Nil(t, messages)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status())
	assert.EqualValues(t, "error getting messages", err.Message())
	assert.EqualValues(t, "server_error", err.Error())
}
