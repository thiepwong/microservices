// package main

// import (
// 	"fmt"

// 	"github.com/kataras/iris"
// 	"github.com/thiepwong/microservices/services/accounts"
// 	"github.com/thiepwong/microservices/services/notificators"
// )

// func main() {
// 	fmt.Println("Da khoi dong dich vu")
// 	account := accounts.NewService()
// 	noti := notificators.NewService()
// 	go func() {
// 		account.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
// 	}()

// 	noti.Run(iris.Addr(":8081"), iris.WithoutServerError(iris.ErrServerClosed))

// }

package main

import (
	"fmt"
	"os"

	"github.com/kataras/iris"
	"github.com/thiepwong/microservices/services/accounts"
	"github.com/thiepwong/microservices/services/notificators"
	"github.com/urfave/cli"
)

var (
	flags []cli.Flag
	host  string
	port  int
)

func init() {
	flags = []cli.Flag{
		cli.StringFlag{
			Name:  "t,target,host",
			Value: "localhost",
			Usage: "Start IP/domain",
		},
		cli.IntFlag{
			Name:  "p,port",
			Value: 8080,
			Usage: "Start port",
		},
	}
}

func alert() {
	fmt.Println("Press 1 to account")
	fmt.Println("Press 2 to notify")
	fmt.Println("Press 0 to exit")
}

func main() {
	alert()
	for {
		sample()
	}
}

func sample() {
	var input int
	n, err := fmt.Scanln(&input)
	if n < 1 || err != nil {
		fmt.Println("invalid input")
		return
	}
	switch input {
	case 1:
		fmt.Println("Starting account service ...")
		go func() {
			account := accounts.NewService()
			account.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
		}()
		alert()
	case 2:
		fmt.Println("Starting notification service ...")
		go func() {
			notify := notificators.NewService()
			notify.Run(iris.Addr(":8081"), iris.WithoutServerError(iris.ErrServerClosed))
		}()
		alert()
	case 3:

		os.Exit(2)
	default:
		alert()
		fmt.Println("def")
	}
}

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

func noArgs(c *cli.Context) error {
	if c.NumFlags() < 2 {
		cli.ShowAppHelp(c)
		return cli.NewExitError("please set both flags", 2)
	}
	fmt.Printf("hacking %s:%d", host, port)
	//	go func() {
	account := accounts.NewService()
	account.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
	//}()

	return nil

}
