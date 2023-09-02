package api

import (
	_ "github.com/joho/godotenv/autoload"

	"encoding/json"
	"net/http"

	"github.com/futugyousuzu/goproject/awsgolang/services"
	verceltool "github.com/futugyousuzu/goproject/awsgolang/vercel"
	model "github.com/futugyousuzu/goproject/awsgolang/viewmodel"
)

func KeyValue(w http.ResponseWriter, r *http.Request) {
	if verceltool.CrosForVercel(w, r) {
		return
	}

	if !verceltool.AuthForVercel(w, r) {
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
	}
}

func createKeyValue(w http.ResponseWriter, r *http.Request) {
	var KeyValue model.KeyValue

	err := json.NewDecoder(r.Body).Decode(&KeyValue)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	keyValueService := services.NewKeyValueService()
	keyValueService.CreateKeyValue(KeyValue.Key, KeyValue.Value)
	w.WriteHeader(200)
}

func getKeyValue(w http.ResponseWriter, r *http.Request) {
	keyValueService := services.NewKeyValueService()
	key := r.URL.Query().Get("key")
	KeyValue := keyValueService.GetValueByKey(key)
	body, _ := json.Marshal(KeyValue)
	w.Write(body)
	w.WriteHeader(200)
}

func getAllKeyValue(w http.ResponseWriter, r *http.Request) {
	keyValueService := services.NewKeyValueService()
	KeyValues := keyValueService.GetAllKeyValues()
	body, _ := json.Marshal(KeyValues)
	w.Write(body)
	w.WriteHeader(200)
}
