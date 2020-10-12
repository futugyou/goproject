package dao

import (
	"github.com/goproject/blog-service/internal/model"
	"github.com/goproject/blog-service/pkg/app"
)

func (d *Dao) CountTag(name string, state uint8) (int, error) {
	tag := model.Tag{Name: name, State: state}
	return tag.Count(d.engine)
}

func (d *Dao) GetTag(id uint32, state uint8) (model.Tag, error) {
	tag := model.Tag{Model: &model.Model{ID: id}, State: state}
	return tag.Get(d.engine)
}

func (d *Dao) GetTagList(name string, state uint8, page, pagesize int) ([]*model.Tag, error) {
	tag := model.Tag{Name: name, State: state}
	pageoffset := app.GetPageOffset(page, pagesize)
	return tag.List(d.engine, pageoffset, pagesize)
}

func (d *Dao) CreateTag(name string, state uint8, createby string) error {
	tag := model.Tag{
		Name:  name,
		State: state,
		Model: &model.Model{CreatedBy: createby},
	}
	return tag.Create(d.engine)
}

func (d *Dao) UpdateTag(id uint32, name string, state uint8, modifiedBy string) error {
	tag := model.Tag{
		Name:  name,
		State: state,
		Model: &model.Model{ModifiedBy: modifiedBy, ID: id},
	}
	data := make(map[string]interface{})
	data["modifiedBy"] = modifiedBy
	data["state"] = state
	if name != "" {
		data["name"] = name
	}

	return tag.Update(d.engine, data)
}

func (d *Dao) DeleteTag(id uint32) error {
	tag := model.Tag{Model: &model.Model{ID: id}}
	return tag.Delete(d.engine)
}
