package model

import (
	"fmt"
	"time"

	"github.com/goproject/blog-service/global"

	otgorm "github.com/eddycjy/opentracing-gorm"
	"github.com/goproject/blog-service/pkg/setting"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Model struct {
	ID         uint32 `gorm:"primary_key" json:"id"`
	CreatedBy  string `json:"created_by"`
	CreatedOn  uint32 `json:"created_on"`
	ModifiedBy string `json:"modified_by"`
	ModifiedOn uint32 `json:"modified_on"`
	DeletedOn  uint32 `json:"deleted_on"`
	IsDel      uint32 `json:"is_del"`
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
	if db.HasTable(&Tag{}) {
		db.AutoMigrate(&Tag{})
	} else {
		db.CreateTable(&Tag{})
	}
	if db.HasTable(&Article{}) {
		db.AutoMigrate(&Article{})
	} else {
		db.CreateTable(&Article{})
	}
	if db.HasTable(&ArticleTag{}) {
		db.AutoMigrate(&ArticleTag{})
	} else {
		db.CreateTable(&ArticleTag{})
	}
	if db.HasTable(&Auth{}) {
		db.AutoMigrate(&Auth{})
	} else {
		db.CreateTable(&Auth{})
	}

	db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
	db.Callback().Delete().Replace("gorm:delete", deleteCallback)

	db.DB().SetMaxIdleConns(dbsetting.MaxIdleConns)
	db.DB().SetMaxOpenConns(dbsetting.MaxOpenConns)
	otgorm.AddGormCallbacks(db)
	return db, nil
}

func addExtraSpaceIfExist(str string) string {
	if str != "" {
		return " " + str
	}
	return ""
}

func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
	if _, ok := scope.Get("gorm:update_column"); !ok {
		_ = scope.SetColumn("ModifiedOn", time.Now().Unix())
	}
}

func updateTimeStampForCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		nowTime := time.Now().Unix()
		if createTimeField, ok := scope.FieldByName("CreatedOn"); ok {
			if createTimeField.IsBlank {
				_ = createTimeField.Set(nowTime)
			}
		}
		if modifyTimeField, ok := scope.FieldByName("ModifiedOn"); ok {
			if modifyTimeField.IsBlank {
				_ = modifyTimeField.Set(nowTime)
			}
		}
	}
}

func deleteCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		var extraoption string
		if str, ok := scope.Get("gorm:delete_option"); ok {
			extraoption = fmt.Sprint(str)
		}
		deleteOnField, hasDeletedonField := scope.FieldByName("DeletedOn")
		isDelFiled, hasIsDelField := scope.FieldByName("IsDel")
		if !scope.Search.Unscoped && hasDeletedonField && hasIsDelField {
			now := time.Now().Unix()
			scope.Raw(fmt.Sprintf(
				"update %v set %v=%v ,%v=%v%v%v",
				scope.QuotedTableName(),
				scope.Quote(deleteOnField.DBName),
				scope.AddToVars(now),
				scope.Quote(isDelFiled.DBName),
				scope.AddToVars(1),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraoption),
			)).Exec()
		} else {
			scope.Raw(fmt.Sprintf(
				"delete from %v%v%v",
				scope.QuotedTableName(),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraoption),
			)).Exec()
		}
	}
}
