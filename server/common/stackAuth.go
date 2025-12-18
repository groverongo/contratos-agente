package common

import (
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func StackAuthValidation(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		auth := c.Request().Header.Get("Authorization")
		if auth != "" && strings.HasPrefix(auth, "Bearer ") {
			tokenString := strings.TrimPrefix(auth, "Bearer ")

			// Decode JWT
			token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &StackAuthClaims{})
			if err == nil {
				if claims, ok := token.Claims.(*StackAuthClaims); ok {
					sub := claims.Subject
					// Set user_id in context
					c.Set("user_id", sub)
					// Also set the header that handlers currently expect
					c.Request().Header.Set("X-Stack-Auth-User-Id", sub)
				}
			}
		}
		return next(c)
	}
}
