package api

import (
	_ "github.com/joho/godotenv/autoload"

	"encoding/json"
	"net/http"

	"github.com/futugyousuzu/goproject/awsgolang/services"
	verceltool "github.com/futugyousuzu/goproject/awsgolang/vercel"
	model "github.com/futugyousuzu/goproject/awsgolang/viewmodel"

	"github.com/futugyou/extensions"
)

func KeyValue(w http.ResponseWriter, r *http.Request) {
	if extensions.Cors(w, r) {
		return
	}

	operatetype := r.URL.Query().Get("type")
	switch operatetype {
	case "getall":
		getAllKeyValue(w, r)
	case "get":
		getKeyValue(w, r)
	case "create":
		createKeyValue(w, r)
	case "awsconfig":
		getResourceGraph(w, r)
	case "newconfig":
		getNewResourceGraph(w, r)
	}
}

func createKeyValue(w http.ResponseWriter, r *http.Request) {
	if !verceltool.AuthForVercel(w, r) {
		return
	}

	var keyValue model.KeyValue
	err := json.NewDecoder(r.Body).Decode(&keyValue)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	keyValueService := services.NewKeyValueService()
	keyValueService.CreateKeyValue(r.Context(), keyValue.Key, keyValue.Value)
	w.WriteHeader(200)
}

func getKeyValue(w http.ResponseWriter, r *http.Request) {
	keyValueService := services.NewKeyValueService()
	key := r.URL.Query().Get("key")
	keyValue := keyValueService.GetValueByKey(r.Context(), key)
	body, _ := json.Marshal(keyValue)
	w.Write(body)
	w.WriteHeader(200)
}

func getAllKeyValue(w http.ResponseWriter, r *http.Request) {
	keyValueService := services.NewKeyValueService()
	keyValues := keyValueService.GetAllKeyValues(r.Context())
	body, _ := json.Marshal(keyValues)
	w.Write(body)
	w.WriteHeader(200)
}

func getResourceGraph(w http.ResponseWriter, r *http.Request) {
	config := services.NewAwsConfigService()
	res := config.GetResourceGraph(r.Context())
	body, _ := json.Marshal(res)
	w.Write(body)
	w.WriteHeader(200)
}

func getNewResourceGraph(w http.ResponseWriter, r *http.Request) {
	config := services.NewAwsConfigServiceWithTableNames("awsconfig_new", "awsconfig_new_relationship")
	res := config.GetResourceGraph(r.Context())
	body, _ := json.Marshal(res)
	w.Write(body)
	w.WriteHeader(200)
}
