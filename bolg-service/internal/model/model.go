package model

import (
	"fmt"

	"github.com/goproject/blog-service/global"

	"github.com/goproject/blog-service/pkg/setting"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Model struct {
	ID          uint32 `gorm:"primary_key" json:"id"`
	CreatedBy   string `json:"created_by"`
	CreatedOn   uint32 `json:"created_on"`
	ModififedBy string `json:"modifed_by"`
	ModififedOn uint32 `json:"modifed_on"`
	Deleted     uint32 `json:"deleted_on"`
	IsDel       uint32 `json:"is_del"`
}

func NewDBEngine(dbsetting *setting.DatabaseSettingS) (*gorm.DB, error) {
	db, err := gorm.Open(dbsetting.DBType,
		fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=Local",
			dbsetting.Username,
			dbsetting.Password,
			dbsetting.Host,
			dbsetting.DBName,
			dbsetting.Charset,
			dbsetting.ParseTime,
		))
	if err != nil {
		return nil, err
	}
	if global.ServerSetting.RunModel == "debug" {
		db.LogMode(true)
	}
	db.SingularTable(true)
	db.DB().SetMaxIdleConns(dbsetting.MaxIdleConns)
	db.DB().SetMaxOpenConns(dbsetting.MaxOpenConns)
	return db, nil
}
