package api

import (
	_ "github.com/joho/godotenv/autoload"

	"net/http"

	"github.com/futugyou/infr-project/api/internal"
	apiadapter "github.com/futugyou/infr-project/api_adapter"
)

func ProjectDispatch(w http.ResponseWriter, r *http.Request) {
	// cors
	if internal.CorsForVercel(w, r) {
		return
	}

	op := r.URL.Query().Get("optype")
	switch op {
	case "create":
		createProject(r, w)
	case "get":
		getProject(r, w)
	case "update":
		updateProject(r, w)
	case "all":
		allProject(r, w)
	case "platform":
		updateProjectPlatform(r, w)
	case "design":
		updateProjectDesign(r, w)
	default:
		w.Write([]byte("system error"))
		w.WriteHeader(500)
		return
	}
}

func allProject(r *http.Request, w http.ResponseWriter) {
	apiadapter.GetAllProject(w, r)
}

func updateProject(r *http.Request, w http.ResponseWriter) {
	id := r.URL.Query().Get("id")
	apiadapter.UpdateProject(id, w, r)
}

func getProject(r *http.Request, w http.ResponseWriter) {
	id := r.URL.Query().Get("id")
	apiadapter.GetProject(id, w, r)
}

func createProject(r *http.Request, w http.ResponseWriter) {
	apiadapter.CreateProject(w, r)
}

func updateProjectPlatform(r *http.Request, w http.ResponseWriter) {
	id := r.URL.Query().Get("id")
	apiadapter.UpdateProjectPlatform(id, w, r)
}

func updateProjectDesign(r *http.Request, w http.ResponseWriter) {
	id := r.URL.Query().Get("id")
	apiadapter.UpdateProjectDesign(id, w, r)
}
