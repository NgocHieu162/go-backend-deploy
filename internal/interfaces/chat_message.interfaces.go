package interfaces

import (
	"context"
	"go-backend/ent"
	"go-backend/internal/common/pagination"
	"go-backend/internal/dto"
	"time"
)

type ChatMessageUsecase interface {
	FindAll(ctx context.Context, input dto.ChatMessageFindAllInput) (any, error)
}

type ChatMessageRepository interface {
	GetAll(ctx context.Context, query pagination.Query, filters dto.ChatMessageFindAllFilters) ([]*ent.ChatMessages, error)
	Count(ctx context.Context, filters dto.ChatMessageFindAllFilters) (int, error)
	CreateMessage(ctx context.Context, userId int, chatGroupId int, messageText string, createdAt time.Time) (*ent.ChatMessages, error)
}
