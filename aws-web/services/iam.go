package services

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/iam"

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
	account := accountService.GetAccountByID(ctx, filter.AccountId)
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
			data.Keys = append(data.Keys, model.IAMDataKey{
				Key:        *key.AccessKeyId,
				Status:     string(key.Status),
				CreateDate: *key.CreateDate,
			})
		}

		result = append(result, data)
	}

	return result, nil
}
