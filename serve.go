package main

import (
	"fmt"
	"os"

	"github.com/thiepwong/microservices/services/images"

	"github.com/thiepwong/microservices/services/notificators"

	"github.com/thiepwong/microservices/common"

	"github.com/kataras/iris"
	"github.com/thiepwong/microservices/services/accounts"
	"github.com/thiepwong/microservices/services/auth"
)

func menuShow() {
	fmt.Println("1. Account service")
	fmt.Println("2. Notification service")
	fmt.Println("3. Authentication service")
	fmt.Println("4. Management service")
	fmt.Println("5. Exit")
	fmt.Println("Type start and name of service to start it, Type stop and name of service to stop it")
}

func alert() {
	fmt.Println("Press 1 to account")
	fmt.Println("Press 2 to notify")
	fmt.Println("Press 0 to exit")
}

func main() {
	menuShow()
	go func() {
		startAccount()
		startAuth()
		startNotificator()
		startImage()
	}()

	for {
		menuHandler()
	}
}

func menuHandler() {
	var cmd, service string
	n, err := fmt.Scanf("%s %s", &cmd, &service)
	if n < 1 || err != nil {
		fmt.Println("invalid input")
		return
	}

	switch service {
	case "account":
		startAccount()

	case "auth":
		startAuth()
	}

	fmt.Println("Da nhan: " + cmd)
	fmt.Println("Da nhan: " + service)
}

func startAccount() {
	go func() {
		_cfgPath := "configs/account.yaml"
		conf, es := common.LoadConfig(_cfgPath)
		if es != nil {
			os.Exit(10)
		}
		serv := accounts.NewService(conf)
		serv.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))

	}()
}

func startAuth() {
	go func() {
		_cfgPath := "configs/account.yaml"
		conf, es := common.LoadConfig(_cfgPath)
		if es != nil {
			os.Exit(10)
		}
		serv := auth.NewService(conf)
		serv.Run(iris.Addr(":8081"), iris.WithoutServerError(iris.ErrServerClosed))

	}()
}

func startNotificator() {
	go func() {
		_cfgPath := "configs/account.yaml"
		conf, es := common.LoadConfig(_cfgPath)
		if es != nil {
			os.Exit(10)
		}
		serv := notificators.NewService(conf)
		serv.Run(iris.Addr(":8082"), iris.WithoutServerError(iris.ErrServerClosed))

	}()
}

func startImage() {
	go func() {
		_cfgPath := "configs/image.yaml"
		conf, es := common.LoadConfig(_cfgPath)
		if es != nil {
			os.Exit(10)
		}
		serv := images.NewService(conf)
		serv.Run(iris.Addr(fmt.Sprintf(":%d", conf.Service.Port)), iris.WithoutServerError(iris.ErrServerClosed))

	}()
}

// func sample() {
// 	var input int
// 	n, err := fmt.Scanln(&input)
// 	if n < 1 || err != nil {
// 		fmt.Println("invalid input")
// 		return
// 	}
// 	switch input {
// 	case 1:
// 		fmt.Println("Starting account service ...")
// 		go func() {
// 			account := accounts.NewService()
// 			account.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
// 		}()
// 		alert()
// 	case 2:
// 		fmt.Println("Starting notification service ...")
// 		go func() {
// 			notify := notificators.NewService()
// 			notify.Run(iris.Addr(":8081"), iris.WithoutServerError(iris.ErrServerClosed))
// 		}()
// 		alert()
// 	case 3:

// 		os.Exit(2)
// 	default:
// 		alert()
// 		fmt.Println("def")
// 	}
// }

// func main() {
// 	app := cli.NewApp()
// 	app.Name = "greet"
// 	app.Usage = "fight the loneliness!"
// 	app.Flags = flags
// 	app.Action = noArgs

// 	err := app.Run(os.Args)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	if 1 == 1 {

// 	}
// }

// func noArgs(c *cli.Context) error {
// 	if c.NumFlags() < 2 {
// 		cli.ShowAppHelp(c)
// 		return cli.NewExitError("please set both flags", 2)
// 	}
// 	fmt.Printf("hacking %s:%d", host, port)
// 	//	go func() {
// 	account := accounts.NewService()
// 	account.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
// 	//}()

// 	return nil

// }
