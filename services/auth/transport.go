package auth

import (
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
	accRep := NewAuthRepository(mongdb, "accounts")
	accSrv := NewAuthService(accRep, cfg)
	acc := mvc.New(app.Party("/auth")) //.AllowMethods(iris.MethodOptions, iris.MethodGet, iris.MethodPost))
	acc.Register(accSrv)
	acc.Handle(new(AuthRoute))

}

type AuthRoute struct {
	common.Context
	Service AuthService
}

func (r *AuthRoute) BeforeActivation(b mvc.BeforeActivation) {
	//r.ApiSecure()
	b.Handle("POST", "/signin", "PostSignIn")
}

func (r *AuthRoute) PostSignIn() {
	var _loginModel SignInModel
	err := r.Ctx.ReadJSON(&_loginModel)
	if err != nil {
		r.Response(406, err.Error(), nil)
		return
	}

	if _loginModel.Username == "" || _loginModel.Password == "" {
		r.Response(428, "Username and Password is required!", nil)
		return
	}

	res, e := r.Service.SignIn(&_loginModel)
	if e != nil {
		r.Response(500, e.Error(), nil)
		return
	}

	r.Response(200, "", res)
}
