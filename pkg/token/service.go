package tokenHelpers

import (
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)


var ErrExpiredToken = errors.New("token has expired")
var ErrInvalidToken = errors.New("invalid token")
var ErrTokenInternalServerError = errors.New("internal server error")


const SecretKey = "ThisIsMySecretKey"


type JWTClaims struct {
	jwt.StandardClaims
	Name string `json:"name"`
}


func GenerateToken(name string, id int64) (signedToken string, err error) {

	claims := &JWTClaims{
		Name: name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour).Unix(),
			Id:        strconv.Itoa(int(id)),
		},
	}

	signedToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SecretKey))

	if err != nil {
		log.Panic(err)
		return
	}

	return signedToken, nil
}


func ValidateToken(signedToken string) (*JWTClaims, error) {

	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SecretKey), nil
		},
	)

	if err != nil {
		log.Print("Method Part: Validate Token")
		return nil, ErrTokenInternalServerError
	}

	claims, ok := token.Claims.(*JWTClaims)

	if !ok {
		log.Print("Method Part")
		return nil, ErrInvalidToken
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		log.Print("Method Part")
		return nil, ErrExpiredToken
	}

	return claims, nil

}