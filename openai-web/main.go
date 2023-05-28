package main

import (
	"os"

	_ "github.com/joho/godotenv/autoload"

	_ "github.com/futugyousuzu/go-openai-web/routers"
	_ "github.com/futugyousuzu/go-openai-web/services"
	_ "github.com/futugyousuzu/openai-tokenizer"

	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/filter/cors"
	"github.com/futugyousuzu/go-openai-web/middleware/oauth"
)

func init() {
	web.InsertFilter("*", web.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"'GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"Access-Control-Allow-Origin", "Origin", "Authorization", "Access-Control-Allow-Headers", "Content-Type"},
		ExposeHeaders:    []string{"Access-Control-Allow-Origin", "Content-Length", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
	}))

	web.InsertFilter("*", web.BeforeRouter, oauth.OAuthConfig(&oauth.Options{
		AuthServerURL: os.Getenv("auth_server_url"),
		ClientID:      os.Getenv("client_id"),
		ClientSecret:  os.Getenv("client_secret"),
		Scopes:        os.Getenv("scopes"),
		RedirectURL:   os.Getenv("redirect_url"),
		AuthURL:       os.Getenv("auth_url"),
		TokenURL:      os.Getenv("token_url"),
	}))
}

func main() {
	// fmt.Println(web.BConfig.RunMode)
	web.BConfig.RouterCaseSensitive = false

	if web.BConfig.RunMode == "dev" {
		web.BConfig.WebConfig.DirectoryIndex = true
		web.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	web.SetStaticPath("/", "static")

	// tree := web.PrintTree()
	// methods := tree["Data"].(web.M)
	// for k, v := range methods {
	// 	fmt.Printf("%s => %v\n", k, v)
	// }

	web.Run()
}
