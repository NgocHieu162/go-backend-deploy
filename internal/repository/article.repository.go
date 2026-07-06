package repository

import (
	"context"
	"go-backend/ent"
	"go-backend/ent/articles"
	"go-backend/internal/common/pagination"
	"go-backend/internal/dto"
	"go-backend/internal/interfaces"
)

type articleRepository struct {
	entClient *ent.Client
}

func NewArticleRepository(entClient *ent.Client) interfaces.ArticleRepository {
	return &articleRepository{
		entClient: entClient,
	}
}

func (a *articleRepository) Create(ctx context.Context, body dto.ArticleCreateReq) (*ent.Articles, error) {
	entCreate := a.entClient.Articles.Create()

	entCreate = entCreate.SetUserID(1)
	entCreate = entCreate.SetTitle(body.Title)

	if body.Content != nil {
		entCreate = entCreate.SetContent(*body.Content)
	}

	if body.ImageUrl != nil {
		entCreate = entCreate.SetImageURL(*body.ImageUrl)
	}

	return entCreate.Save(ctx)
}

// func (a *articleRepository) CreateGorm(ctx context.Context, body dto.ArticleCreateReq) (any, error) {
// 	article := models.Articles{
// 		Title:    body.Title,
// 		Content:  body.Content,
// 		ImageUrl: body.ImageUrl,
// 	}

// 	result := gorm.WithResult()
// 	err := gorm.G[models.Articles](a.db, result).Create(ctx, &article)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return result, nil
// }

// FindAll implements [interfaces.ArticleRepository].
func (a *articleRepository) GetAll(ctx context.Context, query pagination.Query, filters dto.ArticleFindAllFilters) ([]*ent.Articles, error) {
	entQuery := a.entClient.Articles.Query()

	handlerFilter(entQuery, filters)

	entQuery = entQuery.WithUsers()
	entQuery = entQuery.Limit(query.PageSize)
	entQuery = entQuery.Offset(query.Offset)

	return entQuery.All(ctx)
}

// Count implements [interfaces.ArticleRepository].
func (a *articleRepository) Count(ctx context.Context, filters dto.ArticleFindAllFilters) (int, error) {
	entQuery := a.entClient.Articles.Query()
	handlerFilter(entQuery, filters)
	return entQuery.Count(ctx)
}

func handlerFilter(entQuery *ent.ArticlesQuery, filters dto.ArticleFindAllFilters) {
	if filters.Id > 0 {
		entQuery = entQuery.Where(articles.IDEQ(filters.Id))
	}

	if filters.Content != "" {
		entQuery = entQuery.Where(articles.ContentContainsFold(filters.Content))
	}

	if filters.Views != nil {
		entQuery = entQuery.Where(articles.ViewsEQ(*filters.Views))
	}
}

// Delete implements [interfaces.ArticleRepository].
func (a *articleRepository) Delete(ctx context.Context, id int) error {
	result := a.entClient.Articles.DeleteOneID(id)
	return result.Exec(ctx)
}