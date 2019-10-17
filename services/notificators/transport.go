package notificators

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

	notiRepo := NewNotificatorRepository(mongoSession, redis, cfg)
	notiServ := NewNotificatorService(notiRepo, cfg)
	noti := mvc.New(app.Party("/notificator", cors, common.PreFlight).AllowMethods(iris.MethodOptions, iris.MethodGet, iris.MethodPost))
	noti.Register(notiServ)
	noti.Handle(new(NotificatorRoute))

}

type NotificatorRoute struct {
	common.Context
	Service NotificatorService
}

func (r *NotificatorRoute) BeforeActivation(b mvc.BeforeActivation) {
	//r.ApiSecure()
	b.Handle("POST", "/sendmail", "PostSendMail")
	b.Handle("POST", "/sendsms", "PostSendSms")
	// b.Handle("POST", "/mail-forgot-password", "PostMailForgotPassword")
	// b.Handle("POST", "/sms-forgot-password", "PostSmsForgotPassword")

}

func (r *NotificatorRoute) PostSendMail() {
	_mail := &MailModel{}
	err := r.Ctx.ReadJSON(_mail)
	if err != nil {
		r.Response(406, err.Error(), nil)
		return
	}

	if _mail.Mail == "" {
		r.Response(428, "Parammeters is required, please input and submit again!", nil)
		return
	}

	if _mail.Lang == "" {
		_mail.Lang = "vi"
	}

	res, err := r.Service.SendEmail(_mail)
	if err != nil {
		r.Response(500, err.Error(), nil)
		return
	}

	r.Response(200, "Sendmail successfully!", res)
}

func (r *NotificatorRoute) PostSendSms() {
	_sms := &SmsModel{}
	err := r.Ctx.ReadJSON(_sms)
	if err != nil {
		r.Response(406, err.Error(), nil)
		return
	}

	if _sms.Mobile == "" {
		r.Response(428, "Parammeters is required, please input and submit again!", nil)
		return
	}

	res, err := r.Service.SendSMS(_sms)
	if err != nil {
		r.Response(500, err.Error(), nil)
		return
	}

	r.Response(200, "Sms was sent successfully", res)

}

// func (r *NotificatorRoute) PostMailForgotPassword() {
// 	_sms := &MailModel{}
// 	err := r.Ctx.ReadJSON(_sms)
// 	if err != nil {
// 		r.Response(406, err.Error(), nil)
// 		return
// 	}

// 	if _sms.Mobile == "" {
// 		r.Response(428, "Parammeters is required, please input and submit again!", nil)
// 		return
// 	}

// 	res, err := r.Service.SendSMS(_sms)
// 	if err != nil {
// 		r.Response(500, err.Error(), nil)
// 		return
// 	}

// 	r.Response(200, "Sms was sent successfully", res)

// }

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
