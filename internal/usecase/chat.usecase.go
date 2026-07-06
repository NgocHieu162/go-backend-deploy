package usecase

import (
	"context"
	"errors"
	"fmt"
	"go-backend/ent"
	"go-backend/internal/dto"
	"go-backend/internal/interfaces"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type ChatUsecase struct {
	tokenUsecase               interfaces.TokenUsecase
	userRepository             interfaces.UserRepository
	chatMessageRepository      interfaces.ChatMessageRepository
	chatGroupRepository        interfaces.ChatGroupRepository
	chatGroupMembersRepository interfaces.ChatGroupMembersRepository
	unitOfWorkRepository       interfaces.UnitOfWorkRepository
}

func NewChatUsecase(tokenUsecase interfaces.TokenUsecase, userRepository interfaces.UserRepository, chatGroupRepository interfaces.ChatGroupRepository, chatGroupMembersRepository interfaces.ChatGroupMembersRepository, chatMessageRepository interfaces.ChatMessageRepository, unitOfWorkRepository interfaces.UnitOfWorkRepository) interfaces.ChatUsecase {
	return &ChatUsecase{
		tokenUsecase:               tokenUsecase,
		userRepository:             userRepository,
		chatMessageRepository:      chatMessageRepository,
		chatGroupRepository:        chatGroupRepository,
		chatGroupMembersRepository: chatGroupMembersRepository,
		unitOfWorkRepository:       unitOfWorkRepository,
	}
}

// CreateGroup implements [interfaces.ChatUsecase].
func (c *ChatUsecase) CreateGroup(ctx context.Context, accessToken string, targetUserIds []int, name string) (*ent.ChatGroups, error) {
	claimAccessToken, err := c.tokenUsecase.VerifyAccessToken(accessToken, jwt.WithoutClaimsValidation())
	if err != nil {
		return nil, err
	}

	user, err := c.userRepository.FindUserById(ctx, claimAccessToken.UserID)
	if err != nil {
		return nil, err
	}
	// fmt.Println("user", user)
	targetUserIds = append(targetUserIds, user.ID)

	targetUserIdsUnique := []int{}
	seen := map[int]bool{}
	for _, id := range targetUserIds {
		if seen[id] {
			continue
		}
		seen[id] = true
		targetUserIdsUnique = append(targetUserIdsUnique, id)
	}

	var chatGroupExists *ent.ChatGroups
	if len(targetUserIdsUnique) == 2 {
		// chat 1 - 1
		chatGroupExists, err = c.chatGroupRepository.CheckChatGroupOneOneExist(ctx, targetUserIdsUnique)
		if err != nil && !ent.IsNotFound(err) {
			return nil, err
		}

		if chatGroupExists == nil {
			err = c.unitOfWorkRepository.Do(ctx, func(ctxTx context.Context) error {
				chatGroupExists, err := c.chatGroupRepository.CreateGroup(ctxTx, "", user.ID)
				if err != nil {
					return err
				}

				err = c.chatGroupMembersRepository.CreateGroupMemberMany(ctxTx, chatGroupExists.ID, targetUserIdsUnique)
				if err != nil {
					return err
				}

				return nil
			})

			if err != nil {
				return nil, err
			}

		}

	} else {
		err = c.unitOfWorkRepository.Do(ctx, func(ctxTx context.Context) error {
			chatGroupExists, err = c.chatGroupRepository.CreateGroup(ctxTx, name, user.ID)
			if err != nil {
				return err
			}

			err = c.chatGroupMembersRepository.CreateGroupMemberMany(ctxTx, chatGroupExists.ID, targetUserIdsUnique)
			if err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			return nil, err
		}
	}

	return chatGroupExists, nil
}

// JoinGroup implements [interfaces.ChatUsecase].
func (c *ChatUsecase) JoinGroup(ctx context.Context, accessToken string, chatGroupId int) (*ent.ChatGroups, error) {
	claimAccessToken, err := c.tokenUsecase.VerifyAccessToken(accessToken, jwt.WithoutClaimsValidation())
	if err != nil {
		return nil, err
	}

	user, err := c.userRepository.FindUserById(ctx, claimAccessToken.UserID)
	if err != nil {
		return nil, err
	}

	chatGroup, err := c.chatGroupRepository.FindOneById(ctx, chatGroupId)
	if err != nil {
		return nil, err
	}

	isUserInChatGroup := false
	for _, member := range chatGroup.Edges.ChatGroupMembers {
		if member.Edges.Users.ID == user.ID {
			isUserInChatGroup = true
		}
	}

	if !isUserInChatGroup {
		return nil, errors.New("User not exists in chat group")
	}
	return chatGroup, nil
}

// SendMessage implements [interfaces.ChatUsecase].
func (c *ChatUsecase) SendMessage(ctx context.Context, accessToken string, chatGroupId int, message string, createdAt time.Time) (*dto.SendMessageReturn, error) {
	claimAccessToken, err := c.tokenUsecase.VerifyAccessToken(accessToken, jwt.WithoutClaimsValidation())
	if err != nil {
		return nil, err
	}

	user, err := c.userRepository.FindUserById(ctx, claimAccessToken.UserID)
	if err != nil {
		return nil, err
	}

	chatGroup, err := c.chatGroupRepository.FindOneById(ctx, chatGroupId)
	if err != nil {
		return nil, err
	}

	isUserInChatGroup := false
	for _, member := range chatGroup.Edges.ChatGroupMembers {
		if member.Edges.Users.ID == user.ID {
			isUserInChatGroup = true
		}
	}

	if !isUserInChatGroup {
		return nil, errors.New("User not exists in chat group")
	}

	go func() {
		_, err = c.chatMessageRepository.CreateMessage(ctx, user.ID, chatGroupId, message, createdAt)
		if err != nil {
			fmt.Println(err)
			return
		}
	}()

	return &dto.SendMessageReturn{
		MessageText: message,
		ChatGroupId: chatGroupId,
		UserId:      user.ID,
	}, nil
}
