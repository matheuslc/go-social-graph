package auth

import (
	"fmt"
	"gosocialgraph/pkg/repository"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

//go:generate mockgen -source=./login_service.go -destination=../mock/auth/login_service.go

type JwtCustomClaims struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
	jwt.StandardClaims
}

type AuthService struct {
	Repository repository.UserReader
}

func (sv AuthService) Run(username string, password string) (string, string, error) {
	user, err := sv.Repository.FindByUsername(username)
	if err != nil {
		return "", "", nil
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(fmt.Sprintf("%s%s", repository.Salt, password))); err != nil {
		return "", "", nil
	}

	claims := &JwtCustomClaims{
		user.ID.String(),
		user.Username,
		"admin",
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Microsecond + 1).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", "", err
	}

	refreshClaims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	rt, err := refreshToken.SignedString([]byte("secret"))
	if err != nil {
		return "", "", err
	}

	return t, rt, err
}
