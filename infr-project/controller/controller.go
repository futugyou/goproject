package controller

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/futugyou/infr-project/extensions"
)

type Controller struct {
}

func NewController() *Controller {
	return &Controller{}
}

func handleError(w http.ResponseWriter, err error, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(err.Error())
}

func writeJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func handleRequest[S any, T any](
	w http.ResponseWriter,
	r *http.Request,
	createService func(ctx context.Context) (S, error),
	handler func(ctx context.Context, service S, req T) (interface{}, error),
) {
	handleRequestUseSpecValidate(
		w,
		r,
		createService,
		func(v *validator.Validate, req *T) error {
			return v.Struct(req)
		},
		handler,
	)
}

func handleRequestUseSpecValidate[S any, T any](
	w http.ResponseWriter,
	r *http.Request,
	createService func(ctx context.Context) (S, error),
	validor func(v *validator.Validate, req *T) error,
	handler func(ctx context.Context, service S, req T) (interface{}, error),
) {
	ctx := r.Context()
	service, err := createService(ctx)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	var req *T
	if r.Method != http.MethodGet && r.Method != http.MethodDelete && r.Method != http.MethodOptions {
		if r.Body != nil {
			req = new(T)
			if err := json.NewDecoder(r.Body).Decode(req); err != nil && err != io.EOF {
				handleError(w, err, 400)
				return
			}
		}
	}

	if req == nil {
		req = new(T)
	}

	if err := validor(extensions.Validate, req); err != nil {
		handleError(w, err, 400)
		return
	}

	res, err := handler(ctx, service, *req)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	writeJSONResponse(w, res, 200)
}
