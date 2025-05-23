package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context/param"
)

func init() {

	beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:AudioController"] = append(beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:AudioController"],
		beego.ControllerComments{
			Method:           "CreateAudioTranscription",
			Router:           `/transcription`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(
				param.New("request"),
			),
			Filters: nil,
			Params:  nil})

	beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:AudioController"] = append(beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:AudioController"],
		beego.ControllerComments{
			Method:           "CreateAudioTranslation",
			Router:           `/translation`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(
				param.New("request"),
			),
			Filters: nil,
			Params:  nil})

	beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:ChatController"] = append(beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:ChatController"],
		beego.ControllerComments{
			Method:           "CreateChat",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(
				param.New("request"),
			),
			Filters: nil,
			Params:  nil})

	beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:ChatController"] = append(beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:ChatController"],
		beego.ControllerComments{
			Method:           "CreateChatWithSSE",
			Router:           `/sse`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:CompletionController"] = append(beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:CompletionController"],
		beego.ControllerComments{
			Method:           "CreateCompletion",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:CompletionController"] = append(beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:CompletionController"],
		beego.ControllerComments{
			Method:           "CreateCompletionWithSSE",
			Router:           `/sse`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:EditController"] = append(beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:EditController"],
		beego.ControllerComments{
			Method:           "CreateEdit",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:ExampleController"] = append(beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:ExampleController"],
		beego.ControllerComments{
			Method:           "Examples",
			Router:           `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:ExampleController"] = append(beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:ExampleController"],
		beego.ControllerComments{
			Method:           "CreateExample",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:ExampleController"] = append(beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:ExampleController"],
		beego.ControllerComments{
			Method:           "InitExamples",
			Router:           `/init`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:ExampleController"] = append(beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:ExampleController"],
		beego.ControllerComments{
			Method:           "ResetExamples",
			Router:           `/reset`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:FineTuneController"] = append(beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:FineTuneController"],
		beego.ControllerComments{
			Method:           "ListFineTuneEvent",
			Router:           `/:fine_tune_id/events`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:ModelController"] = append(beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:ModelController"],
		beego.ControllerComments{
			Method:           "ListModel",
			Router:           `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:ObjectController"] = append(beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:ObjectController"],
		beego.ControllerComments{
			Method:           "Post",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:ObjectController"] = append(beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:ObjectController"],
		beego.ControllerComments{
			Method:           "GetAll",
			Router:           `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:ObjectController"] = append(beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:ObjectController"],
		beego.ControllerComments{
			Method:           "Get",
			Router:           `/:objectId`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:ObjectController"] = append(beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:ObjectController"],
		beego.ControllerComments{
			Method:           "Put",
			Router:           `/:objectId`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:ObjectController"] = append(beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:ObjectController"],
		beego.ControllerComments{
			Method:           "Delete",
			Router:           `/:objectId`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:TestController"] = append(beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:TestController"],
		beego.ControllerComments{
			Method:           "Test",
			Router:           `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:UserController"],
		beego.ControllerComments{
			Method:           "Post",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:UserController"],
		beego.ControllerComments{
			Method:           "GetAll",
			Router:           `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:UserController"],
		beego.ControllerComments{
			Method:           "Get",
			Router:           `/:uid`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:UserController"],
		beego.ControllerComments{
			Method:           "Put",
			Router:           `/:uid`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:UserController"],
		beego.ControllerComments{
			Method:           "Delete",
			Router:           `/:uid`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:UserController"],
		beego.ControllerComments{
			Method:           "Login",
			Router:           `/login`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:UserController"],
		beego.ControllerComments{
			Method:           "Logout",
			Router:           `/logout`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}
