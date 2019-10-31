package images

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/mvc"
	"github.com/thiepwong/microservices/common"
)

func RegisterRoute(app *iris.Application, cors context.Handler, cfg *common.Config) {

	imgService := NewImageService(cfg)
	acc := mvc.New(app.Party("/images", cors, common.PreFlight).AllowMethods(iris.MethodOptions, iris.MethodGet, iris.MethodPost))
	acc.Register(imgService)
	acc.Handle(new(AccountRoute))
}

type AccountRoute struct {
	common.Context
	Service ImageService
}

func (r *AccountRoute) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("POST", "/upload/{sid:uint64}", "PostUpload", common.AccessAuth)
	b.Handle("GET", "/list/{sid:uint64}", "GetUploadList", common.AccessAuth)
	// b.Handle("GET", "/{sid:uint64}", "PostGetAllImage", common.AccessAuth)
}

func (r *AccountRoute) PostUpload(sid uint64) {
	var _img Image
	err := r.Ctx.ReadJSON(&_img)
	if err != nil {
		r.Response(406, err.Error(), nil)
		return
	}
	res, err := r.Service.Upload(&_img, sid)
	if err != nil {
		r.Response(500, err.Error(), nil)
		return
	}

	r.Response(200, "Image was successful uploaded!", res)
}

func (r *AccountRoute) GetUploadList(sid uint64) {

	res, err := r.Service.List(sid)
	if err != nil {
		r.Response(500, err.Error(), nil)
		return
	}

	r.Response(200, "", res)
}
