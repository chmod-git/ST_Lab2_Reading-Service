package app

import (
	"github.com/gin-gonic/gin"
	"testing-project/controllers"
)

func routes() {
	router.GET("/messages/:message_id", controllers.GetMessage)
	router.GET("/messages", controllers.GetAllMessages)
	router.GET("/health", func(c *gin.Context) {
		c.Status(200)
	})
}
