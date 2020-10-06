package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/goproject/blog-service/internal/service"
	"github.com/goproject/blog-service/pkg/app"
	"github.com/goproject/blog-service/pkg/convert"
	"github.com/goproject/blog-service/pkg/errcode"
)

type Tag struct{}

func NewTage() Tag {
	return Tag{}
}

func (a Tag) Get(c *gin.Context) {}

// @Summary get tags
// @Produce json
// @Param name query string false "tag name" maxlength(100)
// @Param state query int false "state" Enums(0,1) defaulrt(1)
// @Param page query int false "page"
// @Param page_size query int false "page_size"
// @Success 200 {object} model.TagSwagger "success"
// @Failure 400 {object} errcode.Error "error"
// @Failure 500 {object} errcode.Error "error"
// @Router /api/v1/tags [get]
func (a Tag) List(c *gin.Context) {
	params := service.TagListRequest{}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &params)
	if !valid {
		// global.Logger.Errorf("app bindandvalid err: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	svc := service.New(c.Request.Context())
	pager := app.Pager{Page: app.GetPage(c), PageSize: app.GetPageSize(c)}
	totalRows, err := svc.CountTag(&service.CountTagRequest{Name: params.Name, State: params.State})
	if err != nil {
		// global.Logger.Errorf("service  CountTag err: %v", err)
		response.ToErrorResponse(errcode.ErrorCountTagFail)
		return
	}

	tags, err := svc.GetTagList(&params, &pager)
	if err != nil {
		// global.Logger.Errorf("service  GetTagList err: %v", err)
		response.ToErrorResponse(errcode.ErrorGetListFail)
		return
	}
	response.ToResponseList(tags, totalRows)
}
func (a Tag) Create(c *gin.Context) {
	params := service.CreateTagRequest{}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &params)
	if !valid {
		// global.Logger.Errorf("app bindandvalid err: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	svc := service.New(c.Request.Context())
	err := svc.CreateTag(&params)
	if err != nil {
		// global.Logger.Errorf("service  CreateTag err: %v", err)
		response.ToErrorResponse(errcode.ErrorCreateTagFail)
		return
	}
	response.ToResponse(gin.H{})
}
func (a Tag) Update(c *gin.Context) {

	params := service.UpdateTagRequest{Id: convert.StrTo(c.Param("id")).MustUInt32()}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &params)
	if !valid {
		// global.Logger.Errorf("app bindandvalid err: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	svc := service.New(c.Request.Context())

	err := svc.UpdateTag(&params)
	if err != nil {
		// global.Logger.Errorf("service  UpdateTag err: %v", err)
		response.ToErrorResponse(errcode.ErrorUpdateTagFail)
		return
	}
	response.ToResponse(gin.H{})
}
func (a Tag) Delete(c *gin.Context) {

	params := service.DeleteTagRequest{Id: convert.StrTo(c.Param("id")).MustUInt32()}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &params)
	if !valid {
		// global.Logger.Errorf("app bindandvalid err: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	svc := service.New(c.Request.Context())

	err := svc.DeleteTag(&params)
	if err != nil {
		// global.Logger.Errorf("service  DeleteTag err: %v", err)
		response.ToErrorResponse(errcode.ErrorDeleteTagFail)
		return
	}
	response.ToResponse(gin.H{})
}
