package main

import (
	"fmt"
	"openai/pkg/controller"
	_ "openai/pkg/routers"

	"github.com/beego/beego/v2/server/web"
)

func main() {
	fmt.Println(web.BConfig.RunMode)
	if web.BConfig.RunMode == "dev" {
		web.BConfig.WebConfig.DirectoryIndex = true
		web.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	web.SetStaticPath("/", "static")
	ctrl := &controller.UserController{}
	web.AutoPrefix("api", ctrl)
	web.Run()

}
