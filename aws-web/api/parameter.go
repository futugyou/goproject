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

func Parameter(w http.ResponseWriter, r *http.Request) {
	if extensions.Cors(w, r) {
		return
	}

	if !verceltool.AuthForVercel(w, r) {
		return
	}

	operatetype := r.URL.Query().Get("type")
	switch operatetype {
	case "compare":
		compareParameter(w, r)
	case "get":
		getParameter(w, r)
	case "getall":
		getallParameter(w, r)
	case "sync":
		syncParameter(w, r)
	}
}

func compareParameter(w http.ResponseWriter, r *http.Request) {
	parameterService := services.NewParameterService()

	sourceid := r.URL.Query().Get("sourceid")
	destid := r.URL.Query().Get("destid")

	if len(sourceid) == 0 || len(destid) == 0 {
		w.Write([]byte("sourceid or destid can not be empty"))
		w.WriteHeader(400)
		return
	}

	parameters := parameterService.CompareParameterByIDs(r.Context(), sourceid, destid)

	body, _ := json.Marshal(parameters)
	w.Write(body)
	w.WriteHeader(200)
}

func getParameter(w http.ResponseWriter, r *http.Request) {
	parameterService := services.NewParameterService()
	id := r.URL.Query().Get("id")
	parameter := parameterService.GetParameterByID(r.Context(), id)

	body, _ := json.Marshal(parameter)
	w.Write(body)
	w.WriteHeader(200)
}

func getallParameter(w http.ResponseWriter, r *http.Request) {
	parameterService := services.NewParameterService()
	var parameters []model.ParameterViewModel

	pageString := r.URL.Query().Get("page")
	limitString := r.URL.Query().Get("limit")
	alias := r.URL.Query().Get("alias")
	region := r.URL.Query().Get("region")
	key := r.URL.Query().Get("key")

	page, _ := strconv.ParseInt(pageString, 10, 64)
	limit, _ := strconv.ParseInt(limitString, 10, 64)
	if page == 0 {
		page = 1
	}

	if limit == 0 {
		limit = 10
	}

	filer := model.ParameterFilter{
		AccountAlias: alias,
		Region:       region,
		Key:          key,
	}

	paging := core.Paging{Page: page, Limit: limit}
	parameters = parameterService.GetParametersByCondition(r.Context(), paging, filer)

	body, _ := json.Marshal(parameters)
	w.Write(body)
	w.WriteHeader(200)
}

func syncParameter(w http.ResponseWriter, r *http.Request) {
	parameterService := services.NewParameterService()

	var sync model.SyncModel

	err := json.NewDecoder(r.Body).Decode(&sync)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = parameterService.SyncParameterByID(r.Context(), sync.Id)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
}
