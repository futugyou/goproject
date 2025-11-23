package api

import (
	_ "github.com/joho/godotenv/autoload"

	"net/http"

	"github.com/futugyou/extensions"

	"github.com/futugyou/infr-project/controller"
	tool "github.com/futugyou/infr-project/extensions"
)

func ResourceDispatch(w http.ResponseWriter, r *http.Request) {
	r = passingToken(r)
	// cors
	if extensions.Cors(w, r) {
		return
	}

	op := r.URL.Query().Get("optype")
	version := r.URL.Query().Get("version")
	if len(version) == 0 {
		version = "v1"
	}

	ctrl := controller.NewResourceController()
	queryctrl := controller.NewResourceQueryController()
	switch version {
	case "v1":
		switch op {
		case "create":
			createResource(ctrl, r, w)
		case "get":
			getResource(queryctrl, r, w)
		case "update":
			updateResource(ctrl, r, w)
		case "delete":
			deleteResource(ctrl, r, w)
		case "history":
			historyResource(ctrl, r, w)
		case "all":
			allResource(queryctrl, r, w)
		default:
			w.Write([]byte("page not found"))
			w.WriteHeader(404)
		}
	default:
		w.Write([]byte("page not found"))
		w.WriteHeader(404)
	}
}

func historyResource(ctrl *controller.ResourceController, r *http.Request, w http.ResponseWriter) {
	id := r.URL.Query().Get("id")
	ctrl.GetResourceHistory(id, w, r)
}

func allResource(ctrl *controller.ResourceQueryController, r *http.Request, w http.ResponseWriter) {
	ctrl.GetAllResource(w, r)
}

func deleteResource(ctrl *controller.ResourceController, r *http.Request, w http.ResponseWriter) {
	if !tool.AuthForVercel(w, r) {
		return
	}

	id := r.URL.Query().Get("id")
	ctrl.DeleteResource(id, w, r)
}

func updateResource(ctrl *controller.ResourceController, r *http.Request, w http.ResponseWriter) {
	if !tool.AuthForVercel(w, r) {
		return
	}

	id := r.URL.Query().Get("id")
	ctrl.UpdateResource(id, w, r)
}

func getResource(ctrl *controller.ResourceQueryController, r *http.Request, w http.ResponseWriter) {
	id := r.URL.Query().Get("id")
	ctrl.GetResource(id, w, r)
}

func createResource(ctrl *controller.ResourceController, r *http.Request, w http.ResponseWriter) {
	if !tool.AuthForVercel(w, r) {
		return
	}

	ctrl.CreateResource(w, r)
}
