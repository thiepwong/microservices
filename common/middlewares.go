package common

import "github.com/kataras/iris"

func PreFlight(c iris.Context) {

	if c.Request().Method == "OPTIONS" {
		return
	}

	c.Next()

}
