package usecase

import (
	"context"
	"go-backend/internal/common/pagination"
	"go-backend/internal/common/response"
	"go-backend/internal/dto"
	"go-backend/internal/interfaces"
	"math"
)

type ChatGroupUsecase struct {
	ChatGroupRepository interfaces.ChatGroupRepository
}

func NewChatGroupUsecase(chatGroupRepository interfaces.ChatGroupRepository) interfaces.ChatGroupUsecase {
	return &ChatGroupUsecase{
		ChatGroupRepository: chatGroupRepository,
	}
}

// FindAll implements [usecase.ChatGroupUsecase].
func (a *ChatGroupUsecase) FindAll(ctx context.Context, input dto.ChatGroupFindAllInput) (any, error) {
	data, err := a.ChatGroupRepository.GetAll(ctx, input.Query, input.ChatGroupFindAllFilters)

	if err != nil {
		return nil, response.NewBadRequestException(err.Error())
	}

	// totalItem: tổng số lượng item
	totalItem, err := a.ChatGroupRepository.Count(ctx, input.ChatGroupFindAllFilters)
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
