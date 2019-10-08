package notificators

import (
	"github.com/kataras/iris"
	"github.com/thiepwong/microservices/common"
)

func RegisterRoute(app *iris.Application, cfg *common.Config) {
	// mongdb, err := db.GetMongoDb(cfg.Database.Mongo)
	// if err != nil {
	// }

}

//	mvcResult := controllers.NewMvcResult(nil)

// 	// Register NotificatorRoute
// 	notiRep := NewNotificatorReportsitory(mongdb, "notificator")
// 	notiSrv := NewNotificatorService(accRep, cfg)
// 	noti := mvc.New(app.Party("/notificator")) //.AllowMethods(iris.MethodOptions, iris.MethodGet, iris.MethodPost))
// 	noti.Register(accSrv)
// 	noti.Handle(new(NotificatorRoute))

// }

// type NotificatorRoute struct {
// 	common.Context
// 	Service NotificatorService
// }

// func (r *NotificatorRoute) BeforeActivation(b mvc.BeforeActivation) {
// 	//r.ApiSecure()
// 	b.Handle("GET", "/profile/{id:string}", "GetProfile")
// 	b.Handle("POST", "/register", "PostRegister")
// }

// func (r *NotificatorRoute) GetProfile(id string) {
// 	res, e := r.Service.Profile(id)
// 	if e != nil {
// 		r.Ctx.Text("Da bi loi")
// 	}
// 	fmt.Println(res)
// 	r.Response(200, "Da gui thanh cong", res)
// }

// func (r *NotificatorRoute) PostRegister() {
// 	var _registerModel RegisterModel
// 	err := r.Ctx.ReadJSON(&_registerModel)
// 	if err != nil {
// 		r.Response(406, err.Error(), nil)
// 		return
// 	}

// 	if _registerModel.Username == "" || _registerModel.Password == "" {
// 		r.Response(428, "Username and Password is required!", nil)
// 		return
// 	}

// 	res, e := r.Service.Register(&_registerModel)
// 	if e != nil {
// 		r.Response(500, e.Error(), nil)
// 		return
// 	}

// 	r.Response(200, "", res)
// }
