package repository

import (
	"context"
	"go-backend/ent"
	"go-backend/ent/chatmessages"
	"go-backend/internal/common/pagination"
	"go-backend/internal/dto"
	"go-backend/internal/interfaces"
	"time"
)

type ChatMessageRepository struct {
	entClient *ent.Client
}

func NewChatMessageRepository(entClient *ent.Client) interfaces.ChatMessageRepository {
	return &ChatMessageRepository{
		entClient: entClient,
	}
}

// GetAll implements [interfaces.ChatMessageRepository].
func (a *ChatMessageRepository) GetAll(ctx context.Context, query pagination.Query, filters dto.ChatMessageFindAllFilters) ([]*ent.ChatMessages, error) {
	entQuery := a.entClient.ChatMessages.Query()

	if filters.MessageText != "" {
		entQuery = entQuery.Where(chatmessages.MessageTextContainsFold(filters.MessageText))
	}

	entQuery = entQuery.Where(chatmessages.ChatGroupIDEQ(filters.ChatGroupId))
	entQuery = entQuery.WithUsers()
	entQuery = entQuery.Limit(query.PageSize)
	entQuery = entQuery.Offset(query.Offset)

	entQuery = entQuery.Order(ent.Desc(chatmessages.FieldCreatedAt))

	return entQuery.All(ctx)
}

// Count implements [interfaces.ChatMessageRepository].
func (a *ChatMessageRepository) Count(ctx context.Context, filters dto.ChatMessageFindAllFilters) (int, error) {
	entQuery := a.entClient.ChatMessages.Query()
	if filters.MessageText != "" {
		entQuery = entQuery.Where(chatmessages.MessageTextContainsFold(filters.MessageText))
	}
	entQuery = entQuery.Where(chatmessages.ChatGroupIDEQ(filters.ChatGroupId))

	return entQuery.Count(ctx)
}

// CreateMessage implements [interfaces.ChatMessageRepository].
func (a *ChatMessageRepository) CreateMessage(ctx context.Context, userId int, chatGroupId int, messageText string, createdAt time.Time) (*ent.ChatMessages, error) {
	entCreate := a.entClient.ChatMessages.Create()

	entCreate = entCreate.SetChatGroupID(chatGroupId)
	entCreate = entCreate.SetUserID(userId)
	entCreate = entCreate.SetMessageText(messageText)
	entCreate = entCreate.SetCreatedAt(createdAt)

	return entCreate.Save(ctx)
}
