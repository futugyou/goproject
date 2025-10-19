package services

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/iam"
	tool "github.com/futugyou/extensions"

	"github.com/futugyousuzu/goproject/awsgolang/awsenv"
	model "github.com/futugyousuzu/goproject/awsgolang/viewmodel"
)

type IAMService struct {
}

func NewIAMService() *IAMService {
	return &IAMService{}
}

func (s *IAMService) SearchIAMData(ctx context.Context, filter model.IAMDataFilter) ([]model.IAMData, error) {
	accountService := NewAccountService()
	account := accountService.GetAccountByID(ctx, filter.AccountId, false)
	if !account.Valid {
		return nil, fmt.Errorf("account %s is expired", account.Id)
	}
	awsenv.CfgWithProfileAndRegion(account.AccessKeyId, account.SecretAccessKey, account.Region)
	svc := iam.NewFromConfig(awsenv.Cfg)

	iamlist, err := svc.ListUsers(ctx, &iam.ListUsersInput{})
	if err != nil {
		return nil, err
	}

	var result = []model.IAMData{}
	for _, user := range iamlist.Users {
		data := model.IAMData{
			UserName:   *user.UserName,
			CreateDate: *user.CreateDate,
			LastUsed:   user.PasswordLastUsed,
			Keys:       []model.IAMDataKey{},
		}
		input := &iam.ListAccessKeysInput{
			UserName: user.UserName,
		}

		listkeys, err := svc.ListAccessKeys(ctx, input)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		for _, key := range listkeys.AccessKeyMetadata {
			if key.AccessKeyId != nil {
				data.Keys = append(data.Keys, model.IAMDataKey{
					Key:        tool.MaskString(*key.AccessKeyId, 5, 0.5),
					Status:     string(key.Status),
					CreateDate: *key.CreateDate,
				})
			}
		}

		result = append(result, data)
	}

	return result, nil
}
