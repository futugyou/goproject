package server

import (
	"context"
	"encoding/json"
	"os"

	"github.com/futugyousuzu/identity/mongo"
	"github.com/go-oauth2/oauth2/v4/models"
)

func initClient(stor *mongo.ClientStore) {
	clientstring := os.Getenv("init_clients")
	clients := make([]models.Client, 0)
	json.Unmarshal([]byte(clientstring), &clients)
	for _, client := range clients {
		stor.Set(context.Background(), &client)
	}
}
