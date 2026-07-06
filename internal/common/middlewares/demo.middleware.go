package middlewares

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type User struct {
	Id int
	Email string 
	Pass string
}

func A(ctx *gin.Context) {
	fmt.Println("Mid A Before")
	user := User{
		Id: 1,
		Email: "hello",
		Pass: "123",
	}
	ctx.Set("user", user)
	ctx.Next()
	fmt.Println("Mid A After")

}
func B(ctx *gin.Context) {
	fmt.Println("Mid B Before")
	ctx.Next()
	fmt.Println("Mid B After")
}
func C(ctx *gin.Context) {
	fmt.Println("Mid C Before")
	ctx.Next()
	fmt.Println("Mid C After")
}