// @APIVersion 1.0.0
// @Title openai web API
// @Description provider some api for openai demo.
package routers

import (
	"openai/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	ns := beego.NewNamespace("/api/v1",
		beego.NSNamespace("/model",
			beego.NSInclude(
				&controllers.ModelController{}),
		),
		beego.NSNamespace("/chat",
			beego.NSInclude(
				&controllers.ChatController{}),
		),
		beego.NSNamespace("/fine-tune",
			beego.NSInclude(
				&controllers.FineTuneController{}),
		),
	)
	beego.AddNamespace(ns)
}
