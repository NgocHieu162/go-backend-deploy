package handler

import (
	"encoding/json"
	"fmt"
	"go-backend/internal/common/pagination"
	"go-backend/internal/common/response"
	"go-backend/internal/dto"
	"go-backend/internal/interfaces"

	"github.com/gin-gonic/gin"
)

type ChatMessageHandler struct {
	chatMessageUsecase interfaces.ChatMessageUsecase
}

func NewChatMessageHandler(chatMessageUsecase interfaces.ChatMessageUsecase) *ChatMessageHandler {
	return &ChatMessageHandler{
		chatMessageUsecase: chatMessageUsecase,
	}
}

func (a *ChatMessageHandler) FindAll(ctx *gin.Context) {
	queryPagination := pagination.Get(ctx.Query("page"), ctx.Query("pageSize"))

	filterString := ctx.DefaultQuery("filters", "{}")
	var filters dto.ChatMessageFindAllFilters
	json.Unmarshal([]byte(filterString), &filters)
	fmt.Printf("%+v \n\n", filters)

	input := dto.ChatMessageFindAllInput{
		Query:                     *queryPagination,
		ChatMessageFindAllFilters: filters,
	}
	fmt.Printf("%+v \n\n", queryPagination)

	result, err := a.chatMessageUsecase.FindAll(ctx.Request.Context(), input)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.Success(result, "", 0, ctx)
}
