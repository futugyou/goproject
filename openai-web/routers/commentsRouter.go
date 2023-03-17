package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:AudioController"] = append(beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:AudioController"],
        beego.ControllerComments{
            Method: "CreateAudioTranscription",
            Router: `/transcription`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(
				param.New("request"),
			),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:AudioController"] = append(beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:AudioController"],
        beego.ControllerComments{
            Method: "CreateAudioTranslation",
            Router: `/translation`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(
				param.New("request"),
			),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:ChatController"] = append(beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:ChatController"],
        beego.ControllerComments{
            Method: "CreateChat",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(
				param.New("request"),
			),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:FineTuneController"] = append(beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:FineTuneController"],
        beego.ControllerComments{
            Method: "ListFineTuneEvent",
            Router: `/:fine_tune_id/events`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:ModelController"] = append(beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:ModelController"],
        beego.ControllerComments{
            Method: "ListModel",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:ObjectController"] = append(beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:ObjectController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:ObjectController"] = append(beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:ObjectController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:ObjectController"] = append(beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:ObjectController"],
        beego.ControllerComments{
            Method: "Get",
            Router: `/:objectId`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:ObjectController"] = append(beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:ObjectController"],
        beego.ControllerComments{
            Method: "Put",
            Router: `/:objectId`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:ObjectController"] = append(beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:ObjectController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: `/:objectId`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:QuestionController"] = append(beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:QuestionController"],
        beego.ControllerComments{
            Method: "CreateQAndA",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(
				param.New("request"),
			),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:UserController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:UserController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:UserController"],
        beego.ControllerComments{
            Method: "Get",
            Router: `/:uid`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:UserController"],
        beego.ControllerComments{
            Method: "Put",
            Router: `/:uid`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:UserController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: `/:uid`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:UserController"],
        beego.ControllerComments{
            Method: "Login",
            Router: `/login`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/futugyousuzu/go-openai-web/controllers:UserController"],
        beego.ControllerComments{
            Method: "Logout",
            Router: `/logout`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
