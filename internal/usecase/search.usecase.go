package usecase

import (
	"context"
	"go-backend/internal/interfaces"
)

type SearchUsecase struct {
	searchRepository interfaces.SearchRepository
}

func NewSearchUsecase(searchRepository interfaces.SearchRepository) interfaces.SearchUsecase {
	return &SearchUsecase{
		searchRepository: searchRepository,
	}
}

// FindAll implements [usecase.SearchUsecase].
func (a *SearchUsecase) FindAll(ctx context.Context, textSearch string) (any, error) {
	return a.searchRepository.FindAll(ctx, textSearch)
}
