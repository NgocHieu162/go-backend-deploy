package interfaces

import (
	"context"
	"go-backend/ent"
	"go-backend/internal/common/pagination"
	"go-backend/internal/dto"
)

type ChatGroupUsecase interface {
	FindAll(ctx context.Context, input dto.ChatGroupFindAllInput) (any, error)
}

type ChatGroupRepository interface {
	FindAll(ctx context.Context) (any, error)
	FindOneById(ctx context.Context, id int) (*ent.ChatGroups, error)
	CheckChatGroupOneOneExist(ctx context.Context, ids []int) (*ent.ChatGroups, error)
	CreateGroup(ctx context.Context, name string, userId int) (*ent.ChatGroups, error)
	GetAll(ctx context.Context, query pagination.Query, filters dto.ChatGroupFindAllFilters) ([]*ent.ChatGroups, error)
	Count(ctx context.Context, filters dto.ChatGroupFindAllFilters) (int, error)
}
