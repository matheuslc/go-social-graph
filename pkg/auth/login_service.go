package auth

import (
	"context"
	"fmt"
	"gosocialgraph/pkg/repository"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

//go:generate mockgen -source=./login_service.go -destination=../mock/auth/login_service.go

var (
	ErrWrongUsernameOrPassword = fmt.Errorf("Wrong username or password")
)

type JwtCustomClaims struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`

	jwt.StandardClaims
}

type AuthService struct {
	Repository repository.UserReader
}

func (sv AuthService) Run(ctx context.Context, username string, password string) (string, string, error) {
	user, err := sv.Repository.FindByUsername(ctx, username)
	if err != nil {
		return "", "", nil
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(fmt.Sprintf("%s%s", repository.Salt, password))); err != nil {
		return "", "", ErrWrongUsernameOrPassword
	}

	claims := &JwtCustomClaims{
		user.ID.String(),
		user.Username,
		"user",
		jwt.StandardClaims{
			Subject:   user.ID.String(),
			Issuer:    "api",
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(JwtScret))
	if err != nil {
		return "", "", err
	}

	refreshClaims := &jwt.StandardClaims{
		Subject:   user.ID.String(),
		Issuer:    "api",
		ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	rt, err := refreshToken.SignedString([]byte(JwtScret))
	if err != nil {
		return "", "", err
	}

	return t, rt, err
}
