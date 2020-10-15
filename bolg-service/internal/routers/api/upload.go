package api

import (
	"github.com/gin-gonic/gin"
	"github.com/goproject/blog-service/global"
	"github.com/goproject/blog-service/internal/service"
	"github.com/goproject/blog-service/pkg/app"
	"github.com/goproject/blog-service/pkg/convert"
	"github.com/goproject/blog-service/pkg/errcode"
	"github.com/goproject/blog-service/pkg/upload"
)

func NewUpload() {

}

func UploadFile(c *gin.Context) {
	response := app.NewResponse(c)
	file, fileHeader, err := c.Request.FormFile("file")
	fileType := convert.StrTo(c.PostForm("type")).MustInt()
	if err != nil {
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error()))
		return
	}
	if fileHeader == nil || fileType <= 0 {
		response.ToErrorResponse(errcode.InvalidParams)
		return
	}
	svc := service.New(c.Request.Context())
	fileInfo, err := svc.UploadFile(upload.FileType(fileType), file, fileHeader)
	if err != nil {
		global.Logger.Errorf(c, "svc.uploadFile err : %v", err)
		response.ToErrorResponse(errcode.ErrorUploadFileFail.WithDetails(err.Error()))
		return
	}

	response.ToResponse(gin.H{"file_access_url": fileInfo.AccessUrl})
}
