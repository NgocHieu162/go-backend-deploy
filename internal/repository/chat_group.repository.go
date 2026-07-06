package repository

import (
	"context"
	"go-backend/ent"
	"go-backend/ent/chatgroupmembers"
	"go-backend/ent/chatgroups"
	"go-backend/internal/common/pagination"
	"go-backend/internal/dto"
	"go-backend/internal/interfaces"
)

type ChatGroupRepository struct {
	entClient *ent.Client
}

func NewChatGroupRepository(entClient *ent.Client) interfaces.ChatGroupRepository {
	return &ChatGroupRepository{
		entClient: entClient,
	}
}

// FindAll implements [repository.ChatGroupRepository].
func (a *ChatGroupRepository) FindAll(ctx context.Context) (any, error) {
	return nil, nil
}

// CheckChatGroupOneOneExist implements [interfaces.ChatGroupRepository].
func (a *ChatGroupRepository) CheckChatGroupOneOneExist(ctx context.Context, ids []int) (*ent.ChatGroups, error) {
	entClient := a.entClient.ChatGroups.Query()
	entClient = entClient.Where(
		chatgroups.NameIsNil(),
		chatgroups.HasChatGroupMembersWith(
			chatgroupmembers.UserIDIn(ids...),
		),
	)

	return entClient.Only(ctx)
}

// CreateGroup implements [interfaces.ChatGroupRepository].
func (a *ChatGroupRepository) CreateGroup(ctx context.Context, name string, userId int) (*ent.ChatGroups, error) {
	entClientTx := GetClientTx(ctx, a.entClient)

	entCreate := entClientTx.ChatGroups.Create()

	if name != "" {
		entCreate = entCreate.SetName(name)
	}
	entCreate = entCreate.SetUserID(userId)

	return entCreate.Save(ctx)
}

// GetAll implements [interfaces.ChatGroupRepository].
func (a *ChatGroupRepository) GetAll(ctx context.Context, query pagination.Query, filters dto.ChatGroupFindAllFilters) ([]*ent.ChatGroups, error) {
	entQuery := a.entClient.ChatGroups.Query()

	if filters.Name != "" {
		entQuery = entQuery.Where(chatgroups.NameContainsFold(filters.Name))
	}

	entQuery = entQuery.WithUsers()
	entQuery = entQuery.WithChatGroupMembers(func(cgmq *ent.ChatGroupMembersQuery) {
		cgmq.WithUsers()
	})

	entQuery = entQuery.Limit(query.PageSize)
	entQuery = entQuery.Offset(query.Offset)

	entQuery = entQuery.Order(ent.Desc(chatgroups.FieldCreatedAt))

	return entQuery.All(ctx)
}

// Count implements [interfaces.ChatGroupRepository].
func (a *ChatGroupRepository) Count(ctx context.Context, filters dto.ChatGroupFindAllFilters) (int, error) {
	entQuery := a.entClient.ChatGroups.Query()
	if filters.Name != "" {
		entQuery = entQuery.Where(chatgroups.NameContainsFold(filters.Name))
	}
	return entQuery.Count(ctx)
}

// FindOneById implements [interfaces.ChatGroupRepository].
func (a *ChatGroupRepository) FindOneById(ctx context.Context, id int) (*ent.ChatGroups, error) {
	entQuery := a.entClient.ChatGroups.Query()
	entQuery = entQuery.Where(chatgroups.IDEQ(id))
	entQuery = entQuery.WithChatGroupMembers(func(cgmq *ent.ChatGroupMembersQuery) {
		cgmq.WithUsers()
	})

	return entQuery.Only(ctx)
}
