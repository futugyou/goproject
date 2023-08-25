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
)

func EcsService(w http.ResponseWriter, r *http.Request) {
	if verceltool.CrosForVercel(w, r) {
		return
	}

	if !verceltool.AuthForVercel(w, r) {
		return
	}

	operatetype := r.URL.Query().Get("type")
	switch operatetype {
	case "getall":
		getallFuntion(w, r)
	case "get":
		getFuntion(w, r)
	case "compare":
		compareFuntion(w, r)
	}
}

func getallFuntion(w http.ResponseWriter, r *http.Request) {
	ecsService := services.NewEcsClusterService()

	pageString := r.URL.Query().Get("page")
	limitString := r.URL.Query().Get("limit")
	account_id := r.URL.Query().Get("account_id")

	page, _ := strconv.ParseInt(pageString, 10, 64)
	limit, _ := strconv.ParseInt(limitString, 10, 64)
	if page == 0 {
		page = 1
	}

	if limit == 0 {
		limit = 10
	}

	filer := model.EcsClusterFilter{
		AccountId: account_id,
	}

	paging := core.Paging{Page: page, Limit: limit}
	services, err := ecsService.GetAllServices(paging, filer)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
		return
	}

	body, _ := json.Marshal(services)
	w.Write(body)
	w.WriteHeader(200)
}

func getFuntion(w http.ResponseWriter, r *http.Request) {
	ecsService := services.NewEcsClusterService()

	id := r.URL.Query().Get("id")
	service, err := ecsService.GetServiceDetailById(id)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
		return
	}

	body, _ := json.Marshal(service)
	w.Write(body)
	w.WriteHeader(200)
}

func compareFuntion(w http.ResponseWriter, r *http.Request) {
	ecsService := services.NewEcsClusterService()
	var compare model.EcsTaskCompare

	err := json.NewDecoder(r.Body).Decode(&compare)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := ecsService.CompareTaskDefinitions(compare)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
		return
	}

	body, _ := json.Marshal(result)
	w.Write(body)
	w.WriteHeader(200)
}
