package service

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"github.com/ahhhmadtlz/expense-tracker/internal/domain/user/entity"
	"github.com/ahhhmadtlz/expense-tracker/internal/domain/user/param"
	"github.com/ahhhmadtlz/expense-tracker/internal/pkg/richerror"
)

func (s Service) Register(ctx context.Context, req param.RegisterRequest)(param.RegisterResponse,error){
	const op=richerror.Op("userService.Register")

	s.logger.Info("User registration attempt",
		"phonenumber",req.PhoneNumber,
		"name",req.Name,
	)

	hashedPassword,err:=bcrypt.GenerateFromPassword([]byte(req.Password),bcrypt.DefaultCost)

	if err !=nil{
		s.logger.Error("Failed to hash password","error",err.Error())
		return param.RegisterResponse{},richerror.New(op).WithMessage("failed to hash password").WithKind(richerror.KindUnexpected).WithErr(err)
	}

	user:=entity.User{
		Name: req.Name,
		PhoneNumber: req.PhoneNumber,
		Password: string(hashedPassword),
		Role: entity.UserRole,
	}

	createdUser,err:=s.repo.RegisterUser(ctx,user)
	if err!=nil{
		s.logger.Error("Failed to create user",
			"phoneNumber", req.PhoneNumber,
			"error", err.Error(),
		)
		return  param.RegisterResponse{},richerror.New(op).WithMessage("failed to register user").WithKind(richerror.KindUnexpected).WithErr(err)
	}

	//  successful registration
	s.logger.Info("User registered successfully",
		"user_id", createdUser.ID,
		"phoneNumber", createdUser.PhoneNumber,
	)
	return  param.RegisterResponse{
		User:param.UserInfo{
			ID:createdUser.ID,
			Name:createdUser.Name,
			PhoneNumber: createdUser.PhoneNumber,
		},
	},nil

}