package service

import (
	"errors"
	"mime/multipart"
	"os"

	"github.com/goproject/blog-service/global"

	"github.com/goproject/blog-service/pkg/upload"
)

type FileInfo struct {
	Name      string
	AccessUrl string
}

func (svc *Service) UploadFile(fileType upload.FileType, file multipart.File, fileHeader *multipart.FileHeader) (*FileInfo, error) {
	fileName := upload.GetFileName(fileHeader.Filename)
	uploadSaveParh := upload.GetSavePath()
	dst := uploadSaveParh + "/" + fileName
	if !upload.CheckContainExt(fileType, fileName) {
		return nil, errors.New("file suffix is not supported.")
	}
	if upload.CheckSavePath(uploadSaveParh) {
		if err := upload.CreateSavePath(uploadSaveParh, os.ModePerm); err != nil {
			return nil, errors.New("failed to create save directory")
		}
	}
	if upload.CheckMaxSize(fileType, file) {
		return nil, errors.New("exceeded maximu file limit.")
	}
	if upload.CheckPermission(uploadSaveParh) {
		return nil, errors.New("insufficient file permissions.")
	}
	if err := upload.SaveFile(fileHeader, dst); err != nil {
		return nil, err
	}
	accessUrl := global.AppSetting.UploadServerUrl + "/" + fileName
	return &FileInfo{Name: fileName, AccessUrl: accessUrl}, nil
}
