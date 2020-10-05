package v1

import (
	"log"

	"github.com/goproject/blog-service/pkg/errcode"

	"github.com/gin-gonic/gin"
	"github.com/goproject/blog-service/pkg/app"
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
	app.NewResponse(c).ToErrorResponse(errcode.ServiceError)
	log.Print("this is tags list")
}
func (a Tag) Create(c *gin.Context) {}
func (a Tag) Update(c *gin.Context) {}
func (a Tag) Delete(c *gin.Context) {}
