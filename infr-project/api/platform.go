package api

import (
	_ "github.com/joho/godotenv/autoload"

	"net/http"

	"github.com/futugyou/infr-project/controller"
	"github.com/futugyou/infr-project/extensions"
)

func PlatformDispatch(w http.ResponseWriter, r *http.Request) {
	// cors
	if extensions.CorsForVercel(w, r) {
		return
	}

	op := r.URL.Query().Get("optype")
	ctrl := controller.NewController()
	switch op {
	case "create":
		createPlatform(ctrl, r, w)
	case "get":
		getPlatform(ctrl, r, w)
	case "update":
		updatePlatform(ctrl, r, w)
	case "delete":
		deletePlatform(ctrl, r, w)
	case "all":
		allPlatform(ctrl, r, w)
	case "hook":
		hookPlatform(ctrl, r, w)
	case "project":
		createPlatformProject(ctrl, r, w)
	case "prodel":
		deletePlatformProject(ctrl, r, w)
	case "proup":
		updatePlatformProject(ctrl, r, w)
	default:
		w.Write([]byte("system error"))
		w.WriteHeader(500)
		return
	}
}

func updatePlatformProject(ctrl *controller.Controller, r *http.Request, w http.ResponseWriter) {
	id := r.URL.Query().Get("id")
	projectId := r.URL.Query().Get("project_id")
	ctrl.CreatePlatformProject(id, projectId, w, r)
}

func deletePlatformProject(ctrl *controller.Controller, r *http.Request, w http.ResponseWriter) {
	id := r.URL.Query().Get("id")
	projectId := r.URL.Query().Get("project_id")
	ctrl.DeletePlatformProject(id, projectId, w, r)
}

func createPlatformProject(ctrl *controller.Controller, r *http.Request, w http.ResponseWriter) {
	id := r.URL.Query().Get("id")
	ctrl.CreatePlatformProject(id, "", w, r)
}

func hookPlatform(ctrl *controller.Controller, r *http.Request, w http.ResponseWriter) {
	id := r.URL.Query().Get("id")
	projectId := r.URL.Query().Get("project_id")
	ctrl.UpdatePlatformHook(id, projectId, w, r)
}

func allPlatform(ctrl *controller.Controller, r *http.Request, w http.ResponseWriter) {
	ctrl.GetAllPlatform(w, r)
}

func deletePlatform(ctrl *controller.Controller, r *http.Request, w http.ResponseWriter) {
	id := r.URL.Query().Get("id")
	ctrl.DeletePlatform(id, w, r)
}

func updatePlatform(ctrl *controller.Controller, r *http.Request, w http.ResponseWriter) {
	id := r.URL.Query().Get("id")
	ctrl.UpdatePlatform(id, w, r)
}

func getPlatform(ctrl *controller.Controller, r *http.Request, w http.ResponseWriter) {
	id := r.URL.Query().Get("id")
	ctrl.GetPlatform(id, w, r)
}

func createPlatform(ctrl *controller.Controller, r *http.Request, w http.ResponseWriter) {
	ctrl.CreatePlatform(w, r)
}
