package api

import (
	"strconv"

	_ "github.com/joho/godotenv/autoload"

	"encoding/json"
	"net/http"

	"github.com/futugyousuzu/goproject/awsgolang/core"
	"github.com/futugyousuzu/goproject/awsgolang/services"
	verceltool "github.com/futugyousuzu/goproject/awsgolang/vercel"
)

func AccountGetall(w http.ResponseWriter, r *http.Request) {
	if verceltool.CrosForVercel(w, r) {
		return
	}

	if !verceltool.AuthForVercel(w, r) {
		return
	}

	accountService := services.NewAccountService()
	var accounts []services.UserAccount

	pageString := r.URL.Query().Get("page")
	limitString := r.URL.Query().Get("limit")
	page, _ := strconv.ParseInt(pageString, 10, 64)
	limit, _ := strconv.ParseInt(limitString, 10, 64)

	if page != 0 && limit != 0 {
		paging := core.Paging{Page: page, Limit: limit}
		accounts = accountService.GetAccountsByPaging(paging)
	} else {
		accounts = accountService.GetAllAccounts()
	}

	body, _ := json.Marshal(accounts)
	w.Write(body)
	w.WriteHeader(200)
}
