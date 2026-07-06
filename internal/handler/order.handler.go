package handler

import (
	"go-backend/internal/common/helpers"
	"go-backend/internal/common/response"
	"go-backend/internal/dto"
	"go-backend/internal/interfaces"
	"io"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderUsecase interfaces.OrderUsecase
}

func NewOrderHandler(orderUsecase interfaces.OrderUsecase) *OrderHandler {
	return &OrderHandler{
		orderUsecase: orderUsecase,
	}
}

func (a *OrderHandler) FindAll(ctx *gin.Context) {
	result, err := a.orderUsecase.FindAll(ctx.Request.Context())
	if err != nil {
		ctx.Error(err)
		return
	}

	response.Success(result, "", 0, ctx)
}

func (a *OrderHandler) Create(ctx *gin.Context) {
	user, err := helpers.GetUsers(ctx)
	if err != nil {
		ctx.Error(response.NewBadRequestException(err.Error()))
		return
	}

	var body dto.CreateOrder
	err = ctx.ShouldBindJSON(&body)
	if err != nil {
		if err == io.EOF {
			ctx.Error(response.NewBadRequestException("Body required"))
			return
		}
		ctx.Error(response.NewBadRequestException(err.Error()))
		return
	}

	body.UserId = user.ID

	result, err := a.orderUsecase.CreateRequest(ctx.Request.Context(), body)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.Success(result, "", 0, ctx)
}
