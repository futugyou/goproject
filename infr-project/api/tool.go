package api

import (
	"net/http"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"github.com/redis/go-redis/v9"

	"github.com/futugyou/extensions"

	"github.com/futugyou/infr-project/controller"
)

func ToolsDispatch(w http.ResponseWriter, r *http.Request) {
	// cors
	if extensions.Cors(w, r) {
		return
	}

	op := r.URL.Query().Get("optype")
	ctrl := controller.NewController()
	switch op {
	case "redis":
		redistool(ctrl, r, w)
	default:
		w.Write([]byte("system error"))
		w.WriteHeader(500)
		return
	}
}

func redistool(_ *controller.Controller, r *http.Request, w http.ResponseWriter) {
	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		w.Write([]byte("linkMsg:" + err.Error()))
		w.WriteHeader(500)
		return
	}
	opt.MaxRetries = 3
	opt.DialTimeout = 10 * time.Second
	opt.ReadTimeout = -1
	opt.WriteTimeout = -1
	opt.DB = 0

	client := redis.NewClient(opt)

	ctx := r.Context()

	err = client.Set(ctx, "foo", "bar", 0).Err()
	if err != nil {
		w.Write([]byte("WriteMsg:" + err.Error()))
		w.WriteHeader(500)
		return
	}

	val, err := client.Get(ctx, "foo").Result()
	if err != nil {
		w.Write([]byte("ReadMsg:" + err.Error()))
		w.WriteHeader(500)
		return
	}

	w.Write([]byte("ResultMsg:" + val))
	w.WriteHeader(200)
}
