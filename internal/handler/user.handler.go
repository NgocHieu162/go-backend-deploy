package handler

import (
	"encoding/json"
	"fmt"
	"go-backend/internal/common/helpers"
	"go-backend/internal/common/pagination"
	"go-backend/internal/common/response"
	"go-backend/internal/dto"
	"go-backend/internal/interfaces"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUsecase interfaces.UserUsecase
}

func NewUserHandler(UserUsecase interfaces.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: UserUsecase,
	}
}

func (a *UserHandler) FindAll(ctx *gin.Context) {
	queryPagination := pagination.Get(ctx.Query("page"), ctx.Query("pageSize"))

	filterString := ctx.DefaultQuery("filters", "{}")
	var filters dto.UserFindAllFilters
	json.Unmarshal([]byte(filterString), &filters)
	fmt.Printf("%+v \n\n", filters)

	input := dto.UserFindAllInput{
		Query:              *queryPagination,
		UserFindAllFilters: filters,
	}
	fmt.Printf("%+v \n\n", queryPagination)

	result, err := a.userUsecase.FindAll(ctx.Request.Context(), input)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.Success(result, "", 0, ctx)
}

func (a *UserHandler) FindOne(ctx *gin.Context) {
	idString := ctx.Param("id")
	if idString == "" {
		ctx.Error(response.NewBadRequestException("required id"))
		return
	}

	id, err := strconv.Atoi(idString)
	if err != nil {
		ctx.Error(response.NewBadRequestException("id khong hop le"))
		return
	}

	result, err := a.userUsecase.FindOne(ctx.Request.Context(), id)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.Success(result, "", 0, ctx)
}

func (a *UserHandler) UploadAvatarLocal(ctx *gin.Context) {
	user, err := helpers.GetUsers(ctx)

	if err != nil {
		ctx.Error(response.NewBadRequestException(err.Error()))
		return
	}

	fileHeader, err := ctx.FormFile("avatar")

	if err != nil {
		ctx.Error(response.NewBadRequestException(err.Error()))
		return
	}

	result, err := a.userUsecase.UploadAvatarLocal(ctx.Request.Context(), fileHeader, user)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.Success(result, "", 0, ctx)
}

func (a *UserHandler) UploadAvatarCloud(ctx *gin.Context) {
	user, err := helpers.GetUsers(ctx)
	if err != nil {
		ctx.Error(response.NewBadRequestException(err.Error()))
		return
	}

	fileHeader, err := ctx.FormFile("avatar")
	if err != nil {
		ctx.Error(response.NewBadRequestException(err.Error()))
		return
	}

	result, err := a.userUsecase.UploadAvatarCloud(ctx.Request.Context(), fileHeader, user)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.Success(result, "", 0, ctx)
}
