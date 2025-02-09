package api

import (
	"strconv"

	_ "github.com/joho/godotenv/autoload"

	"net/http"

	"github.com/futugyou/extensions"

	"github.com/futugyou/infr-project/controller"
	tool "github.com/futugyou/infr-project/extensions"
)

func ProjectDispatch(w http.ResponseWriter, r *http.Request) {
	// cors
	if extensions.Cors(w, r) {
		return
	}

	op := r.URL.Query().Get("optype")
	version := r.URL.Query().Get("version")
	if len(version) == 0 {
		version = "v1"
	}

	ctrl := controller.NewProjectController()
	switch version {
	case "v1":
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
			w.Write([]byte("page not found"))
			w.WriteHeader(404)
		}
	default:
		w.Write([]byte("page not found"))
		w.WriteHeader(404)
	}
}

func allProject(ctrl *controller.ProjectController, r *http.Request, w http.ResponseWriter) {
	pageStr := r.URL.Query().Get("page")
	sizeStr := r.URL.Query().Get("size")
	var page *int
	if p, err := strconv.Atoi(pageStr); err != nil && p > 0 {
		page = &p
	}
	var size *int
	if p, err := strconv.Atoi(sizeStr); err != nil && p > 0 {
		size = &p
	}
	ctrl.GetAllProject(w, r, page, size)
}

func updateProject(ctrl *controller.ProjectController, r *http.Request, w http.ResponseWriter) {
	if !tool.AuthForVercel(w, r) {
		return
	}

	id := r.URL.Query().Get("id")
	ctrl.UpdateProject(id, w, r)
}

func getProject(ctrl *controller.ProjectController, r *http.Request, w http.ResponseWriter) {
	id := r.URL.Query().Get("id")
	ctrl.GetProject(id, w, r)
}

func createProject(ctrl *controller.ProjectController, r *http.Request, w http.ResponseWriter) {
	if !tool.AuthForVercel(w, r) {
		return
	}

	ctrl.CreateProject(w, r)
}

func updateProjectPlatform(ctrl *controller.ProjectController, r *http.Request, w http.ResponseWriter) {
	if !tool.AuthForVercel(w, r) {
		return
	}

	id := r.URL.Query().Get("id")
	ctrl.UpdateProjectPlatform(id, w, r)
}

func updateProjectDesign(ctrl *controller.ProjectController, r *http.Request, w http.ResponseWriter) {
	if !tool.AuthForVercel(w, r) {
		return
	}

	id := r.URL.Query().Get("id")
	ctrl.UpdateProjectDesign(id, w, r)
}
