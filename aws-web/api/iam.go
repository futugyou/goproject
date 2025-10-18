package api

import (
	_ "github.com/joho/godotenv/autoload"

	"encoding/json"
	"net/http"

	"github.com/futugyousuzu/goproject/awsgolang/services"
	model "github.com/futugyousuzu/goproject/awsgolang/viewmodel"

	"github.com/futugyou/extensions"
)

func IAM(w http.ResponseWriter, r *http.Request) {
	if extensions.Cors(w, r) {
		return
	}

	operatetype := r.URL.Query().Get("type")
	switch operatetype {
	case "search":
		searchIAMData(w, r)
	}
}
func searchIAMData(w http.ResponseWriter, r *http.Request) {
	service := services.NewIAMService()

	accountId := r.URL.Query().Get("accountId")
	filter := model.IAMDataFilter{AccountId: accountId}
	datas, err := service.SearchIAMData(r.Context(), filter)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
		return
	}
	body, _ := json.Marshal(datas)
	w.Write(body)
	w.WriteHeader(200)
}
