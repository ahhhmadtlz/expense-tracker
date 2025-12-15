package service

import (
	"context"

	"github.com/ahhhmadtlz/expense-tracker/internal/domain/user/param"
	"github.com/ahhhmadtlz/expense-tracker/internal/pkg/richerror"
)

func (s Service) RefreshAccessToken(ctx context.Context, req param.RefreshAccessTokenRequest)(param.RefreshAccessTokenResponse,error){
	const op=richerror.Op("userService.RefreshAccessToken")

	// Log refresh attempt
	s.logger.Info("Refresh token request received")

	claims,err:=s.auth.ParseRefreshToken(req.RefreshToken)


	if err != nil {
		s.logger.Warn("Invalid refresh token attempt",
			"error", err.Error(),
		)
    return param.RefreshAccessTokenResponse{}, richerror.New(op).
        WithMessage("invalid or expired refresh token").
        WithKind(richerror.KindInvalid).
        WithErr(err)
	}
	
	user,err:=s.repo.GetUserByID(ctx,claims.UserID)

	if err!=nil{
		s.logger.Error("Failed to retrieve user",
			"user_id", claims.UserID,
			"error", err.Error(),
		)

		return  param.RefreshAccessTokenResponse{},richerror.New(op).WithMessage("failed to retrieve user").
		WithKind(richerror.KindUnexpected).
		WithErr(err)
	}

	accessToken,err:=s.auth.CreateAccessToken(user)
	if err !=nil {
		s.logger.Error("Failed to create access token",
			"user_id", user.ID,
			"error", err.Error(),
		)
		return param.RefreshAccessTokenResponse{},richerror.New(op).
    WithMessage("failed to create access token").
    WithKind(richerror.KindUnexpected).
    WithErr(err)
	}

	// Log successful token refresh
	s.logger.Info("Access token refreshed successfully",
		"user_id", user.ID,
	)
	
	return param.RefreshAccessTokenResponse{
		AccessToken: accessToken,
	},nil

}