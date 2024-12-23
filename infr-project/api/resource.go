package api

import (
	_ "github.com/joho/godotenv/autoload"

	"net/http"

	"github.com/futugyou/extensions"

	"github.com/futugyou/infr-project/controller"
	tool "github.com/futugyou/infr-project/extensions"
)

func ResourceDispatch(w http.ResponseWriter, r *http.Request) {
	// cors
	if extensions.Cors(w, r) {
		return
	}

	op := r.URL.Query().Get("optype")
	ctrl := controller.NewResourceController()
	queryctrl := controller.NewResourceQueryController()
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
		w.Write([]byte("system error"))
		w.WriteHeader(500)
		return
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
