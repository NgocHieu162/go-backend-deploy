package dto

import "go-backend/internal/common/pagination"

type ArticleCreateReq struct {
	Title    string  `json:"title" binding:"required"`
	Content  *string `json:"content" binding:"omitempty"`
	ImageUrl *string `json:"image_url" binding:"omitempty"`
}

type ArticleFindAllFilters struct {
	Id      int    `json:"id"`
	Content string `json:"content"`
	Views   *int   `json:"views"`
}

type ArticleFindAllInput struct {
	Query pagination.Query
	ArticleFindAllFilters
}

type ArticleFindAllReq struct{
	pagination.Query
	Filters string `query:"filters" example:"{ \"id\":1, \"content\":\"string\", \"views\":55 }"`
}