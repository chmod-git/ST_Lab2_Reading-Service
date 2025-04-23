package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing-project/domain"
	"testing-project/services"
	"testing-project/utils/error_utils"
)

var (
	getMessageService    func(msgId int64) (*domain.Message, error_utils.MessageErr)
	getAllMessageService func() ([]domain.Message, error_utils.MessageErr)
)

type serviceMock struct{}

func (sm *serviceMock) GetMessage(msgId int64) (*domain.Message, error_utils.MessageErr) {
	return getMessageService(msgId)
}

func (sm *serviceMock) GetAllMessages() ([]domain.Message, error_utils.MessageErr) {
	return getAllMessageService()
}

// "GetMessage" test cases

func TestGetMessage_Success(t *testing.T) {
	services.MessagesService = &serviceMock{}
	getMessageService = func(msgId int64) (*domain.Message, error_utils.MessageErr) {
		return &domain.Message{
			Id:    1,
			Title: "the title",
			Body:  "the body",
		}, nil
	}
	msgId := "1"
	r := gin.Default()
	req, _ := http.NewRequest(http.MethodGet, "/messages/"+msgId, nil)
	rr := httptest.NewRecorder()
	r.GET("/messages/:message_id", GetMessage)
	r.ServeHTTP(rr, req)

	var message domain.Message
	err := json.Unmarshal(rr.Body.Bytes(), &message)
	assert.Nil(t, err)
	assert.NotNil(t, message)
	assert.EqualValues(t, http.StatusOK, rr.Code)
	assert.EqualValues(t, 1, message.Id)
	assert.EqualValues(t, "the title", message.Title)
	assert.EqualValues(t, "the body", message.Body)
}

func TestGetMessage_Invalid_Id(t *testing.T) {
	msgId := "abc"
	r := gin.Default()
	req, _ := http.NewRequest(http.MethodGet, "/messages/"+msgId, nil)
	rr := httptest.NewRecorder()
	r.GET("/messages/:message_id", GetMessage)
	r.ServeHTTP(rr, req)

	apiErr, err := error_utils.NewApiErrFromBytes(rr.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusBadRequest, apiErr.Status())
	assert.EqualValues(t, "message id should be a number", apiErr.Message())
	assert.EqualValues(t, "bad_request", apiErr.Error())
}

func TestGetMessage_Message_Not_Found(t *testing.T) {
	services.MessagesService = &serviceMock{}
	getMessageService = func(msgId int64) (*domain.Message, error_utils.MessageErr) {
		return nil, error_utils.NewNotFoundError("message not found")
	}
	msgId := "1"
	r := gin.Default()
	req, _ := http.NewRequest(http.MethodGet, "/messages/"+msgId, nil)
	rr := httptest.NewRecorder()
	r.GET("/messages/:message_id", GetMessage)
	r.ServeHTTP(rr, req)

	apiErr, err := error_utils.NewApiErrFromBytes(rr.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusNotFound, apiErr.Status())
	assert.EqualValues(t, "message not found", apiErr.Message())
	assert.EqualValues(t, "not_found", apiErr.Error())
}

func TestGetMessage_Message_Database_Error(t *testing.T) {
	services.MessagesService = &serviceMock{}
	getMessageService = func(msgId int64) (*domain.Message, error_utils.MessageErr) {
		return nil, error_utils.NewInternalServerError("database error")
	}
	msgId := "1"
	r := gin.Default()
	req, _ := http.NewRequest(http.MethodGet, "/messages/"+msgId, nil)
	rr := httptest.NewRecorder()
	r.GET("/messages/:message_id", GetMessage)
	r.ServeHTTP(rr, req)

	apiErr, err := error_utils.NewApiErrFromBytes(rr.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusInternalServerError, apiErr.Status())
	assert.EqualValues(t, "database error", apiErr.Message())
	assert.EqualValues(t, "server_error", apiErr.Error())
}

// "GetAllMessages" test cases

func TestGetAllMessages_Success(t *testing.T) {
	services.MessagesService = &serviceMock{}
	getAllMessageService = func() ([]domain.Message, error_utils.MessageErr) {
		return []domain.Message{
			{
				Id:    1,
				Title: "first title",
				Body:  "first body",
			},
			{
				Id:    2,
				Title: "second title",
				Body:  "second body",
			},
		}, nil
	}
	r := gin.Default()
	req, err := http.NewRequest(http.MethodGet, "/messages", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.GET("/messages", GetAllMessages)
	r.ServeHTTP(rr, req)

	var messages []domain.Message
	theErr := json.Unmarshal(rr.Body.Bytes(), &messages)
	if theErr != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	assert.Nil(t, err)
	assert.NotNil(t, messages)
	assert.EqualValues(t, messages[0].Id, 1)
	assert.EqualValues(t, messages[0].Title, "first title")
	assert.EqualValues(t, messages[0].Body, "first body")
	assert.EqualValues(t, messages[1].Id, 2)
	assert.EqualValues(t, messages[1].Title, "second title")
	assert.EqualValues(t, messages[1].Body, "second body")
}

func TestGetAllMessages_Failure(t *testing.T) {
	services.MessagesService = &serviceMock{}
	getAllMessageService = func() ([]domain.Message, error_utils.MessageErr) {
		return nil, error_utils.NewInternalServerError("error getting messages")
	}
	r := gin.Default()
	req, err := http.NewRequest(http.MethodGet, "/messages", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.GET("/messages", GetAllMessages)
	r.ServeHTTP(rr, req)

	apiErr, err := error_utils.NewApiErrFromBytes(rr.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusInternalServerError, apiErr.Status())
	assert.EqualValues(t, "error getting messages", apiErr.Message())
	assert.EqualValues(t, "server_error", apiErr.Error())
}
