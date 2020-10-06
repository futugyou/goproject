package service

import (
	"github.com/goproject/blog-service/internal/model"
	"github.com/goproject/blog-service/pkg/app"
)

type CountTagRequest struct {
	Name  string `form:"name" binding:"max=100"`
	State uint8  `form:"state,default=1" binding:"oneof=0 1"`
}

type TagListRequest struct {
	Name  string `form:"name" binding:"max=100"`
	State uint8  `form:"state,default=1" binding:"oneof=0 1"`
}

type CreateTagRequest struct {
	Name     string `form:"name" binding:"required,min=3,max=100"`
	CreateBy string `form:"created_by" binding:"required,min=3,max=100"`
	State    uint8  `form:"state,default=1" binding:"oneof=0 1"`
}

type UpdateTagRequest struct {
	Id         uint32 `form:"id" binding:"required,gte=1"`
	Name       string `form:"name" binding:"required,min=3,max=100"`
	ModifiedBy string `form:"modified_by" binding:"required,min=3,max=100"`
	State      uint8  `form:"state,default=1" binding:"oneof=0 1"`
}

type DeleteTagRequest struct {
	Id uint32 `form:"id" binding:"required,gte=1"`
}

func (s *Service) CountTag(param *CountTagRequest) (int, error) {
	return s.dao.CountTag(param.Name, param.State)
}

func (s *Service) GetTagList(param *TagListRequest, pager *app.Pager) ([]*model.Tag, error) {
	return s.dao.GetTagList(param.Name, param.State, pager.Page, pager.PageSize)
}
func (s *Service) CreateTag(param *CreateTagRequest) error {
	return s.dao.CreateTag(param.Name, param.State, param.CreateBy)
}

func (s *Service) UpdateTag(param *UpdateTagRequest) error {
	return s.dao.UpdateTag(param.Id, param.Name, param.State, param.ModifiedBy)
}

func (s *Service) DeleteTag(param *DeleteTagRequest) error {
	return s.dao.DeleteTag(param.Id)
}
