package scheduler

import (
	"context"
	"os"

	"github.com/futugyousuzu/identity/token"
	"github.com/robfig/cron"
)

func JwksGenerate() {
	c := cron.New()
	spec := "0 */30 * * * *"
	c.AddFunc(spec, func() {
		signed_key_id := os.Getenv("signed_key_id")
		token.NewJwksStore().CreateJwks(context.Background(), signed_key_id)
	})
	c.Start()
}
