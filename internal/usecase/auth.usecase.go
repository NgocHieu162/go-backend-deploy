package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"go-backend/ent"
	"go-backend/internal/common/env"
	"go-backend/internal/common/response"
	"go-backend/internal/dto"
	"go-backend/internal/interfaces"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"google.golang.org/api/idtoken"
)

type AuthUsecase struct {
	userRepository interfaces.UserRepository
	tokenUsecase   interfaces.TokenUsecase
	oauth2Config   *oauth2.Config
}

func NewAuthUsecase(userRepository interfaces.UserRepository, tokenUsecase interfaces.TokenUsecase, env *env.Env) interfaces.AuthUsecase {
	return &AuthUsecase{
		userRepository: userRepository,
		tokenUsecase:   tokenUsecase,
		oauth2Config: &oauth2.Config{
			ClientID:     env.GoogleClientID,
			ClientSecret: env.GoogleClientSecret,
			RedirectURL:  env.GoogleRedirectURL,
			Scopes: []string{
				"openid",
				"email",
				"profile",
			},
			Endpoint: google.Endpoint,
		},
	}
}

func (a *AuthUsecase) Register(ctx context.Context, body dto.AuthRegisterReq) (any, error) {
	// Kiểm tra user tồn tại hay chưa
	isExists, err := a.userRepository.ExistByEmail(ctx, body.Email)
	if err != nil {
		return nil, err
	}

	// Nếu tồn tại thì trả lỗi
	if isExists {
		return nil, response.NewBadRequestException("Already Register. Please Sign In.")
	}

	// hash password
	hashPass, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, response.NewBadRequestException(err.Error())
	}
	body.Password = string(hashPass)

	// Tạo người dùng mới
	userNew, err := a.userRepository.CreateUser(ctx, body)
	if err != nil {
		return nil, response.NewBadRequestException(err.Error())
	}

	return userNew, nil
}

func (a *AuthUsecase) Login(ctx context.Context, body dto.AuthLoginReq) (*dto.AuthLoginReturn, error) {
	user, err := a.userRepository.FindUserByEmail(ctx, body.Email)
	if err != nil {
		return nil, response.NewBadRequestException(err.Error())
	}

	if user.Password == nil {
		return nil, response.NewBadRequestException("Vui long dang nhap bang google de cap nhat password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(body.Password))
	if err != nil {
		fmt.Printf("Loi")
		return nil, response.NewBadRequestException("Mật khẩu không chính xác")
	}

	// trả về accessToken và refreshToken
	accessToken, err := a.tokenUsecase.CreateAccessToken(user.ID)
	if err != nil {
		return nil, response.NewBadRequestException(err.Error())
	}

	refreshToken, err := a.tokenUsecase.CreateRefreshToken(user.ID)
	if err != nil {
		return nil, response.NewBadRequestException(err.Error())
	}

	return &dto.AuthLoginReturn{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// GetInfo implements [interfaces.AuthUsecase].
func (a *AuthUsecase) GetInfo(ctx context.Context, user *ent.Users) (*ent.Users, error) {
	return user, nil
}

// RefreshToken implements [interfaces.AuthUsecase].
func (a *AuthUsecase) RefreshToken(ctx context.Context, accessToken string, refreshToken string) (*dto.AuthRefreshTokenReturn, error) {
	claimAccessToken, err := a.tokenUsecase.VerifyAccessToken(accessToken, jwt.WithoutClaimsValidation())
	if err != nil {
		return nil, response.NewUnauthorizedException(err.Error())
	}

	claimRefreshToken, err := a.tokenUsecase.VerifyRefreshToken(refreshToken)
	if err != nil {
		return nil, response.NewUnauthorizedException(err.Error())
	}

	if claimAccessToken.UserID != claimRefreshToken.UserID {
		return nil, response.NewUnauthorizedException("2 Token không cùng 1 user")
	}

	user, err := a.userRepository.FindUserById(ctx, claimAccessToken.UserID)
	if err != nil {
		return nil, response.NewUnauthorizedException(err.Error())
	}

	accessTokenNew, err := a.tokenUsecase.CreateAccessToken(user.ID)
	if err != nil {
		return nil, response.NewUnauthorizedException(err.Error())
	}

	refreshTokenNew, err := a.tokenUsecase.CreateRefreshToken(user.ID)
	if err != nil {
		return nil, response.NewUnauthorizedException(err.Error())
	}

	return &dto.AuthRefreshTokenReturn{
		AccessToken:  accessTokenNew,
		RefreshToken: refreshTokenNew,
	}, nil
}

// GoogleLogin implements [interfaces.AuthUsecase].
func (a *AuthUsecase) GoogleLogin(ctx context.Context) (*dto.AuthGoogleLoginReturn, error) {
	bytes := make([]byte, 10)
	// fmt.Println(bytes)

	rand.Read(bytes)
	// fmt.Println(bytes)

	state := hex.EncodeToString(bytes)
	// fmt.Println(state)

	url := a.oauth2Config.AuthCodeURL(state)
	// fmt.Println("url", url)

	return &dto.AuthGoogleLoginReturn{
		State: state,
		Url:   url,
	}, nil
}

// GoogleCallback implements [interfaces.AuthUsecase].
func (a *AuthUsecase) GoogleCallback(ctx context.Context, input *dto.AuthGoogleCallbackInput) (*dto.AuthLoginReturn, error) {
	// fmt.Println(input.Code)
	// fmt.Println(input.StateCookie)
	// fmt.Println(input.StateGoogle)

	if input.StateCookie != input.StateGoogle {
		return nil, errors.New("State Invalid")
	}

	token, err := a.oauth2Config.Exchange(ctx, input.Code)
	if err != nil {
		return nil, err
	}

	id_token := token.Extra("id_token").(string)
	idTokenPayload, err := idtoken.Validate(ctx, id_token, a.oauth2Config.ClientID)
	if err != nil {
		return nil, err
	}

	email := idTokenPayload.Claims["email"].(string)
	emailVerified := idTokenPayload.Claims["email_verified"].(bool)
	name := idTokenPayload.Claims["name"].(string)
	picture := idTokenPayload.Claims["picture"].(string)
	googleId := idTokenPayload.Claims["sub"].(string)

	if !emailVerified {
		return nil, errors.New("Email not verify")
	}

	userExists, err := a.userRepository.FindUserByEmail(ctx, email)
	if err != nil && !ent.IsNotFound(err) {
		return nil, err
	}

	if userExists == nil {
		input := dto.AuthGoogleRegisterReq{
			Email:    email,
			FullName: name,
			Avatar:   picture,
			GoogleId: googleId,
		}
		userExists, err = a.userRepository.CreateGoogleUser(ctx, input)
		if err != nil {
			return nil, err
		}
	}

	accessToken, err := a.tokenUsecase.CreateAccessToken(userExists.ID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := a.tokenUsecase.CreateRefreshToken(userExists.ID)
	if err != nil {
		return nil, err
	}

	return &dto.AuthLoginReturn{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
