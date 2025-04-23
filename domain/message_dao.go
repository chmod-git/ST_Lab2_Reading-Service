package domain

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"strconv"
	"testing-project/utils/error_utils"
)

var (
	MessageRepo messageRepoInterface = &messageRepo{}
	ctx                              = context.Background()
)

type messageRepoInterface interface {
	Get(int64) (*Message, error_utils.MessageErr)
	GetAll() ([]Message, error_utils.MessageErr)
	Save(*Message) error_utils.MessageErr
	Delete(int64) error_utils.MessageErr
	Initialize(string, string, string) *redis.Client
}

type messageRepo struct {
	client *redis.Client
}

func (mr *messageRepo) Initialize(addr, password, db string) *redis.Client {
	dbIndex, _ := strconv.Atoi(db)
	mr.client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       dbIndex,
	})

	_, err := mr.client.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Ошибка подключения к Redis:", err)
	}

	fmt.Println("Успешное подключение к Redis")
	return mr.client
}

func NewMessageRepository(client *redis.Client) messageRepoInterface {
	return &messageRepo{client: client}
}

func (mr *messageRepo) Get(messageId int64) (*Message, error_utils.MessageErr) {
	data, err := mr.client.Get(ctx, fmt.Sprintf("message:%d", messageId)).Result()
	if err == redis.Nil {
		return nil, error_utils.NewNotFoundError("message not found")
	} else if err != nil {
		return nil, error_utils.NewInternalServerError("redis get error")
	}

	var msg Message
	if err := json.Unmarshal([]byte(data), &msg); err != nil {
		return nil, error_utils.NewInternalServerError("json unmarshal error")
	}
	return &msg, nil
}

func (mr *messageRepo) GetAll() ([]Message, error_utils.MessageErr) {
	keys, err := mr.client.Keys(ctx, "message:*").Result()
	if err != nil {
		return nil, error_utils.NewInternalServerError("error fetching keys")
	}

	messages := make([]Message, 0)
	for _, key := range keys {
		data, err := mr.client.Get(ctx, key).Result()
		if err != nil {
			continue
		}
		var msg Message
		if err := json.Unmarshal([]byte(data), &msg); err != nil {
			continue
		}
		messages = append(messages, msg)
	}
	if len(messages) == 0 {
		return nil, error_utils.NewNotFoundError("no messages found")
	}
	return messages, nil
}

func (mr *messageRepo) Save(msg *Message) error_utils.MessageErr {
	data, err := json.Marshal(msg)
	if err != nil {
		return error_utils.NewInternalServerError("json marshal error")
	}
	key := fmt.Sprintf("message:%d", msg.Id)
	err = mr.client.Set(ctx, key, data, 0).Err()
	if err != nil {
		return error_utils.NewInternalServerError("redis save error")
	}
	return nil
}

func (mr *messageRepo) Delete(messageId int64) error_utils.MessageErr {
	key := fmt.Sprintf("message:%d", messageId)
	_, err := mr.client.Del(ctx, key).Result()
	if err != nil {
		return error_utils.NewInternalServerError("redis delete error")
	}
	return nil
}
