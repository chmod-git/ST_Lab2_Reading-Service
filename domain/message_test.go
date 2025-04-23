package domain_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
	"testing-project/domain"
)

func TestGetMessage_Success(t *testing.T) {
	db, mock := redismock.NewClientMock()
	repo := domain.NewMessageRepository(db)

	msg := domain.Message{
		Id:        1,
		Title:     "Test",
		Body:      "Message",
		CreatedAt: time.Now(),
	}
	data, _ := json.Marshal(msg)

	mock.ExpectGet("message:1").SetVal(string(data))

	result, err := repo.Get(1)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, msg.Id, result.Id)
	assert.Equal(t, msg.Title, result.Title)
}

func TestGetMessage_NotFound(t *testing.T) {
	db, mock := redismock.NewClientMock()
	repo := domain.NewMessageRepository(db)

	mock.ExpectGet("message:1").RedisNil()

	result, err := repo.Get(1)

	assert.Nil(t, result)
	assert.Equal(t, "message not found", err.Message())
}

func TestGetAllMessages_Success(t *testing.T) {
	db, mock := redismock.NewClientMock()
	repo := domain.NewMessageRepository(db)

	msg := domain.Message{
		Id:        2,
		Title:     "Title",
		Body:      "Body",
		CreatedAt: time.Now(),
	}
	data, _ := json.Marshal(msg)

	mock.ExpectKeys("message:*").SetVal([]string{"message:2"})
	mock.ExpectGet("message:2").SetVal(string(data))

	result, err := repo.GetAll()

	assert.Nil(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, msg.Id, result[0].Id)
}

func TestGetAllMessages_Empty(t *testing.T) {
	db, mock := redismock.NewClientMock()
	repo := domain.NewMessageRepository(db)

	mock.ExpectKeys("message:*").SetVal([]string{})

	result, err := repo.GetAll()

	assert.Nil(t, result)
	assert.Equal(t, "no messages found", err.Message())
}

func TestSaveMessage_Success(t *testing.T) {
	db, mock := redismock.NewClientMock()
	repo := domain.NewMessageRepository(db)

	msg := &domain.Message{
		Id:        10,
		Title:     "Hello",
		Body:      "World",
		CreatedAt: time.Now(),
	}
	data, _ := json.Marshal(msg)
	key := "message:10"

	mock.ExpectSet(key, data, 0).SetVal("OK")

	err := repo.Save(msg)

	assert.Nil(t, err)
}

func TestDeleteMessage_Success(t *testing.T) {
	db, mock := redismock.NewClientMock()
	repo := domain.NewMessageRepository(db)

	mock.ExpectDel("message:12").SetVal(1)

	err := repo.Delete(12)

	assert.Nil(t, err)
}
