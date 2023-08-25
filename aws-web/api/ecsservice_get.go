package api

import (
	_ "github.com/joho/godotenv/autoload"

	"encoding/json"
	"net/http"

	"github.com/futugyousuzu/goproject/awsgolang/services"
	verceltool "github.com/futugyousuzu/goproject/awsgolang/vercel"
)

func EcsServiceGet(w http.ResponseWriter, r *http.Request) {
	if verceltool.CrosForVercel(w, r) {
		return
	}

	if !verceltool.AuthForVercel(w, r) {
		return
	}

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
