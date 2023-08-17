package api

import (
	_ "github.com/joho/godotenv/autoload"

	"encoding/json"
	"net/http"

	"github.com/futugyousuzu/goproject/awsgolang/services"
	verceltool "github.com/futugyousuzu/goproject/awsgolang/vercel"
)

func ParameterCompare(w http.ResponseWriter, r *http.Request) {
	if verceltool.CrosForVercel(w, r) {
		return
	}

	if !verceltool.AuthForVercel(w, r) {
		return
	}

	parameterService := services.NewParameterService()

	sourceid := r.URL.Query().Get("sourceid")
	destid := r.URL.Query().Get("destid")

	if len(sourceid) == 0 || len(destid) == 0 {
		w.Write([]byte("sourceid or destid can not be empty"))
		w.WriteHeader(400)
		return
	}

	parameters := parameterService.CompareParameterByIDs(sourceid, destid)

	body, _ := json.Marshal(parameters)
	w.Write(body)
	w.WriteHeader(200)
}
