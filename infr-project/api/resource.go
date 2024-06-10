package api

import (
	_ "github.com/joho/godotenv/autoload"

	"net/http"

	"github.com/futugyou/infr-project/api/internal"
	apiadapter "github.com/futugyou/infr-project/api_adapter"
)

func ResourceDispatch(w http.ResponseWriter, r *http.Request) {
	// cors
	if internal.CorsForVercel(w, r) {
		return
	}

	op := r.URL.Query().Get("optype")
	switch op {
	case "create":
		createResource(r, w)
	case "get":
		getResource(r, w)
	case "update":
		updateResource(r, w)
	case "delete":
		deleteResource(r, w)
	case "history":
		historyResource(r, w)
	case "all":
		allResource(r, w)
	default:
		w.Write([]byte("system error"))
		w.WriteHeader(500)
		return
	}
}

func historyResource(r *http.Request, w http.ResponseWriter) {
	id := r.URL.Query().Get("id")
	apiadapter.GetResourceHistory(id, w, r)
}

func allResource(r *http.Request, w http.ResponseWriter) {
	apiadapter.GetAllResource(w, r)
}

func deleteResource(r *http.Request, w http.ResponseWriter) {
	id := r.URL.Query().Get("id")
	apiadapter.DeleteResource(id, w, r)
}

func updateResource(r *http.Request, w http.ResponseWriter) {
	id := r.URL.Query().Get("id")
	apiadapter.UpdateResource(id, w, r)
}

func getResource(r *http.Request, w http.ResponseWriter) {
	id := r.URL.Query().Get("id")
	apiadapter.GetResource(id, w, r)
}

func createResource(r *http.Request, w http.ResponseWriter) {
	apiadapter.CreateResource(w, r)
}
