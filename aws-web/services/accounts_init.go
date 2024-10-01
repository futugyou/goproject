package services

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"log"

	"github.com/futugyousuzu/goproject/awsgolang/entity"
	"github.com/google/uuid"
)

func (a *AccountService) AccountInit(ctx context.Context) {
	result := make([]entity.AccountEntity, 0)
	var accounts []byte
	var err error

	if accounts, err = os.ReadFile("./data/accounts.json"); err != nil {
		log.Println(err)
		return
	}

	if err = json.Unmarshal(accounts, &result); err != nil {
		log.Println(err)
		return
	}

	if err = a.repository.DeleteAll(ctx); err != nil {
		log.Println(err)
		return
	}

	for i := 0; i < len(result); i++ {
		result[i].Id = uuid.New().String()
		result[i].CreatedAt = time.Now().Unix()
	}

	if err = a.repository.InsertMany(ctx, result); err != nil {
		log.Println(err)
	}
}
