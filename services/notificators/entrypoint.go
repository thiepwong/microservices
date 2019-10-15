package notificators

import (
	"github.com/kataras/iris"
	"github.com/thiepwong/microservices/common"
)

func NewService(cfg *common.Config) *iris.Application {
	app := iris.Default()

	crs := func(ctx iris.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Credentials", "true")
		ctx.Header("Access-Control-Allow-Headers", "Access-Control-Allow-Origin,Content-Type, Authorization")
		ctx.Next()
	}

	app.Logger().SetLevel("debug")
	RegisterRoute(app,crs, cfg)
	// app.Get("/", func(ctx iris.Context) {
	// 	ctx.Text("Xin chao nhe! Day la trang tao account")
	// })
	return app
}
