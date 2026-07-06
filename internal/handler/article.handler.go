package handler

import (
	"encoding/json"
	"fmt"
	"go-backend/internal/common/pagination"
	"go-backend/internal/common/response"
	"go-backend/internal/dto"
	"go-backend/internal/interfaces"
	"io"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ArticleHandler struct {
	articleUsecase interfaces.ArticleUsecase
}

func NewArticleHandler(	articleUsecase interfaces.ArticleUsecase) *ArticleHandler {
	return &ArticleHandler{
		articleUsecase: articleUsecase,
	}
}

func (a *ArticleHandler) Create(ctx *gin.Context) {
	var body dto.ArticleCreateReq
	err := ctx.ShouldBindJSON(&body)
	if err != nil{
		if err == io.EOF{
			ctx.Error(response.NewBadRequestException("Body required"))
			return
		}
		ctx.Error(response.NewBadRequestException(err.Error()))
		return 
	}

	result, err := a.articleUsecase.Create(ctx.Request.Context(), body)
	if err != nil{
		ctx.Error(err)
		return
	}

	response.Success(result, "", 0, ctx)
}

func (a *ArticleHandler) FindAll(ctx *gin.Context){
	queryPagination := pagination.Get(ctx.Query("page"), ctx.Query("pageSize"))

	filterString := ctx.DefaultQuery("filters", "{}")
	var filters dto.ArticleFindAllFilters
	json.Unmarshal([]byte(filterString), &filters)
	fmt.Printf("%+v \n\n", filters)

	input := dto.ArticleFindAllInput{
		Query: *queryPagination,
		ArticleFindAllFilters: filters,
	}
	fmt.Printf("%+v \n\n", queryPagination)

	result, err := a.articleUsecase.FindAll(ctx.Request.Context(), input)
	if err != nil{
		ctx.Error(err)
		return
	}

	response.Success(result, "", 0, ctx)
}

func (a* ArticleHandler)Delete(ctx *gin.Context){
	idString := ctx.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil{
		ctx.Error(response.NewBadRequestException(err.Error()))
		return 
	}

	result, err := a.articleUsecase.Delete(ctx.Request.Context(), id)
	if err != nil{
		ctx.Error(err)
		return
	}

	response.Success(result, "", 0, ctx)
}