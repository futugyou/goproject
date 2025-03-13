// @APIVersion 1.0.0
// @Title openai web API
// @Description provider some api for openai demo.
package routers

import (
	"github.com/futugyousuzu/go-openai-web/controllers"

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
		beego.NSNamespace("/autio",
			beego.NSInclude(
				&controllers.AudioController{}),
		),
		beego.NSNamespace("/examples",
			beego.NSInclude(
				&controllers.ExampleController{}),
		),
		beego.NSNamespace("/completions",
			beego.NSInclude(
				&controllers.CompletionController{}),
		),
		beego.NSNamespace("/edits",
			beego.NSInclude(
				&controllers.EditController{}),
		),
	)
	beego.AddNamespace(ns)
}
