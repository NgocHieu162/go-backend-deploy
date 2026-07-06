package usecase

import (
	"context"
	"go-backend/internal/common/pagination"
	"go-backend/internal/common/response"
	"go-backend/internal/dto"
	"go-backend/internal/interfaces"
	"math"
)

type ChatMessageUsecase struct {
	chatMessageRepository interfaces.ChatMessageRepository
}

func NewChatMessageUsecase(chatMessageRepository interfaces.ChatMessageRepository) interfaces.ChatMessageUsecase {
	return &ChatMessageUsecase{
		chatMessageRepository: chatMessageRepository,
	}
}

// FindAll implements [usecase.chatMessageUsecase].
func (a *ChatMessageUsecase) FindAll(ctx context.Context, input dto.ChatMessageFindAllInput) (any, error) {
	data, err := a.chatMessageRepository.GetAll(ctx, input.Query, input.ChatMessageFindAllFilters)

	if err != nil {
		return nil, response.NewBadRequestException(err.Error())
	}

	// totalItem: tổng số lượng item
	totalItem, err := a.chatMessageRepository.Count(ctx, input.ChatMessageFindAllFilters)
	if err != nil {
		return nil, response.NewBadRequestException(err.Error())
	}

	// totalPage: tổng số trang
	totalPage := int(math.Ceil(float64(totalItem) / float64(input.Query.PageSize)))

	res := pagination.PaginationRes[any]{
		Items:     data,
		Page:      input.Query.Page,
		PageSize:  input.Query.PageSize,
		TotalItem: totalItem,
		TotalPage: totalPage,
	}

	return &res, nil
}
