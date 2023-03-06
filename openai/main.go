package main

import (
	_ "openai/routers"

	"github.com/beego/beego/v2/server/web"
)

func main() {
	// fmt.Println(web.BConfig.RunMode)
	web.BConfig.RouterCaseSensitive = false

	if web.BConfig.RunMode == "dev" {
		web.BConfig.WebConfig.DirectoryIndex = true
		web.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	web.SetStaticPath("/", "static")
	// web.AutoPrefix("api", &controller.UserController{})
	// web.AutoPrefix("api", &controller.ModelController{})
	// web.AutoPrefix("api", &controller.ChatController{})
	// web.AutoPrefix("api", &controller.FineTuneController{})
	// tree := web.PrintTree()
	// methods := tree["Data"].(web.M)
	// for k, v := range methods {
	// 	fmt.Printf("%s => %v\n", k, v)
	// }
	web.Run()
}
