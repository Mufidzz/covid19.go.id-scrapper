package Utils

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

var issuerID = "9cf010930e65e34fa7afa03dca069117b1337ac6"
var signingKey = []byte("07dcf32f768528e2312377cc235429243b93578e")

func GenerateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS512)

	c := token.Claims.(jwt.MapClaims)
	c["iat"] = time.Now().Unix()
	c["expr"] = time.Now().Add(time.Minute * 5).Unix()
	c["iss"] = issuerID

	tS, err := token.SignedString(signingKey)

	if err != nil {
		return "", err
	}
	return tS, nil
}
