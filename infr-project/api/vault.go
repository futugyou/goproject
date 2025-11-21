package api

import (
	"strconv"
	"strings"

	_ "github.com/joho/godotenv/autoload"

	"net/http"

	"github.com/futugyou/extensions"

	"github.com/futugyou/infr-project/controller"
	tool "github.com/futugyou/infr-project/extensions"

	"github.com/futugyou/vaultservice/viewmodel"
)

func VaultDispatch(w http.ResponseWriter, r *http.Request) {
	// cors
	if extensions.Cors(w, r) {
		return
	}

	op := r.URL.Query().Get("optype")
	version := r.URL.Query().Get("version")
	if len(version) == 0 {
		version = "v1"
	}

	ctrl := controller.NewVaultController()
	switch version {
	case "v1":
		switch op {
		case "batch":
			createVaults(ctrl, r, w)
		case "single":
			createVault(ctrl, r, w)
		case "show":
			showVault(ctrl, r, w)
		case "get":
			getVault(ctrl, r, w)
		case "update":
			updateVault(ctrl, r, w)
		case "delete":
			deleteVault(ctrl, r, w)
		case "import":
			importVault(ctrl, r, w)
		default:
			w.Write([]byte("page not found"))
			w.WriteHeader(404)
		}
	default:
		w.Write([]byte("page not found"))
		w.WriteHeader(404)
	}
}

func getVault(ctrl *controller.VaultController, r *http.Request, w http.ResponseWriter) {
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
	request := viewmodel.SearchVaultsRequest{
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

func deleteVault(ctrl *controller.VaultController, r *http.Request, w http.ResponseWriter) {
	if !tool.AuthForVercel(w, r) {
		return
	}

	id := r.URL.Query().Get("id")
	ctrl.DeleteVault(w, r, id)
}

func updateVault(ctrl *controller.VaultController, r *http.Request, w http.ResponseWriter) {
	if !tool.AuthForVercel(w, r) {
		return
	}

	id := r.URL.Query().Get("id")
	ctrl.ChangeVault(w, r, id)
}

func showVault(ctrl *controller.VaultController, r *http.Request, w http.ResponseWriter) {
	if !tool.AuthForVercel(w, r) {
		return
	}

	id := r.URL.Query().Get("id")
	ctrl.ShowVaultRawValue(w, r, id)
}

func createVaults(ctrl *controller.VaultController, r *http.Request, w http.ResponseWriter) {
	if !tool.AuthForVercel(w, r) {
		return
	}

	ctrl.CreateVaults(w, r)
}

func createVault(ctrl *controller.VaultController, r *http.Request, w http.ResponseWriter) {
	if !tool.AuthForVercel(w, r) {
		return
	}

	ctrl.CreateVault(w, r)
}

func importVault(ctrl *controller.VaultController, r *http.Request, w http.ResponseWriter) {
	if !tool.AuthForVercel(w, r) {
		return
	}

	ctrl.ImportVaults(w, r)
}
