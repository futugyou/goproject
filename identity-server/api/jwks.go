package api

import (
	"net/http"

	"github.com/futugyousuzu/identity-server/operate"
	"github.com/futugyousuzu/identity-server/service"
)

type JwksApi struct {
	*operate.Operator
}

func NewJwksApi(operator *operate.Operator) *JwksApi {
	return &JwksApi{
		operator,
	}
}

func (j *JwksApi) Jwks(w http.ResponseWriter, r *http.Request) {
	service := service.NewJwksService(j.Operator)
	result, _ := service.GetPublicJwks(r.Context())
	w.Write([]byte(result))
	w.WriteHeader(200)
}
