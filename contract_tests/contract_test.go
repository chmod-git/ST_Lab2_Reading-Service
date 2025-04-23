package contract_tests

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"net/http"
	"testing"
	"time"
)

type Message struct {
	Id        int64     `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}

func TestGetMessage(t *testing.T) {
	baseURL := "http://localhost:8090"
	messageID := 1

	resp, err := http.Get(fmt.Sprintf("%s/messages/%d", baseURL, messageID))
	if err != nil {
		t.Fatalf("Помилка при надсиланні запиту: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Очікувався статус 200 OK, отримано %d", resp.StatusCode)
	}

	var msg Message
	if err := json.NewDecoder(resp.Body).Decode(&msg); err != nil {
		t.Fatalf("Помилка при декодуванні відповіді: %v", err)
	}

	if msg.Id != int64(messageID) {
		t.Errorf("Очікувався ID %d, отримано %d", messageID, msg.Id)
	}
}

func TestGetAllMessages(t *testing.T) {
	baseURL := "http://localhost:8090"

	resp, err := http.Get(fmt.Sprintf("%s/messages", baseURL))
	if err != nil {
		t.Fatalf("Помилка при надсиланні запиту: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Очікувався статус 200 OK, отримано %d", resp.StatusCode)
	}

	var messages []Message
	if err := json.NewDecoder(resp.Body).Decode(&messages); err != nil {
		t.Fatalf("Помилка при декодуванні списку повідомлень: %v", err)
	}

	if len(messages) == 0 {
		t.Errorf("Очікувався непорожній список повідомлень")
	}
}

func TestGetMessageNotFound(t *testing.T) {
	baseURL := "http://localhost:8090"
	messageID := 9999

	resp, err := http.Get(fmt.Sprintf("%s/messages/%d", baseURL, messageID))
	if err != nil {
		t.Fatalf("Помилка при надсиланні запиту: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Очікувався статус 404 Not Found, отримано %d", resp.StatusCode)
	}
}

func publishEvent(eventType string, message Message) error {
	event := map[string]interface{}{
		"event": eventType,
		"data":  message,
	}
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	return ch.Publish(
		"",
		"my_queue",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		},
	)
}

func TestCreateMessageEvent(t *testing.T) {
	baseURL := "http://localhost:8090"
	unique := time.Now().UnixNano()
	title := fmt.Sprintf("Event title %d", unique)
	body := fmt.Sprintf("Event body %d", unique)

	message := Message{
		Id:        unique,
		Title:     title,
		Body:      body,
		CreatedAt: time.Now(),
	}
	err := publishEvent("created", message)

	if err != nil {
		t.Fatalf("Failed to publish create event: %v", err)
	}
	time.Sleep(2 * time.Second)

	resp, err := http.Get(fmt.Sprintf("%s/messages/%d", baseURL, message.Id))
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var received Message
	if err := json.NewDecoder(resp.Body).Decode(&received); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if received.Id != message.Id {
		t.Errorf("Expected ID %d, got %d", message.Id, received.Id)
	}
	if received.Title != message.Title {
		t.Errorf("Expected Title %q, got %q", message.Title, received.Title)
	}
	if received.Body != message.Body {
		t.Errorf("Expected Body %q, got %q", message.Body, received.Body)
	}
}

func TestUpdateMessageEvent(t *testing.T) {
	baseURL := "http://localhost:8090"
	unique := time.Now().UnixNano()
	title := fmt.Sprintf("Event title %d", unique)
	body := fmt.Sprintf("Event body %d", unique)

	message := Message{
		Id:        unique,
		Title:     title,
		Body:      body,
		CreatedAt: time.Now(),
	}

	_ = publishEvent("created", message)
	time.Sleep(2 * time.Second)

	message.Title += " Updated"
	message.Body += " Updated"
	err := publishEvent("updated", message)
	if err != nil {
		t.Fatalf("Failed to publish update event: %v", err)
	}
	time.Sleep(2 * time.Second)

	resp, err := http.Get(fmt.Sprintf("%s/messages/%d", baseURL, message.Id))
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var updated Message
	if err := json.NewDecoder(resp.Body).Decode(&updated); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if updated.Id != message.Id {
		t.Errorf("Expected ID %d, got %d", message.Id, updated.Id)
	}
	if updated.Title != message.Title {
		t.Errorf("Expected Title %q, got %q", message.Title, updated.Title)
	}
	if updated.Body != message.Body {
		t.Errorf("Expected Body %q, got %q", message.Body, updated.Body)
	}
}

func TestDeleteMessageEvent(t *testing.T) {
	baseURL := "http://localhost:8090"
	unique := time.Now().UnixNano()
	title := fmt.Sprintf("Event title %d", unique)
	body := fmt.Sprintf("Event body %d", unique)

	message := Message{
		Id:        unique,
		Title:     title,
		Body:      body,
		CreatedAt: time.Now(),
	}

	_ = publishEvent("created", message)
	time.Sleep(2 * time.Second)

	err := publishEvent("deleted", message)
	if err != nil {
		t.Fatalf("Failed to publish delete event: %v", err)
	}
	time.Sleep(2 * time.Second)

	resp, err := http.Get(fmt.Sprintf("%s/messages/%d", baseURL, message.Id))
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status 404 Not Found, got %d", resp.StatusCode)
	}
}
