package middlewares

import (
	"errors"
	"go-backend/internal/common/response"
	"go-backend/internal/interfaces"
	"net/http"

	"github.com/golang-jwt/jwt/v5"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	tokenUsecase   interfaces.TokenUsecase
	userRepository interfaces.UserRepository
}

func NewAuthMiddleware(tokenUsecase interfaces.TokenUsecase, userRepository interfaces.UserRepository) *AuthMiddleware {
	return &AuthMiddleware{
		tokenUsecase:   tokenUsecase,
		userRepository: userRepository,
	}
}

func (a *AuthMiddleware) Protect(ctx *gin.Context) {
	accessToken, err := ctx.Cookie("accessToken")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			ctx.Error(response.NewUnauthorizedException(""))
			ctx.Abort()
			return
		}
		ctx.Error(response.NewBadRequestException(""))
		ctx.Abort()
		return
	}

	claim, err := a.tokenUsecase.VerifyAccessToken(accessToken)
	if err != nil {
		// fmt.Println(err)
		if errors.Is(err, jwt.ErrTokenExpired) {
			ctx.Error(response.NewForbiddenException(""))
			ctx.Abort()
			return
		}
		ctx.Error(response.NewUnauthorizedException(""))
		ctx.Abort()
		return
	}

	user, err := a.userRepository.FindUserById(ctx, claim.UserID)
	if err != nil {
		ctx.Error(response.NewUnauthorizedException(""))
		ctx.Abort()
		return
	}

	ctx.Set("user", user)

	// fmt.Println("user", user)
}
