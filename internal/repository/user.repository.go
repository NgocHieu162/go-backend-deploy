package repository

import (
	"context"
	"go-backend/ent"
	"go-backend/ent/users"
	"go-backend/internal/common/pagination"
	"go-backend/internal/dto"
	"go-backend/internal/interfaces"
)

type UserRepository struct {
	entClient *ent.Client
}

func NewUserRepository(entClient *ent.Client) interfaces.UserRepository {
	return &UserRepository{
		entClient: entClient,
	}
}

// GetAll implements [interfaces.UserRepository].
func (a *UserRepository) GetAll(ctx context.Context, query pagination.Query, filters dto.UserFindAllFilters) ([]*ent.Users, error) {
	entQuery := a.entClient.Users.Query()

	if filters.Name != "" {
		entQuery = entQuery.Where(users.FullNameContainsFold(filters.Name))
	}

	entQuery = entQuery.Limit(query.PageSize)
	entQuery = entQuery.Offset(query.Offset)

	return entQuery.All(ctx)
}

// Count implements [interfaces.UserRepository].
func (a *UserRepository) Count(ctx context.Context, filters dto.UserFindAllFilters) (int, error) {
	entQuery := a.entClient.Users.Query()
	if filters.Name != "" {
		entQuery = entQuery.Where(users.FullNameContainsFold(filters.Name))
	}
	return entQuery.Count(ctx)
}

// ExistByEmail implements [interfaces.UserRepository].
func (a *UserRepository) ExistByEmail(ctx context.Context, email string) (bool, error) {
	entQuery := a.entClient.Users.Query()

	entQuery = entQuery.Where(users.EmailEQ(email))
	return entQuery.Exist(ctx)
}

// FindUserByEmail implements [interfaces.UserRepository].
func (a *UserRepository) FindUserByEmail(ctx context.Context, email string) (*ent.Users, error) {
	entQuery := a.entClient.Users.Query()
	entQuery = entQuery.Where(users.EmailEQ(email))

	return entQuery.Only(ctx)
}

// FindUserById implements [interfaces.UserRepository].
func (a *UserRepository) FindUserById(ctx context.Context, id int) (*ent.Users, error) {
	entQuery := a.entClient.Users.Query()
	entQuery = entQuery.Where(users.IDEQ(id))

	return entQuery.Only(ctx)
}

// CreateUser implements [interfaces.UserRepository].
func (a *UserRepository) CreateUser(ctx context.Context, body dto.AuthRegisterReq) (*ent.Users, error) {
	entCreate := a.entClient.Users.Create()

	entCreate = entCreate.SetEmail(body.Email)
	entCreate = entCreate.SetFullName(body.FullName)
	entCreate = entCreate.SetPassword(body.Password)

	return entCreate.Save(ctx)
}

// CreateGoogleUser implements [interfaces.UserRepository].
func (a *UserRepository) CreateGoogleUser(ctx context.Context, input dto.AuthGoogleRegisterReq) (*ent.Users, error) {
	entCreate := a.entClient.Users.Create()

	entCreate = entCreate.SetEmail(input.Email)
	entCreate = entCreate.SetFullName(input.FullName)
	entCreate = entCreate.SetAvatar(input.Avatar)
	entCreate = entCreate.SetGoogleID(input.GoogleId)

	return entCreate.Save(ctx)
}

// UploadAvatarByID implements [interfaces.UserRepository].
func (a *UserRepository) UpdateAvatarByID(ctx context.Context, id int, avatar string) (*ent.Users, error) {
	entUpdate := a.entClient.Users.UpdateOneID(id)
	entUpdate = entUpdate.SetAvatar(avatar)

	return entUpdate.Save(ctx)
}
