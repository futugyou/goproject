package main

import (
	"openai/pkg/controller"

	"github.com/beego/beego/v2/server/web"
)

func main() {
	web.SetStaticPath("/", "static")
	ctrl := &controller.UserController{}
	web.AutoPrefix("api", ctrl)
	web.Run()

}
