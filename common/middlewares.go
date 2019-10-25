package common

import (
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris"
)

func PreFlight(c iris.Context) {

	if c.Request().Method == "OPTIONS" {
		return
	}

	c.Next()

}

func AccessAuth(c iris.Context) {
	token := c.URLParam("token")
	if token == "" {
		return
	}

	_, err := tokenValidate(token)
	if err != nil {
		return
	}

	c.Next()

}

func tokenValidate(tokenString string) (interface{}, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		rsa, err := ReadPublicKey("./1010.pub")
		// key, err := ioutil.ReadFile("your-private-key.pem")
		if err != nil {
			return nil, errors.New("private key could not be loaded")
		}
		return rsa, nil
	})
	if err != nil {
		return nil, errors.New("Token key is invalid")
	}
	fmt.Print(token)
	return token, nil
}
