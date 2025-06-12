# Reading Service

**Reading Service** is a Go-based microservice for reading messages stored in Redis. It listens to RabbitMQ for message creation events and provides a REST API to fetch those messages.

### Technologies Used

* **Go** – Core language
* **Gin** – Web framework
* **Redis** – Message cache
* **RabbitMQ** – Message broker
* **Testify** – Unit testing
* **Redismock** – Redis mocking for tests

### Features

* List all messages: `GET /messages`
* Get message by ID: `GET /messages/:id`
* Consumes messages via RabbitMQ
* Fast reads via Redis caching

### Run Locally

1. Make sure Redis and RabbitMQ are running
2. Start the service:

   ```bash
   go run cmd/main.go
   ```

### Tests

```bash
go test ./... -v
```
