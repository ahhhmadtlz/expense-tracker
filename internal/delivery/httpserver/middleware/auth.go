package middleware

import (
	cfg "github.com/ahhhmadtlz/expense-tracker/internal/config"
	"github.com/ahhhmadtlz/expense-tracker/internal/domain/auth"
	mw "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
)

func Auth(service auth.Service, config auth.Config)echo.MiddlewareFunc {
	return mw.WithConfig(mw.Config{
		ContextKey:cfg.AuthMiddlewareContextKey ,
		SigningKey: []byte(config.SignKey),
		SigningMethod: "HS256",
		ParseTokenFunc: func(c echo.Context,auth string)(any,error){
			claims,err:=service.ParseBearerToken(auth)
			if err!=nil{
				return nil,err
			}
			return claims,nil
		},
	})
}