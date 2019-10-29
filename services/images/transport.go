package images

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/mvc"
	"github.com/thiepwong/microservices/common"
	"github.com/thiepwong/microservices/common/db"
)

func RegisterRoute(app *iris.Application, cors context.Handler, cfg *common.Config) {
	mongoSession := db.GetMongoSession(cfg.Database.Mongo)
	// Register Account Route
	accRep := NewAccountReportsitory(mongoSession, cfg)
	accSrv := NewAccountService(accRep, cfg)
	acc := mvc.New(app.Party("/images", cors, common.PreFlight).AllowMethods(iris.MethodOptions, iris.MethodGet, iris.MethodPost))
	acc.Register(accSrv)
	acc.Handle(new(AccountRoute))
}

type AccountRoute struct {
	common.Context
	Service AccountService
}

func (r *AccountRoute) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("POST", "/upload/{sid:uint64}", "PostUpload", common.AccessAuth)
	b.Handle("GET", "/{sid:uint64}", "PostGetAllImage", common.AccessAuth)
}

func (r *AccountRoute) PostUpload(sid uint64) {
	var _img Image
	err := r.Ctx.ReadJSON(&_img)
	if err != nil {
		r.Response(406, err.Error(), nil)
		return
	}
	res, err := r.Service.Upload(&_img)
	if err != nil {
		r.Response(500, err.Error(), nil)
		return
	}

	r.Response(200, "Register successful! Please verify the account to active Smart ID", res)
}

func (r *AccountRoute) PostGetAllImage(sid uint64) {
	var _profile AuthUpdate
	err := r.Ctx.ReadJSON(&_profile)
	if err != nil {
		r.Response(406, err.Error(), nil)
		return
	}
	if _profile.Email == "" {
		r.Response(428, "Email is required!", nil)
		return
	}

	r.Response(200, "", nil)
}
