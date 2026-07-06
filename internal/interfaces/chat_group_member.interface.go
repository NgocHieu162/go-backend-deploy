package interfaces

import (
	"context"
)

type ChatGroupMembersUsecase interface {
	FindAll(ctx context.Context) (any, error)
}

type ChatGroupMembersRepository interface {
	FindAll(ctx context.Context) (any, error)
	CreateGroupMemberMany(ctx context.Context, chatGroupId int, ids []int) error
}
