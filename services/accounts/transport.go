package accounts

import (
	"fmt"

	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/thiepwong/microservices/common"
	"github.com/thiepwong/microservices/common/db"
	"github.com/thiepwong/smartid/pkg/logger"
)

func RegisterRoute(app *iris.Application, cfg *common.Config) {
	mongdb, err := db.GetMongoDb(cfg.Database.Mongo)
	if err != nil {
		logger.LogErr.Println(err.Error())
	}

	//	mvcResult := controllers.NewMvcResult(nil)

	// Register Account Route
	accRep := NewAccountReportsitory(mongdb, "accounts")
	accSrv := NewAccountService(accRep, cfg)
	acc := mvc.New(app.Party("/accounts")) //.AllowMethods(iris.MethodOptions, iris.MethodGet, iris.MethodPost))
	acc.Register(accSrv)
	acc.Handle(new(AccountRoute))

}

type AccountRoute struct {
	common.Context
	Service AccountService
}

func (r *AccountRoute) BeforeActivation(b mvc.BeforeActivation) {
	//r.ApiSecure()
	b.Handle("GET", "/profile/{id:string}", "GetProfile")
	b.Handle("POST", "/register", "PostRegister")
}

func (r *AccountRoute) GetProfile(id string) {
	res, e := r.Service.Profile(id)
	if e != nil {
		r.Ctx.Text("Da bi loi")
	}
	fmt.Println(res)
	r.Response(200, "Da gui thanh cong", res)
}

func (r *AccountRoute) PostRegister() {
	var _registerModel RegisterModel
	err := r.Ctx.ReadJSON(&_registerModel)
	if err != nil {
		r.Response(406, err.Error(), nil)
		return
	}

	if _registerModel.Username == "" || _registerModel.Password == "" {
		r.Response(428, "Username and Password is required!", nil)
		return
	}

	res, e := r.Service.Register(&_registerModel)
	if e != nil {
		r.Response(500, e.Error(), nil)
		return
	}

	r.Response(200, "", res)
}
