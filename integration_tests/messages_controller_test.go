package integration_tests

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing-project/controllers"
	"testing-project/domain"
	"testing-project/utils/error_utils"
	"time"
)

type mockMessageRepo struct {
	mock.Mock
}

func (m *mockMessageRepo) Get(id int64) (*domain.Message, error_utils.MessageErr) {
	args := m.Called(id)

	var msg *domain.Message
	if args.Get(0) != nil {
		msg = args.Get(0).(*domain.Message)
	}

	var err error_utils.MessageErr
	if args.Get(1) != nil {
		err = args.Get(1).(error_utils.MessageErr)
	}

	return msg, err
}
func (m *mockMessageRepo) GetAll() ([]domain.Message, error_utils.MessageErr) {
	args := m.Called()

	var messages []domain.Message
	if args.Get(0) != nil {
		messages = args.Get(0).([]domain.Message)
	}

	var err error_utils.MessageErr
	if args.Get(1) != nil {
		err = args.Get(1).(error_utils.MessageErr)
	}

	return messages, err
}
func (m *mockMessageRepo) Save(msg *domain.Message) error_utils.MessageErr {
	args := m.Called(msg)
	return args.Get(0).(error_utils.MessageErr)
}
func (m *mockMessageRepo) Delete(id int64) error_utils.MessageErr {
	args := m.Called(id)
	return args.Get(0).(error_utils.MessageErr)
}
func (m *mockMessageRepo) Initialize(a, b, c string) *redis.Client { return nil }

func TestGetMessage_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	expected := &domain.Message{
		Id:        1,
		Title:     "Test Title",
		Body:      "Test Body",
		CreatedAt: time.Now(),
	}

	mockRepo := new(mockMessageRepo)
	mockRepo.On("Get", int64(1)).Return(expected, nil)
	domain.MessageRepo = mockRepo

	req, _ := http.NewRequest(http.MethodGet, "/messages/1", nil)
	resp := httptest.NewRecorder()

	router := gin.Default()
	router.GET("/messages/:message_id", controllers.GetMessage)

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var actual domain.Message
	json.Unmarshal(resp.Body.Bytes(), &actual)
	assert.Equal(t, expected.Title, actual.Title)
	assert.Equal(t, expected.Body, actual.Body)
}

func TestGetMessage_InvalidId(t *testing.T) {
	gin.SetMode(gin.TestMode)

	req, _ := http.NewRequest(http.MethodGet, "/messages/abc", nil)
	resp := httptest.NewRecorder()

	router := gin.Default()
	router.GET("/messages/:message_id", controllers.GetMessage)

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestGetMessage_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(mockMessageRepo)
	mockRepo.On("Get", int64(42)).Return(nil, error_utils.NewNotFoundError("message not found"))

	domain.MessageRepo = mockRepo

	req, _ := http.NewRequest(http.MethodGet, "/messages/42", nil)
	resp := httptest.NewRecorder()

	router := gin.Default()
	router.GET("/messages/:message_id", controllers.GetMessage)

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
}

func TestGetAllMessages_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(mockMessageRepo)

	expectedMessages := []domain.Message{
		{
			Id:        1,
			Title:     "First Message",
			Body:      "This is the first message",
			CreatedAt: time.Now(),
		},
		{
			Id:        2,
			Title:     "Second Message",
			Body:      "This is the second message",
			CreatedAt: time.Now(),
		},
	}

	mockRepo.On("GetAll").Return(expectedMessages, nil)
	domain.MessageRepo = mockRepo

	req, _ := http.NewRequest(http.MethodGet, "/messages", nil)
	resp := httptest.NewRecorder()

	router := gin.Default()
	router.GET("/messages", controllers.GetAllMessages)

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var actual []domain.Message
	json.Unmarshal(resp.Body.Bytes(), &actual)
	assert.Equal(t, 2, len(actual))
	assert.Equal(t, expectedMessages[0].Title, actual[0].Title)
	assert.Equal(t, expectedMessages[0].Body, actual[0].Body)
	assert.Equal(t, expectedMessages[1].Title, actual[1].Title)
	assert.Equal(t, expectedMessages[1].Body, actual[1].Body)
}

func TestGetAllMessages_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(mockMessageRepo)
	mockRepo.On("GetAll").Return(nil, error_utils.NewNotFoundError("no messages found"))
	domain.MessageRepo = mockRepo

	req, _ := http.NewRequest(http.MethodGet, "/messages", nil)
	resp := httptest.NewRecorder()

	router := gin.Default()
	router.GET("/messages", controllers.GetAllMessages)

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
}
