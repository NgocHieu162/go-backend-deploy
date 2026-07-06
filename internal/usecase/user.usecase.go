package usecase

import (
	"context"
	"fmt"
	"go-backend/ent"
	"go-backend/internal/common/pagination"
	"go-backend/internal/common/response"
	"go-backend/internal/common/storage"
	"go-backend/internal/dto"
	"go-backend/internal/interfaces"
	"math"
	"mime/multipart"
	"path/filepath"
)

type UserUsecase struct {
	userRepository    interfaces.UserRepository
	localFileStorage  storage.FileStorage
	cloudinaryStorage storage.FileStorage
}

func NewUserUsecase(userRepository interfaces.UserRepository, localFileStorage storage.FileStorage, cloudinaryStorage storage.FileStorage) interfaces.UserUsecase {
	return &UserUsecase{
		userRepository:    userRepository,
		localFileStorage:  localFileStorage,
		cloudinaryStorage: cloudinaryStorage,
	}
}

// FindAll implements [usecase.UserUsecase].
func (a *UserUsecase) FindAll(ctx context.Context, input dto.UserFindAllInput) (any, error) {
	data, err := a.userRepository.GetAll(ctx, input.Query, input.UserFindAllFilters)

	if err != nil {
		return nil, response.NewBadRequestException(err.Error())
	}

	// totalItem: tổng số lượng item
	totalItem, err := a.userRepository.Count(ctx, input.UserFindAllFilters)
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

// FindAll implements [usecase.UserUsecase].
func (a *UserUsecase) FindOne(ctx context.Context, id int) (any, error) {
	return a.userRepository.FindUserById(ctx, id)
}

// UploadAvatarLocal implements [interfaces.UserUsecase].
func (a *UserUsecase) UploadAvatarLocal(ctx context.Context, fileHeader *multipart.FileHeader, user *ent.Users) (any, error) {
	input := dto.UploadInput{
		Prefix:     "local-",
		FileHeader: fileHeader,
		Folder:     "images",
	}

	uploadResult, err := a.localFileStorage.Upload(input)
	if err != nil {
		return nil, response.NewBadRequestException(err.Error())
	}

	userUpdate, err := a.userRepository.UpdateAvatarByID(ctx, user.ID, uploadResult.FileName)
	if err != nil {
		err := a.localFileStorage.Delete(filepath.Join(input.Folder, uploadResult.FileName))
		if err != nil {
			fmt.Println("LocalFileStorage.Delete", err.Error())
		}

		err = a.cloudinaryStorage.Delete(uploadResult.FullPath)
		if err != nil {
			fmt.Println("CloudinaryStorage.Delete", err.Error())
		}

		return nil, response.NewBadRequestException(err.Error())
	}

	// Nếu trước đó user có avatar thì xóa avatar cũ
	if user.Avatar != nil {
		// Xóa hình cũ
		err := a.localFileStorage.Delete(filepath.Join(input.Folder, *user.Avatar))
		if err != nil {
			fmt.Println("LocalFileStorage.Delete", err.Error())
		}

		err = a.cloudinaryStorage.Delete(*user.Avatar)
		if err != nil {
			fmt.Println("CloudinaryStorage.Delete", err.Error())
		}
	}

	return userUpdate, nil
}

// UploadAvatarCloud implements [interfaces.UserUsecase].
func (a *UserUsecase) UploadAvatarCloud(ctx context.Context, fileHeader *multipart.FileHeader, user *ent.Users) (any, error) {
	input := dto.UploadInput{
		Prefix:     "cloud-",
		FileHeader: fileHeader,
		Folder:     "images",
	}

	uploadResult, err := a.cloudinaryStorage.Upload(input)
	if err != nil {
		return nil, response.NewBadRequestException(err.Error())
	}

	userUpdate, err := a.userRepository.UpdateAvatarByID(ctx, user.ID, uploadResult.FullPath)
	if err != nil {
		err := a.cloudinaryStorage.Delete(uploadResult.FullPath)
		if err != nil {
			fmt.Println("CloudinaryStorage.Delete", err.Error())
		}

		err = a.localFileStorage.Delete(filepath.Join(input.Folder, uploadResult.FileName))
		if err != nil {
			fmt.Println("LocalFileStorage.Delete", err.Error())
		}

		return nil, response.NewBadRequestException(err.Error())
	}

	// Nếu trước đó user có avatar thì xóa avatar cũ
	if user.Avatar != nil {
		// Xóa hình cũ
		err := a.cloudinaryStorage.Delete(*user.Avatar)
		if err != nil {
			fmt.Println("CloudinaryStorage.Delete", err.Error())
		}

		err = a.localFileStorage.Delete(filepath.Join(input.Folder, *user.Avatar))
		if err != nil {
			fmt.Println("LocalFileStorage.Delete", err.Error())
		}
	}

	return userUpdate, nil
}
