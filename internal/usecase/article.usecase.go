package usecase

import (
	"context"
	"go-backend/internal/common/pagination"
	"go-backend/internal/common/response"
	"go-backend/internal/dto"
	"go-backend/internal/interfaces"
	"math"
)

type articleUsecase struct {
	articleRepository interfaces.ArticleRepository
}

func NewArticleUsecase(articleRepository interfaces.ArticleRepository) interfaces.ArticleUsecase {
	return &articleUsecase{
		articleRepository: articleRepository,
	}
}

func (a *articleUsecase) Create(ctx context.Context, body dto.ArticleCreateReq) (any, error) {
	data, err := a.articleRepository.Create(ctx, body)

	if err != nil {
		return nil, response.NewBadRequestException(err.Error())
	}

	return data, nil
}

func (a *articleUsecase) FindAll(ctx context.Context, input dto.ArticleFindAllInput) (*pagination.PaginationRes[any], error) {
	data, err := a.articleRepository.GetAll(ctx, input.Query, input.ArticleFindAllFilters)

	if err != nil {
		return nil, response.NewBadRequestException(err.Error())
	}

	// totalItem: tổng số lượng item
	totalItem, err := a.articleRepository.Count(ctx, input.ArticleFindAllFilters)
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

// Delete implements [interfaces.ArticleUsecase].
func (a *articleUsecase) Delete(ctx context.Context, id int) (any, error) {
	err := a.articleRepository.Delete(ctx, id)
	if err != nil{
		return nil, response.NewBadRequestException(err.Error())
	}
	
	return true, nil
}
