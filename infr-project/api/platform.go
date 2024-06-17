package api

import (
	_ "github.com/joho/godotenv/autoload"

	"net/http"

	"github.com/futugyou/infr-project/api/internal"
	apiadapter "github.com/futugyou/infr-project/api_adapter"
)

func PlatformDispatch(w http.ResponseWriter, r *http.Request) {
	// cors
	if internal.CorsForVercel(w, r) {
		return
	}

	op := r.URL.Query().Get("optype")
	switch op {
	case "create":
		createPlatform(r, w)
	case "get":
		getPlatform(r, w)
	case "update":
		updatePlatform(r, w)
	case "delete":
		deletePlatform(r, w)
	case "all":
		allPlatform(r, w)
	case "hook":
		hookPlatform(r, w)
	case "project":
		createPlatformProject(r, w)
	case "prodel":
		deletePlatformProject(r, w)
	case "proup":
		updatePlatformProject(r, w)
	default:
		w.Write([]byte("system error"))
		w.WriteHeader(500)
		return
	}
}

func updatePlatformProject(r *http.Request, w http.ResponseWriter) {
	id := r.URL.Query().Get("id")
	projectId := r.URL.Query().Get("project_id")
	apiadapter.CreatePlatformProject(id, projectId, w, r)
}

func deletePlatformProject(r *http.Request, w http.ResponseWriter) {
	id := r.URL.Query().Get("id")
	projectId := r.URL.Query().Get("project_id")
	apiadapter.DeletePlatformProject(id, projectId, w, r)
}

func createPlatformProject(r *http.Request, w http.ResponseWriter) {
	id := r.URL.Query().Get("id")
	apiadapter.CreatePlatformProject(id, "", w, r)
}

func hookPlatform(r *http.Request, w http.ResponseWriter) {
	id := r.URL.Query().Get("id")
	projectId := r.URL.Query().Get("project_id")
	apiadapter.UpdatePlatformHook(id, projectId, w, r)
}

func allPlatform(r *http.Request, w http.ResponseWriter) {
	apiadapter.GetAllPlatform(w, r)
}

func deletePlatform(r *http.Request, w http.ResponseWriter) {
	id := r.URL.Query().Get("id")
	apiadapter.DeletePlatform(id, w, r)
}

func updatePlatform(r *http.Request, w http.ResponseWriter) {
	id := r.URL.Query().Get("id")
	apiadapter.UpdatePlatform(id, w, r)
}

func getPlatform(r *http.Request, w http.ResponseWriter) {
	id := r.URL.Query().Get("id")
	apiadapter.GetPlatform(id, w, r)
}

func createPlatform(r *http.Request, w http.ResponseWriter) {
	apiadapter.CreatePlatform(w, r)
}
