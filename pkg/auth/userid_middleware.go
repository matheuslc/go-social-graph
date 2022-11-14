package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func UserIDAtContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, err := GetJWSFromRequest(c.Request())
		if err != nil {
			fmt.Println("not a JWT request")
			return next(c)
		}

		parsed, err := ParseJWT(token)
		if err != nil {
			return err
		}

		claims, ok := parsed.Claims.(jwt.MapClaims)
		if !ok {
			return fmt.Errorf("could not parse jwt claims")
		}

		c.Set("userID", claims["id"])
		return next(c)
	}
}
