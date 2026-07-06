package handler

import (
	"context"
	"fmt"
	"go-backend/internal/dto"
	"go-backend/internal/interfaces"
	"time"

	server "github.com/zishang520/socket.io/servers/socket/v3"
)

type ChatHandler struct {
	chatUsecase interfaces.ChatUsecase
}

func NewChatHandler(chatUsecase interfaces.ChatUsecase) *ChatHandler {
	return &ChatHandler{
		chatUsecase: chatUsecase,
	}
}

func (c *ChatHandler) CreateGroup(args ...any) {
	payload := args[0].(map[string]interface{})
	accessToken := payload["accessToken"].(string)
	targetUserIdsAny := payload["targetUserIds"].([]interface{})

	name := ""
	nameAny := payload["name"]
	if nameAny != nil {
		name = nameAny.(string)
	}

	ack := args[1].(func([]interface{}, error))

	targetUserIds := []int{}
	for _, userIdAny := range targetUserIdsAny {
		// fmt.Printf("received: %T | %v\n", userIdAny, userIdAny)
		userId := int(userIdAny.(float64))
		targetUserIds = append(targetUserIds, userId)
	}

	chatGroup, err := c.chatUsecase.CreateGroup(context.Background(), accessToken, targetUserIds, name)
	if err != nil {
		res := dto.ChatRes{
			Status:  "error",
			Message: err.Error(),
			Data:    nil,
		}
		ack([]interface{}{res}, err)
	}

	res := dto.ChatRes{
		Status:  "success",
		Message: "Create Chat Group Success",
		Data: map[string]any{
			"chatGroupId": chatGroup.ID,
		},
	}
	ack([]interface{}{res}, err)
}

func (c *ChatHandler) JoinGroup(socket *server.Socket, args ...any) {
	payload := args[0].(map[string]interface{})
	accessToken := payload["accessToken"].(string)
	chatGroupId := int(payload["chatGroupId"].(float64))

	ack := args[1].(func([]interface{}, error))

	chatGroup, err := c.chatUsecase.JoinGroup(context.Background(), accessToken, chatGroupId)
	if err != nil {
		res := dto.ChatRes{
			Status:  "error",
			Message: err.Error(),
			Data:    nil,
		}
		ack([]interface{}{res}, err)
	}

	socket.Join(createRoomName(chatGroup.ID))

	res := dto.ChatRes{
		Status:  "success",
		Message: "Join Group Success",
		Data: map[string]any{
			"chatGroupId": chatGroup.ID,
		},
	}
	ack([]interface{}{res}, err)
}

func (c *ChatHandler) SendMessage(io *server.Server, args ...any) {
	payload := args[0].(map[string]interface{})
	accessToken := payload["accessToken"].(string)
	chatGroupId := int(payload["chatGroupId"].(float64))
	message := payload["message"].(string)
	ack := args[1].(func([]interface{}, error))

	createdAt := time.Now().UTC()

	result, err := c.chatUsecase.SendMessage(context.Background(), accessToken, chatGroupId, message, createdAt)
	if err != nil {
		res := dto.ChatRes{
			Status:  "error",
			Message: err.Error(),
			Data:    nil,
		}
		ack([]interface{}{res}, err)
	}

	io.To(createRoomName(result.ChatGroupId)).Emit("SEND_MESSAGE", map[string]any{
		"messageText": result.MessageText,
		"userId":      result.UserId,
		"chatGroupId": result.ChatGroupId,
		"createdAt":   createdAt.Format(time.RFC3339),
	})
}

func createRoomName(chatGroupId int) server.Room {
	return server.Room(fmt.Sprintf("chat:%d", chatGroupId))
}
