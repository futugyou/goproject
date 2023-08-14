package api

import (
	_ "github.com/joho/godotenv/autoload"

	"net/http"

	"github.com/futugyousuzu/goproject/awsgolang/services"
)

func Cron(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("key") != "E86ZaWHjHnRqn" {
		w.WriteHeader(404)
		return
	}

	parameterService := services.NewParameterService()
	parameterService.SyncAllParameter()
	w.WriteHeader(200)
}
