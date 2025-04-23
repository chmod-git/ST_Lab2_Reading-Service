package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
	"testing-project/domain"
)

var (
	router = gin.Default()
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print(".env файл не знайдено, використовується дефолт")
	}
}

func StartApp() {
	brokerAddr := os.Getenv("RABBITMQ_URL")
	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB := os.Getenv("REDIS_DB")

	domain.MessageRepo.Initialize(redisAddr, redisPassword, redisDB)
	fmt.Println("Redis успішно ініціалізовано")

	go startRabbitListener(brokerAddr)

	routes()

	router.Run(":8090")
}
