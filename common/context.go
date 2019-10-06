package common

import (
	"fmt"
	"time"

	"github.com/kataras/iris"
)

// Result Struct for return
type Result struct {
	System      string      `json:"system"`
	Version     string      `json:"version"`
	RequestTime int64       `json:"request_time"`
	Code        int         `json:"result_code"`
	Message     string      `json:"message"`
	Data        interface{} `json:"data"`
}

// Context struct for Route
type Context struct {
	Ctx    iris.Context
	Secure bool
}

// Response method
func (c *Context) Response(code int, msg string, d interface{}) {

	if code == 0 {
		code = 200
	}

	if msg == "" {
		msg = "Successful"
	}

	result := Result{
		Version:     "0.1.0",
		System:      "smartid.account",
		RequestTime: time.Now().Unix(),
		Code:        code,
		Message:     msg,
		Data:        d,
	}

	c.Ctx.JSON(result)

}

// Request before
func (c *Context) Request(cnt iris.Context) {
	fmt.Println(cnt.URLParam("id"))
	if c.Secure == true {
		fmt.Println("Da su dung che do bao mat")
	} else {
		fmt.Println("Khong bao mat!")
	}
}

//BeginRequest method
func (c *Context) BeginRequest(ctx iris.Context) {
	fmt.Println(ctx.URLParam("id"))
}

func (c *Context) EndRequest(ctx iris.Context) {}

// func (c *Context) Auth() iris.Context.Handlers {

// }

type Auth struct {
	Token  string
	IsAuth bool
}

func (a *Auth) Check() bool {
	if a.IsAuth == true {
		if a.Token == "" {
			return false
		} else {
			return true
		}
	}
	return false
}
