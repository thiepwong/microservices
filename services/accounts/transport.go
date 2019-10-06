package accounts

import (
	"fmt"

	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/thiepwong/microservices/common"
)

func RegisterRoute(app *iris.Application) {
	// pg, err := datasources.GetPg(*config.Database.Postgre)
	// if err != nil {
	// 	logger.LogErr.Println(err.Error())
	// 	os.Exit(2)
	// }

	//	mvcResult := controllers.NewMvcResult(nil)

	//Register Employee Controller
	//	empRep := repositories.NewEmployeeRepository(pg)
	id := "Da truyen vao"
	empSrv := NewAccountService(id)
	emp := mvc.New(app.Party("/accounts")) //.AllowMethods(iris.MethodOptions, iris.MethodGet, iris.MethodPost))
	emp.Register(empSrv)
	emp.Handle(new(AccountRoute))

}

var tst = func() {
	fmt.Println("Da goi vao ham")
}

type AccountRoute struct {
	common.Context
	Service AccountService
}

func (r *AccountRoute) BeforeActivation(b mvc.BeforeActivation) {
	//r.ApiSecure()
	b.Handle("GET", "/profile/{id:string}", "GetProfile")
}

func (r *AccountRoute) GetProfile(id string) {
	res, e := r.Service.Profile(id)
	if e != nil {
		r.Ctx.Text("Da bi loi")
	}
	fmt.Println(res)
	r.Response(200, "Da gui thanh cong", res)
}
