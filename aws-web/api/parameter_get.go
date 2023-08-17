package api

import (
	_ "github.com/joho/godotenv/autoload"

	"encoding/json"
	"net/http"

	"github.com/futugyousuzu/goproject/awsgolang/services"
	verceltool "github.com/futugyousuzu/goproject/awsgolang/vercel"
)

func ParameterGet(w http.ResponseWriter, r *http.Request) {
	if verceltool.CrosForVercel(w, r) {
		return
	}

	if !verceltool.AuthForVercel(w, r) {
		return
	}

	parameterService := services.NewParameterService()
	id := r.URL.Query().Get("id")
	parameter := parameterService.GetParameterByID(id)

	body, _ := json.Marshal(parameter)
	w.Write(body)
	w.WriteHeader(200)
}
