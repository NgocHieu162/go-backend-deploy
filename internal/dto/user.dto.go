package dto

import (
	"go-backend/internal/common/pagination"
	"mime/multipart"
)

type UploadInput struct {
	Prefix     string
	FileHeader *multipart.FileHeader
	Folder     string
}

type UploadReturn struct {
	FileName string
	FullPath string
	Url      string
}

type UserFindAllFilters struct {
	Name string `json:"name"`
}

type UserFindAllInput struct {
	Query pagination.Query
	UserFindAllFilters
}
