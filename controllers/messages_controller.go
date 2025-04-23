package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"testing-project/services"
	"testing-project/utils/error_utils"
)

func getMessageId(msgIdParam string) (int64, error_utils.MessageErr) {
	msgId, msgErr := strconv.ParseInt(msgIdParam, 10, 64)
	if msgErr != nil {
		return 0, error_utils.NewBadRequestError("message id should be a number")
	}
	return msgId, nil
}

func GetMessage(c *gin.Context) {
	msgId, err := getMessageId(c.Param("message_id"))
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	message, getErr := services.MessagesService.GetMessage(msgId)
	if getErr != nil {
		c.JSON(getErr.Status(), getErr)
		return
	}
	c.JSON(http.StatusOK, message)
}

func GetAllMessages(c *gin.Context) {
	messages, getErr := services.MessagesService.GetAllMessages()
	if getErr != nil {
		c.JSON(getErr.Status(), getErr)
		return
	}
	c.JSON(http.StatusOK, messages)
}
