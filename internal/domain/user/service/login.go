package service

import (
	"context"

	"github.com/ahhhmadtlz/expense-tracker/internal/domain/user/param"
	"github.com/ahhhmadtlz/expense-tracker/internal/observability/logger"
	"github.com/ahhhmadtlz/expense-tracker/internal/pkg/richerror"
	"golang.org/x/crypto/bcrypt"
)


func (s Service)Login(ctx context.Context ,req param.LoginRequest)(param.LoginResponse,error){
	const op=richerror.Op("userservice.Login")

	// Log login attempt
	logger.Info("Login attempt",
		"email", req.PhoneNumber,
	)

	user,err:=s.repo.GetUserByPhoneNumber(ctx,req.PhoneNumber)
	if err!=nil{
			logger.Warn("Login failed - user not found",
			"phoneNumber", req.PhoneNumber,
		)
		return param.LoginResponse{},richerror.New(op).WithErr(err).WithMeta("phone_number",req.PhoneNumber)
	}

	err =bcrypt.CompareHashAndPassword([]byte(user.Password),[]byte(req.Password))

	if err!=nil{
			logger.Warn("Login failed - user or password is incorrect",
				"user_id", user.ID,
				"phoneNumber", user.PhoneNumber,
			)
		return param.LoginResponse{},richerror.New(op).
			WithErr(err).
			WithMessage("username or password is incorrect")
	}

	accessToken,err:=s.auth.CreateAccessToken(user)

	if err!=nil{
		logger.Error("Failed to create access token",
			"user_id", user.ID,
			"error", err.Error(),
		)
		return  param.LoginResponse{},richerror.New(op).WithErr(err)
	}

	refreshToken,err:=s.auth.CreateRefreshToken(user)

	if err!=nil{
		logger.Error("Failed to create refresh token",
			"user_id", user.ID,
			"error", err.Error(),
		)
		return  param.LoginResponse{},richerror.New(op).WithErr(err)
	}

	// Log successful login
	logger.Info("User logged in successfully",
		"user_id", user.ID,
		"phoneNumber", user.PhoneNumber,
	)
	return param.LoginResponse{
		User:param.UserInfo{
			ID:user.ID,
			PhoneNumber: user.PhoneNumber,
			Name: user.Name,
		},
		Tokens:param.Tokens{
			AccessToken: accessToken,
			RefreshToken: refreshToken,
		},
	},nil
}