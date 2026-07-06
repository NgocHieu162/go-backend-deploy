package handler

import (
	"fmt"
	"go-backend/internal/common/middlewares"
	"go-backend/internal/common/pagination"
	"go-backend/internal/common/response"
	"go-backend/internal/dto"
	"go-backend/internal/interfaces"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DemoHandler struct {
	demoUsecase interfaces.DemoUsecase
}

func NewDemoHandler(demoInterface interfaces.DemoUsecase) *DemoHandler {
	return &DemoHandler{
		demoUsecase: demoInterface,
	}
}

func (d *DemoHandler) Query(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	pageSize, _ := strconv.Atoi(ctx.Query("pageSize"))

	input := &pagination.Query {
		Page: page,
		PageSize: pageSize,
	}

	data := d.demoUsecase.Query(input)
	response.Success(data, "", 0, ctx)
}

func (d *DemoHandler)Param (ctx *gin.Context){
	fmt.Println("Handler Param (Mid)")

	userRaw, exists := ctx.Get("user")
	if !exists{
		fmt.Println("Loi ko co user")
	}
	user, ok := userRaw.(middlewares.User)

	if !ok{
		fmt.Println("Loi ko phai struct User")
	}

	fmt.Println("Handler nhan duoc", user)
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil{
		ctx.Error(response.NewBadRequestException())
		return
	}

	data := d.demoUsecase.Param(id)
	response.Success(data, "", 0, ctx)
}

func (d *DemoHandler)Body (ctx *gin.Context){
	var body dto.DemoBody
	err := ctx.ShouldBindJSON(&body)
	fmt.Println(err)

	data := d.demoUsecase.Body(&body)
	response.Success(data, "", 0, ctx)
}

func (d *DemoHandler)Header (ctx *gin.Context){
	apiKey := ctx.GetHeader("API_KEY")
	data := d.demoUsecase.Header(apiKey)
	response.Success(data, "", 0, ctx)
}