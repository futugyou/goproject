package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/goproject/blog-service/pkg/app"
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
	params := struct {
		Name  string `form:"name" binding:"max=100"`
		State uint8  `form:"state,default 1" binding:"oneof=0 1"`
	}{}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &params)
	if valid {
		// global.Logger.Errorf("app bindandvalid err: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	response.ToResponse(gin.H{})
}
func (a Tag) Create(c *gin.Context) {}
func (a Tag) Update(c *gin.Context) {}
func (a Tag) Delete(c *gin.Context) {}
