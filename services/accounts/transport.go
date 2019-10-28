package accounts

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
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
	acc := mvc.New(app.Party("/profile", cors, common.PreFlight).AllowMethods(iris.MethodOptions, iris.MethodGet, iris.MethodPost))
	acc.Register(accSrv)
	acc.Handle(new(AccountRoute))
}

type AccountRoute struct {
	common.Context
	Service AccountService
}

func (r *AccountRoute) BeforeActivation(b mvc.BeforeActivation) {
	//r.ApiSecure()
	b.Handle("GET", "/user-profile/{sid:uint64}", "GetProfile", common.AccessAuth)
	b.Handle("POST", "/register", "PostRegister")
	b.Handle("POST", "/{sid:string}/email", "PostUpdateEmail")
	b.Handle("POST", "/{sid:string}/mobile", "PostUpdateMobile")
	b.Handle("POST", "/update/{sid:uint64}", "PostUpdateProfile", common.AccessAuth)

}

func (r *AccountRoute) GetProfile(sid uint64) {
	if sid == 0 {
		r.Response(428, "Smart ID is required for request profile", nil)
		return
	}

	token := r.Ctx.URLParam("token")

	res, e := r.Service.Profile(sid, token)
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

	r.Response(200, "Register successful! Please verify the account to active Smart ID", res)
}

func (r *AccountRoute) PostUpdateEmail(sid string) {
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

	res, err := r.Service.UpdateEmail(&_profile)
	if err != nil {
		r.Response(500, err.Error(), err)
		return
	}

	r.Response(200, "", res)
}

func (r *AccountRoute) PostUpdateMobile(sid string) {
	var _profile AuthUpdate
	err := r.Ctx.ReadJSON(&_profile)
	if err != nil {
		r.Response(406, err.Error(), nil)
	}
	if _profile.Mobile == "" {
		r.Response(428, "Mobile is required!", nil)
	}

	res, err := r.Service.UpdateMobile(&_profile)
	if err != nil {
		r.Response(500, err.Error(), err)
		return
	}
	r.Response(200, "", res)

}

func (r *AccountRoute) PostUpdateProfile(sid uint64) {
	var _profile Profile
	err := r.Ctx.ReadJSON(&_profile)
	_usr := r.Ctx.Values().Get("user")
	__u := _usr.(jwt.MapClaims)
	if err != nil {
		r.Response(406, err.Error(), nil)
		return
	}
	// __sid, err := strconv.ParseUint(__u["sid"].map(string), 10, 64)
	__sid := uint64(__u["sid"].(float64))
	if err != nil {
		r.Response(406, err.Error(), nil)
		return
	}

	res, err := r.Service.UpdateProfile(&_profile, __sid)
	if err != nil {
		r.Response(500, err.Error(), err)
		return
	}
	r.Response(200, "", res)

}
