package jwtoken

import (
	"github.com/BatuhanSerin/final-project/package/config"
	"github.com/golang-jwt/jwt"
)

type JWToken struct {
	Secret string
}

var cfg *config.Config

func NewJWToken(secret string) *JWToken {
	return &JWToken{Secret: secret}
}

func GenerateToken(claims *jwt.Token, secret string) (string, error) {
	hmacSecret := []byte(secret)
	token, err := claims.SignedString(hmacSecret)
	if err != nil {
		return "", err
	}

	return token, nil
}
func ParseToken(tokenString string, secret string) (*jwt.MapClaims, error) {
	hmacSecret := []byte(secret)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return hmacSecret, nil
	})
	if err != nil {
		return nil, err
	}
	decodedClaims := token.Claims.(jwt.MapClaims)

	if token.Valid {

		return &decodedClaims, nil
	}
	return nil, nil
}
