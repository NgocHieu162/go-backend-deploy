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

type ChatGroupHandler struct {
	chatGroupUsecase interfaces.ChatGroupUsecase
}

func NewChatGroupHandler(chatGroupUsecase interfaces.ChatGroupUsecase) *ChatGroupHandler {
	return &ChatGroupHandler{
		chatGroupUsecase: chatGroupUsecase,
	}
}

func (a *ChatGroupHandler) FindAll(ctx *gin.Context) {
	queryPagination := pagination.Get(ctx.Query("page"), ctx.Query("pageSize"))

	filterString := ctx.DefaultQuery("filters", "{}")
	var filters dto.ChatGroupFindAllFilters
	json.Unmarshal([]byte(filterString), &filters)
	fmt.Printf("%+v \n\n", filters)

	input := dto.ChatGroupFindAllInput{
		Query:                   *queryPagination,
		ChatGroupFindAllFilters: filters,
	}
	fmt.Printf("%+v \n\n", queryPagination)

	result, err := a.chatGroupUsecase.FindAll(ctx.Request.Context(), input)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.Success(result, "", 0, ctx)
}
