package api

import (
	_ "github.com/joho/godotenv/autoload"

	"net/http"

	"github.com/futugyousuzu/goproject/awsgolang/services"
	verceltool "github.com/futugyousuzu/goproject/awsgolang/vercel"
)

func AccountDelete(w http.ResponseWriter, r *http.Request) {
	if verceltool.CrosForVercel(w, r) {
		return
	}

	if !verceltool.AuthForVercel(w, r) {
		return
	}

	id := r.URL.Query().Get("id")

	accountService := services.NewAccountService()

	err := accountService.DeleteAccount(id)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
}
