package api

import (
	_ "github.com/joho/godotenv/autoload"

	"encoding/json"
	"net/http"

	"github.com/futugyousuzu/goproject/awsgolang/services"
	verceltool "github.com/futugyousuzu/goproject/awsgolang/vercel"
)

func AccountGet(w http.ResponseWriter, r *http.Request) {
	if verceltool.CrosForVercel(w, r) {
		return
	}

	if !verceltool.AuthForVercel(w, r) {
		return
	}

	accountService := services.NewAccountService()

	id := r.URL.Query().Get("id")

	account := accountService.GetAccountByID(id)

	body, _ := json.Marshal(account)
	w.Write(body)
	w.WriteHeader(200)
}
