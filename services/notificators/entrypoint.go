package notificators

import (
	"github.com/kataras/iris"
)

func NewService() *iris.Application {
	app := iris.New()
	app.Get("/", func(ctx iris.Context) {
		ctx.Text("Day la trang gui mail nhe!!")

	})
	return app
}
