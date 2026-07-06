package usecase

import (
	"context"
	"go-backend/ent"
	"go-backend/internal/common/rabbitmq"
	"go-backend/internal/common/response"
	"go-backend/internal/dto"
	"go-backend/internal/interfaces"
)

type OrderUsecase struct {
	rabbitmq *rabbitmq.RabbitMQ
}

func NewOrderUsecase(rabbitmq *rabbitmq.RabbitMQ) interfaces.OrderUsecase {
	return &OrderUsecase{
		rabbitmq: rabbitmq,
	}
}

// FindAll implements [usecase.OrderUsecase].
func (a *OrderUsecase) FindAll(ctx context.Context) (any, error) {
	return "FindAll", nil
}

// CreateSend implements [interfaces.OrderUsecase].
func (a *OrderUsecase) CreateSend(ctx context.Context, body dto.CreateOrder) (any, error) {
	err := a.rabbitmq.Send(ctx, "CREATE_ORDER_SEND", body)
	if err != nil {
		return nil, response.NewBadRequestException(err.Error())
	}

	return true, nil
}

// CreateRequest implements [interfaces.OrderUsecase].
func (a *OrderUsecase) CreateRequest(ctx context.Context, body dto.CreateOrder) (any, error) {
	var result *ent.Orders

	err := a.rabbitmq.Request(
		ctx,
		"CREATE_ORDER_REQUEST",
		body,
		&result,
	)
	if err != nil {
		return nil, response.NewBadRequestException(err.Error())
	}
	return result, nil
}
