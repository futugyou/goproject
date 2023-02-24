package main

import (
	"openai/pkg/controller"

	"github.com/beego/beego/v2/server/web"
)

func main() {
	web.AutoRouter(&controller.UserController{})
	web.Run()

}
