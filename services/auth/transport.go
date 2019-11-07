package auth

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/mvc"
	"github.com/thiepwong/microservices/common"
	"github.com/thiepwong/microservices/common/db"
)

func RegisterRoute(app *iris.Application, cors context.Handler, cfg *common.Config) {
	mongoSession := db.GetMongoSession(cfg.Database.Mongo)
	redis := db.GetRedisDb(cfg.Database.Redis)

	// Register Account Route
	accRep := NewAuthRepository(mongoSession, redis, cfg)
	accSrv := NewAuthService(accRep, cfg)
	acc := mvc.New(app.Party("/auth", cors, common.PreFlight).AllowMethods(iris.MethodOptions, iris.MethodGet, iris.MethodPost))
	acc.Register(accSrv)
	acc.Handle(new(AuthRoute))

}

type AuthRoute struct {
	common.Context
	Service AuthService
}

func (r *AuthRoute) BeforeActivation(b mvc.BeforeActivation) {
	// sign in with normal method, using username and password
	b.Handle("POST", "/signin", "PostSignIn")
	// sign in with social network token
	b.Handle("POST", "social-network-signin", "PostSocialSignIn")

	// activate an register and create a smart id
	b.Handle("GET", "/activate", "GetActivate")

	// update contact and combine user
	b.Handle("GET", "/update-contact", "GetUpdateContact")

	// Change password
	b.Handle("POST", "/change-password", "PostChangePassword", common.AccessAuth)

	// Create new password
	b.Handle("POST", "/create-new-password", "PostCreateNewPassword")

	b.Handle("GET", "/refresh-token/{rft:string}", "GetRefreshToken")
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
func (r *AuthRoute) GetRefreshToken(rft string) {

	res, e := r.Service.SignInViaRefreshToken(rft)
	if e != nil {
		r.Response(500, e.Error(), nil)
		return
	}

	r.Response(200, "", res)
}

func (r *AuthRoute) GetActivate() {
	var _activate ActivateModel
	_activate.Username = r.Ctx.URLParam("username")
	_activate.ActivateCode = r.Ctx.URLParam("code")

	if _activate.Username == "" || _activate.ActivateCode == "" {
		r.Response(428, "Username and activate code is required!", nil)
		return
	}

	res, err := r.Service.Verify(&_activate)
	if err != nil {
		r.Response(500, err.Error(), nil)
		return
	}
	r.Response(200, "", res)
}

func (r *AuthRoute) GetUpdateContact() {
	var _verify UpdateContact
	_verify.Contact = r.Ctx.URLParam("contact")
	_verify.Code = r.Ctx.URLParam("code")

	if _verify.Contact == "" || _verify.Code == "" {
		r.Response(428, "Contact and verify code is required!", nil)
		return
	}

	res, err := r.Service.UpdateContact(&_verify)
	if err != nil {
		r.Response(500, err.Error(), nil)
		return
	}
	r.Response(200, "", res)
}

func (r *AuthRoute) PostSocialSignIn() {

}

func (r *AuthRoute) PostChangePassword() {
	var _data ChangePasswordModel
	err := r.Ctx.ReadJSON(&_data)
	if err != nil {
		r.Response(406, "Json data is required!", err)
		return
	}

	if _data.Username == "" {
		r.Response(208, "Username is required!", nil)
		return
	}

	if _data.OldPassword == "" {
		r.Response(208, "Old password is required!", nil)
		return
	}

	if _data.NewPassword == "" {
		r.Response(208, "New password is required!", nil)
		return
	}
	res, err := r.Service.ChangePassword(&_data)

	if err != nil {
		r.Response(500, err.Error(), err)
		return
	}

	r.Response(200, "", res)

}

func (r *AuthRoute) PostCreateNewPassword() {

	var _cont VerifyContact

	err := r.Ctx.ReadJSON(&_cont)
	if err != nil {
		r.Response(416, "Request is invalid!", nil)
		return
	}

	if _cont.Contact == "" || _cont.Code == "" {
		r.Response(428, "Contact infomation and verify code is required!", nil)
		return
	}

}
