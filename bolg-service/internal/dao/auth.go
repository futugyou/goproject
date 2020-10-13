package dao

import "github.com/goproject/blog-service/internal/model"

func (a *Dao) GetAuth(appkey, appsecret string) (model.Auth, error) {
	auth := model.Auth{AppKey: appkey, AppSecret: appsecret}
	return auth.Get(a.engine)
}
