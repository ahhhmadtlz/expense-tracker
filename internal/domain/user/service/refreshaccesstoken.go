package service

import (
	"context"

	"github.com/ahhhmadtlz/expense-tracker/internal/domain/user/param"
	"github.com/ahhhmadtlz/expense-tracker/internal/pkg/richerror"
)

func (s Service) RefreshAccessToken(ctx context.Context, req param.RefreshAccessTokenRequest)(param.RefreshAccessTokenResponse,error){
	const op=richerror.Op("userService.RefreshAccessToken")
	claims,err:=s.auth.ParseRefreshToken(req.RefreshToken)


	if err != nil {
    return param.RefreshAccessTokenResponse{}, richerror.New(op).
        WithMessage("invalid or expired refresh token").
        WithKind(richerror.KindInvalid).
        WithErr(err)
	}
	
	user,err:=s.repo.GetUserByID(ctx,claims.UserID)

	if err!=nil{
		return  param.RefreshAccessTokenResponse{},richerror.New(op).WithMessage("failed to retrieve user").
		WithKind(richerror.KindUnexpected).
		WithErr(err)
	}

	accessToken,err:=s.auth.CreateAccessToken(user)
	if err !=nil {
		return param.RefreshAccessTokenResponse{},richerror.New(op).
    WithMessage("failed to create access token").
    WithKind(richerror.KindUnexpected).
    WithErr(err)
	}
	return param.RefreshAccessTokenResponse{
		AccessToken: accessToken,
	},nil

}