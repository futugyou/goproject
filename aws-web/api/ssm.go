package api

import (
	_ "github.com/joho/godotenv/autoload"

	"encoding/json"
	"net/http"

	"github.com/futugyousuzu/goproject/awsgolang/services"
	verceltool "github.com/futugyousuzu/goproject/awsgolang/vercel"
	model "github.com/futugyousuzu/goproject/awsgolang/viewmodel"

	"github.com/futugyou/extensions"
)

func SSM(w http.ResponseWriter, r *http.Request) {
	if extensions.Cors(w, r) {
		return
	}

	operatetype := r.URL.Query().Get("type")
	switch operatetype {
	case "search":
		searchSSMData(w, r)
	}
}

func searchSSMData(w http.ResponseWriter, r *http.Request) {
	if !verceltool.AuthForVercel(w, r) {
		return
	}

	service := services.NewSSMService()
	accountId := r.URL.Query().Get("accountId")
	name := r.URL.Query().Get("name")
	filter := model.SSMDataFilter{AccountId: accountId, Name: name}
	datas, err := service.SearchSSMData(r.Context(), filter)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
		return
	}
	body, _ := json.Marshal(datas)
	w.Write(body)
	w.WriteHeader(200)
}
