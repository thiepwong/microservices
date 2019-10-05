package accounts

import (
	"github.com/kataras/iris"
)

func NewService() *iris.Application {
	app := iris.New()
	app.Get("/", func(ctx iris.Context) {
		ctx.Text("Xin chao nhe! Day la trang tao account")
	})
	return app
}
