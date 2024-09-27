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

func VaultDispatch(w http.ResponseWriter, r *http.Request) {
	// cors
	if extensions.Cors(w, r) {
		return
	}

	op := r.URL.Query().Get("optype")
	ctrl := controller.NewController()
	switch op {
	case "create":
		createVault(ctrl, r, w)
	case "show":
		showVault(ctrl, r, w)
	case "get":
		getVault(ctrl, r, w)
	case "update":
		updateVault(ctrl, r, w)
	case "delete":
		deleteVault(ctrl, r, w)
	default:
		w.Write([]byte("system error"))
		w.WriteHeader(500)
		return
	}
}

func getVault(ctrl *controller.Controller, r *http.Request, w http.ResponseWriter) {
	key := r.URL.Query().Get("key")
	storageMedia := r.URL.Query().Get("storage_media")
	tags := strings.FieldsFunc(r.URL.Query().Get("tags"), func(r rune) bool {
		return r == ','
	})
	if len(tags) == 1 && tags[0] == "" {
		tags = nil
	}
	typeIdentity := r.URL.Query().Get("type_identity")
	vaultType := r.URL.Query().Get("vault_type")

	page := r.URL.Query().Get("page")
	size := r.URL.Query().Get("size")

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 1
	}

	sizeInt, err := strconv.Atoi(size)
	if err != nil {
		sizeInt = 100
	}
	request := viewmodels.SearchVaultsRequest{
		Key:          key,
		StorageMedia: storageMedia,
		VaultType:    vaultType,
		TypeIdentity: typeIdentity,
		Tags:         tags,
		Page:         pageInt,
		Size:         sizeInt,
	}
	ctrl.SearchVaults(w, r, request)
}

func deleteVault(ctrl *controller.Controller, r *http.Request, w http.ResponseWriter) {
	if !tool.AuthForVercel(w, r) {
		return
	}

	id := r.URL.Query().Get("id")
	ctrl.DeleteVault(w, r, id)
}

func updateVault(ctrl *controller.Controller, r *http.Request, w http.ResponseWriter) {
	if !tool.AuthForVercel(w, r) {
		return
	}

	id := r.URL.Query().Get("id")
	ctrl.ChangeVault(w, r, id)
}

func showVault(ctrl *controller.Controller, r *http.Request, w http.ResponseWriter) {
	if !tool.AuthForVercel(w, r) {
		return
	}

	id := r.URL.Query().Get("id")
	ctrl.ShowVaultRawValue(w, r, id)
}

func createVault(ctrl *controller.Controller, r *http.Request, w http.ResponseWriter) {
	if !tool.AuthForVercel(w, r) {
		return
	}

	ctrl.CreateVaults(w, r)
}
