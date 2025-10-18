package api

import (
	"strconv"

	_ "github.com/joho/godotenv/autoload"

	"encoding/json"
	"net/http"

	"github.com/futugyousuzu/goproject/awsgolang/core"
	"github.com/futugyousuzu/goproject/awsgolang/services"
	verceltool "github.com/futugyousuzu/goproject/awsgolang/vercel"
	model "github.com/futugyousuzu/goproject/awsgolang/viewmodel"

	"github.com/futugyou/extensions"
)

func Account(w http.ResponseWriter, r *http.Request) {
	if extensions.Cors(w, r) {
		return
	}

	operatetype := r.URL.Query().Get("type")
	switch operatetype {
	case "getall":
		getAllAccount(w, r)
	case "get":
		getAccount(w, r)
	case "create":
		createAccount(w, r)
	case "delete":
		deleteAccount(w, r)
	case "update":
		updateAccount(w, r)
	}
}

func createAccount(w http.ResponseWriter, r *http.Request) {
	if !verceltool.AuthForVercel(w, r) {
		return
	}
	var account model.UserAccount

	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	accountService := services.NewAccountService()

	err = accountService.CreateAccount(r.Context(), account)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
}

func deleteAccount(w http.ResponseWriter, r *http.Request) {
	if !verceltool.AuthForVercel(w, r) {
		return
	}

	id := r.URL.Query().Get("id")
	accountService := services.NewAccountService()

	err := accountService.DeleteAccount(r.Context(), id)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
}

func getAccount(w http.ResponseWriter, r *http.Request) {
	accountService := services.NewAccountService()

	id := r.URL.Query().Get("id")

	account := accountService.GetAccountByID(r.Context(), id)

	body, _ := json.Marshal(account)
	w.Write(body)
	w.WriteHeader(200)
}

func getAllAccount(w http.ResponseWriter, r *http.Request) {
	accountService := services.NewAccountService()
	var accounts []model.UserAccount

	pageString := r.URL.Query().Get("page")
	limitString := r.URL.Query().Get("limit")
	page, _ := strconv.ParseInt(pageString, 10, 64)
	limit, _ := strconv.ParseInt(limitString, 10, 64)

	if page != 0 && limit != 0 {
		paging := core.Paging{Page: page, Limit: limit}
		accounts = accountService.GetAccountsByPaging(r.Context(), paging)
	} else {
		accounts = accountService.GetAllAccounts(r.Context())
	}

	body, _ := json.Marshal(accounts)
	w.Write(body)
	w.WriteHeader(200)
}

func updateAccount(w http.ResponseWriter, r *http.Request) {
	if !verceltool.AuthForVercel(w, r) {
		return
	}

	var account model.UserAccount
	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	accountService := services.NewAccountService()

	err = accountService.UpdateAccount(r.Context(), account)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
}
