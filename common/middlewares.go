package common

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris"
)

type TokenClaims struct {
	Usr string        `json:"usr"`
	Iss string        `json:"iss"`
	Act string        `json:"act"`
	Sid uint64        `json:"sid"`
	Jit string        `json:"jit"`
	Iat time.Duration `json:"iat"`
	Exp time.Duration `json:"exp"`
	Sys string        `json:"sys"`
}

func PreFlight(c iris.Context) {

	if c.Request().Method == "OPTIONS" {
		return
	}

	c.Next()

}

// func AccessOwnerAuth(c iris.Context) {

// }

func AccessAuth(c iris.Context) {
	token := c.URLParam("token")
	if token == "" {
		return
	}

	tk, err := tokenValidate(token)
	if err != nil {
		var _result = Result{
			Version:     "0.1.0",
			Code:        401,
			Message:     "Access token is invalid, you cannot access this page",
			System:      "smartid.account",
			RequestTime: time.Now().Unix(),
		}
		c.StatusCode(iris.StatusUnauthorized)
		c.JSON(_result)
		// c.WriteString("loi khong co quyen")
		return
		//c.JSON(_result)
	}
	c.Values().Set("user", tk)

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
	return token.Claims, nil
}
