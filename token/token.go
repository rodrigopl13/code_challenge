package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"

	"jobsity-code-challenge/entities"
)

const defaultExpirationTime = time.Hour

type Claims struct {
	UserName string `json:"user_name"`
	jwt.StandardClaims
}

type TokenService struct {
	hmacSecret string
}

func New(secret string) *TokenService {
	return &TokenService{hmacSecret: secret}
}

func (t TokenService) GenerateToken(user entities.User) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_name": user.UserName,
		"exp":       time.Now().Add(defaultExpirationTime).Unix(),
	})
	return token.SignedString([]byte(t.hmacSecret))
}

func (t TokenService) ValidateToken(tokenString string) error {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(fmt.Sprintf("Invalid signing method: %s", token.Header["alg"]))
		}
		return []byte(t.hmacSecret), nil
	})
	if _, ok := token.Claims.(*Claims); ok && token.Valid {
		return nil
	}
	return err
}
