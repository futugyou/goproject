package v1

import (
	"log"

	"github.com/gin-gonic/gin"
)

type Tag struct{}

func NewTage() Tag {
	return Tag{}
}

func (a Tag) Get(c *gin.Context)    {}
func (a Tag) List(c *gin.Context)   { log.Print("this is tags list") }
func (a Tag) Create(c *gin.Context) {}
func (a Tag) Update(c *gin.Context) {}
func (a Tag) Delete(c *gin.Context) {}
