package repository

import (
	"context"
	"go-backend/ent"
	"go-backend/internal/interfaces"
)

type ChatGroupMembersRepository struct {
	entClient *ent.Client
}

func NewChatGroupMembersRepository(entClient *ent.Client) interfaces.ChatGroupMembersRepository {
	return &ChatGroupMembersRepository{
		entClient: entClient,
	}
}

// FindAll implements [repository.ChatGroupMembersRepository].
func (a *ChatGroupMembersRepository) FindAll(ctx context.Context) (any, error) {
	return nil, nil
}

// CreateGroupMemberMany implements [interfaces.ChatGroupMembersRepository].
func (a *ChatGroupMembersRepository) CreateGroupMemberMany(ctx context.Context, chatGroupId int, ids []int) error {
	entClientTx := GetClientTx(ctx, a.entClient)

	builders := []*ent.ChatGroupMembersCreate{}
	for _, id := range ids {
		builders = append(builders, entClientTx.ChatGroupMembers.Create().SetUserID(id).SetChatGroupID(chatGroupId))
	}

	return entClientTx.ChatGroupMembers.CreateBulk(builders...).Exec(ctx)
}
