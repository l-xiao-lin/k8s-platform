package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var secret = []byte("cisco46589")

type MyClaims struct {
	UserID   int64
	Username string
	jwt.StandardClaims
}

// GetToken 生成token
func GetToken(userID int64, username string) (string, error) {

	c := MyClaims{
		UserID:   userID,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(2 * time.Hour).Unix(),
			Issuer:    "k8s-platform",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString(secret)
}

func ParseToken(tokenString string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("无效的token")
}
