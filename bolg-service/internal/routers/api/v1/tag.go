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
func (a Tag) List(c *gin.Context) {
	app.NewResponse(c).ToErrorResponse(errcode.ServiceError)
	log.Print("this is tags list")
}
func (a Tag) Create(c *gin.Context) {}
func (a Tag) Update(c *gin.Context) {}
func (a Tag) Delete(c *gin.Context) {}
