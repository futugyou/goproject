package api

import (
	_ "github.com/joho/godotenv/autoload"

	"net/http"

	"github.com/futugyou/infr-project/api/internal"
	"github.com/futugyou/infr-project/controller"
)

func ProjectDispatch(w http.ResponseWriter, r *http.Request) {
	// cors
	if internal.CorsForVercel(w, r) {
		return
	}

	op := r.URL.Query().Get("optype")
	ctrl := controller.NewController()
	switch op {
	case "create":
		createProject(ctrl, r, w)
	case "get":
		getProject(ctrl, r, w)
	case "update":
		updateProject(ctrl, r, w)
	case "all":
		allProject(ctrl, r, w)
	case "platform":
		updateProjectPlatform(ctrl, r, w)
	case "design":
		updateProjectDesign(ctrl, r, w)
	default:
		w.Write([]byte("system error"))
		w.WriteHeader(500)
		return
	}
}

func allProject(ctrl *controller.Controller, r *http.Request, w http.ResponseWriter) {
	ctrl.GetAllProject(w, r)
}

func updateProject(ctrl *controller.Controller, r *http.Request, w http.ResponseWriter) {
	id := r.URL.Query().Get("id")
	ctrl.UpdateProject(id, w, r)
}

func getProject(ctrl *controller.Controller, r *http.Request, w http.ResponseWriter) {
	id := r.URL.Query().Get("id")
	ctrl.GetProject(id, w, r)
}

func createProject(ctrl *controller.Controller, r *http.Request, w http.ResponseWriter) {
	ctrl.CreateProject(w, r)
}

func updateProjectPlatform(ctrl *controller.Controller, r *http.Request, w http.ResponseWriter) {
	id := r.URL.Query().Get("id")
	ctrl.UpdateProjectPlatform(id, w, r)
}

func updateProjectDesign(ctrl *controller.Controller, r *http.Request, w http.ResponseWriter) {
	id := r.URL.Query().Get("id")
	ctrl.UpdateProjectDesign(id, w, r)
}
