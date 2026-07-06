package interfaces

import (
	"context"
	"go-backend/ent"
	"go-backend/internal/common/pagination"
	"go-backend/internal/dto"
)

type ArticleUsecase interface {
	Create(ctx context.Context, body dto.ArticleCreateReq) (any, error)
	FindAll(ctx context.Context, input dto.ArticleFindAllInput) (*pagination.PaginationRes[any], error)
	Delete(ctx context.Context, id int) (any, error)
}

type ArticleRepository interface {
	Create(ctx context.Context, body dto.ArticleCreateReq) (*ent.Articles, error)
	// CreateGorm(ctx context.Context, body dto.ArticleCreateReq) (any, error)
	GetAll(ctx context.Context, query pagination.Query, filters dto.ArticleFindAllFilters) ([]*ent.Articles, error)
	Count(ctx context.Context, filters dto.ArticleFindAllFilters) (int, error)
	Delete(ctx context.Context, id int) error
}
