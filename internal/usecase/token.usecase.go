package usecase

import (
	"errors"
	"go-backend/internal/common/env"
	"go-backend/internal/dto"
	"go-backend/internal/interfaces"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type tokenUsecase struct {
	env *env.Env
}

func NewTokenUsecase(env *env.Env) interfaces.TokenUsecase {
	return &tokenUsecase{
		env: env,
	}
}

// CreateAccessToken implements [interfaces.TokenUsecase].
func (t *tokenUsecase) CreateAccessToken(userID int) (string, error) {
	claim := dto.CustomClaim{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(t.env.AccessTokenExpiresAt)),
		},
	}

	tokenClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return tokenClaim.SignedString([]byte(t.env.AccessTokenSecret))
}

// VerifyAccessToken implements [interfaces.TokenUsecase].
func (t *tokenUsecase) VerifyAccessToken(accessToken string, option ...jwt.ParserOption) (*dto.CustomClaim, error) {
	claim := &dto.CustomClaim{}
	tokenClaim, err := jwt.ParseWithClaims(
		accessToken,
		claim,
		func(*jwt.Token) (any, error) {
			return []byte(t.env.AccessTokenSecret), nil
		},
		option...,
	)

	if err != nil {
		return nil, err
	}

	if !tokenClaim.Valid {
		return nil, errors.New("Access Token Invalid")
	}

	return claim, nil
}

// CreateRefreshToken implements [interfaces.TokenUsecase].
func (t *tokenUsecase) CreateRefreshToken(userID int) (string, error) {
	claim := dto.CustomClaim{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(t.env.RefreshTokenExpiresAt)),
		},
	}

	tokenClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return tokenClaim.SignedString([]byte(t.env.RefreshTokenSecret))
}

// VerifyRefreshToken implements [interfaces.TokenUsecase].
func (t *tokenUsecase) VerifyRefreshToken(refreshToken string, option ...jwt.ParserOption) (*dto.CustomClaim, error) {
	claim := &dto.CustomClaim{}
	tokenClaim, err := jwt.ParseWithClaims(
		refreshToken,
		claim,
		func(*jwt.Token) (any, error) {
			return []byte(t.env.RefreshTokenSecret), nil
		},
		option...,
	)

	if err != nil {
		return nil, err
	}

	if !tokenClaim.Valid {
		return nil, errors.New("Refresb Token Invalid")
	}

	return claim, nil
}