package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/golang-jwt/jwt"
)

var (
	ErrNoAuthHeader      = errors.New("Authorization header is missing")
	ErrInvalidAuthHeader = errors.New("Authorization header is malformed")
	ErrClaimsInvalid     = errors.New("Provided claims do not match expected scopes")

	JwtScret = os.Getenv("JWT_SECRET")
	JwtSalt  = os.Getenv("JWT_SALT")
)

func NewAuthenticator() openapi3filter.AuthenticationFunc {
	return func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
		if err := Authenticate(ctx, input); err != nil {
			return err
		}

		return nil
	}
}

func Authenticate(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
	if input.SecuritySchemeName != "bearerAuth" {
		return fmt.Errorf("security scheme %s != 'BearerAuth'", input.SecuritySchemeName)
	}

	jws, err := GetJWSFromRequest(input.RequestValidationInput.Request)
	if err != nil {
		return fmt.Errorf("getting jws: %w", err)
	}

	_, err = ParseJWT(jws)
	if err != nil {
		return fmt.Errorf("parsing jwt: %w", err)
	}

	return nil
}

func ParseJWT(token string) (*jwt.Token, error) {
	fn := jwt.Keyfunc(func(token *jwt.Token) (interface{}, error) {
		return []byte(JwtScret), nil
	})

	parsed, err := jwt.Parse(token, fn)
	if err != nil || !parsed.Valid {
		return nil, err
	}

	return parsed, nil
}

func GetJWSFromRequest(req *http.Request) (string, error) {
	authHdr := req.Header.Get("Authorization")
	// Check for the Authorization header.
	if authHdr == "" {
		return "", ErrNoAuthHeader
	}
	// We expect a header value of the form "Bearer <token>", with 1 space after
	// Bearer, per spec.
	prefix := "Bearer "
	if !strings.HasPrefix(authHdr, prefix) {
		return "", ErrInvalidAuthHeader
	}
	return strings.TrimPrefix(authHdr, prefix), nil
}
