package api

import (
	_ "github.com/joho/godotenv/autoload"

	"encoding/json"
	"net/http"

	"github.com/futugyousuzu/goproject/awsgolang/services"
	verceltool "github.com/futugyousuzu/goproject/awsgolang/vercel"
)

func Get(w http.ResponseWriter, r *http.Request) {
	if verceltool.CrosForVercel(w, r) {
		return
	}

	if !verceltool.AuthForVercel(w, r) {
		return
	}

	regions, err := services.GetAllRegionInCurrentAccount()

	if err != nil {
		body, _ := json.Marshal(err)
		w.Write(body)
		w.WriteHeader(500)
		return
	}

	body, _ := json.Marshal(regions)
	w.Write(body)
	w.WriteHeader(200)
}
