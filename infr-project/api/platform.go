package api

import (
	"strconv"
	"strings"

	_ "github.com/joho/godotenv/autoload"

	"net/http"

	"github.com/futugyou/extensions"

	"github.com/futugyou/infr-project/controller"
	tool "github.com/futugyou/infr-project/extensions"
	viewmodels "github.com/futugyou/infr-project/view_models"
)

func PlatformDispatch(w http.ResponseWriter, r *http.Request) {
	// cors
	if extensions.Cors(w, r) {
		return
	}

	op := r.URL.Query().Get("optype")
	version := r.URL.Query().Get("version")
	if len(version) == 0 {
		version = "v1"
	}

	ctrl := controller.NewPlatformController()
	switch version {
	case "v1":
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
		case "hook_del":
			deleteHookPlatform(ctrl, r, w)
		case "project":
			createPlatformProject(ctrl, r, w)
		case "proget":
			getPlatformProject(ctrl, r, w)
		case "prodel":
			deletePlatformProject(ctrl, r, w)
		case "proup":
			updatePlatformProject(ctrl, r, w)
		case "recovery":
			recoveryPlatform(ctrl, r, w)
		default:
			w.Write([]byte("page not found"))
			w.WriteHeader(404)
		}
	case "v2":
		switch op {
		case "proget":
			getPlatformProjectV2(ctrl, r, w)
		default:
			w.Write([]byte("page not found"))
			w.WriteHeader(404)
		}
	default:
		w.Write([]byte("page not found"))
		w.WriteHeader(404)
	}
}
func recoveryPlatform(ctrl *controller.PlatformController, r *http.Request, w http.ResponseWriter) {
	if !tool.AuthForVercel(w, r) {
		return
	}

	id := r.URL.Query().Get("id")
	ctrl.RecoveryPlatform(id, w, r)
}

func updatePlatformProject(ctrl *controller.PlatformController, r *http.Request, w http.ResponseWriter) {
	id := r.URL.Query().Get("id")
	projectId := r.URL.Query().Get("project_id")
	ctrl.UpsertPlatformProject(id, projectId, w, r)
}

func deletePlatformProject(ctrl *controller.PlatformController, r *http.Request, w http.ResponseWriter) {
	if !tool.AuthForVercel(w, r) {
		return
	}

	id := r.URL.Query().Get("id")
	projectId := r.URL.Query().Get("project_id")
	ctrl.DeletePlatformProject(id, projectId, w, r)
}

func createPlatformProject(ctrl *controller.PlatformController, r *http.Request, w http.ResponseWriter) {
	if !tool.AuthForVercel(w, r) {
		return
	}

	id := r.URL.Query().Get("id")
	ctrl.UpsertPlatformProject(id, "", w, r)
}

func hookPlatform(ctrl *controller.PlatformController, r *http.Request, w http.ResponseWriter) {
	if !tool.AuthForVercel(w, r) {
		return
	}

	id := r.URL.Query().Get("id")
	projectId := r.URL.Query().Get("project_id")
	ctrl.UpdatePlatformHook(id, projectId, w, r)
}

func deleteHookPlatform(ctrl *controller.PlatformController, r *http.Request, w http.ResponseWriter) {
	if !tool.AuthForVercel(w, r) {
		return
	}

	id := r.URL.Query().Get("id")
	projectId := r.URL.Query().Get("project_id")
	hookName := r.URL.Query().Get("hook_name")
	sync, _ := strconv.ParseBool(r.URL.Query().Get("sync"))

	request := viewmodels.RemoveWebhookRequest{
		PlatformId: id,
		ProjectId:  projectId,
		HookName:   hookName,
		Sync:       sync,
	}
	ctrl.DeletePlatformHook(request, w, r)
}

func allPlatform(ctrl *controller.PlatformController, r *http.Request, w http.ResponseWriter) {
	name := r.URL.Query().Get("name")
	tags := strings.FieldsFunc(r.URL.Query().Get("tags"), func(r rune) bool {
		return r == ','
	})
	if len(tags) == 1 && tags[0] == "" {
		tags = nil
	}

	pageInt, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		pageInt = 1
	}

	sizeInt, err := strconv.Atoi(r.URL.Query().Get("size"))
	if err != nil {
		sizeInt = 100
	}

	request := viewmodels.SearchPlatformsRequest{
		Name:     name,
		Activate: extensions.StringToBoolPtr(r.URL.Query().Get("activate")),
		Tags:     tags,
		Page:     pageInt,
		Size:     sizeInt,
	}
	ctrl.SearchPlatforms(w, r, request)
}

func deletePlatform(ctrl *controller.PlatformController, r *http.Request, w http.ResponseWriter) {
	if !tool.AuthForVercel(w, r) {
		return
	}

	id := r.URL.Query().Get("id")
	ctrl.DeletePlatform(id, w, r)
}

func updatePlatform(ctrl *controller.PlatformController, r *http.Request, w http.ResponseWriter) {
	if !tool.AuthForVercel(w, r) {
		return
	}

	id := r.URL.Query().Get("id")
	ctrl.UpdatePlatform(id, w, r)
}

func getPlatform(ctrl *controller.PlatformController, r *http.Request, w http.ResponseWriter) {
	id := r.URL.Query().Get("id")
	ctrl.GetPlatform(id, w, r)
}

func getPlatformProject(ctrl *controller.PlatformController, r *http.Request, w http.ResponseWriter) {
	id := r.URL.Query().Get("id")
	projectid := r.URL.Query().Get("project_id")
	ctrl.GetPlatformProject(id, projectid, w, r)
}

func getPlatformProjectV2(ctrl *controller.PlatformController, r *http.Request, w http.ResponseWriter) {
	id := r.URL.Query().Get("id")
	projectid := r.URL.Query().Get("project_id")
	ctrl.GetPlatformProjectV2(id, projectid, w, r)
}

func createPlatform(ctrl *controller.PlatformController, r *http.Request, w http.ResponseWriter) {
	if !tool.AuthForVercel(w, r) {
		return
	}

	ctrl.CreatePlatform(w, r)
}
