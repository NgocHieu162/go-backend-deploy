package interfaces

import (
	"context"
	"go-backend/ent"
	"go-backend/internal/common/pagination"
	"go-backend/internal/dto"
	"mime/multipart"
)

type UserUsecase interface {
	FindAll(ctx context.Context, input dto.UserFindAllInput) (any, error)
	FindOne(ctx context.Context, id int) (any, error)
	UploadAvatarLocal(ctx context.Context, fileHeader *multipart.FileHeader, user *ent.Users) (any, error)
	UploadAvatarCloud(ctx context.Context, fileHeader *multipart.FileHeader, user *ent.Users) (any, error)
}

type UserRepository interface {
	GetAll(ctx context.Context, query pagination.Query, filters dto.UserFindAllFilters) ([]*ent.Users, error)
	Count(ctx context.Context, filters dto.UserFindAllFilters) (int, error)
	ExistByEmail(ctx context.Context, email string) (bool, error)
	FindUserByEmail(ctx context.Context, email string) (*ent.Users, error)
	FindUserById(ctx context.Context, id int) (*ent.Users, error)
	CreateUser(ctx context.Context, body dto.AuthRegisterReq) (*ent.Users, error)
	CreateGoogleUser(ctx context.Context, body dto.AuthGoogleRegisterReq) (*ent.Users, error)
	UpdateAvatarByID(ctx context.Context, id int, avatar string) (*ent.Users, error)
}
