package api

import (
	"fmt"
	"io"

	_ "github.com/joho/godotenv/autoload"

	"encoding/json"
	"net/http"

	"github.com/futugyousuzu/goproject/awsgolang/services"
	"github.com/futugyousuzu/goproject/awsgolang/tools"
	verceltool "github.com/futugyousuzu/goproject/awsgolang/vercel"
	model "github.com/futugyousuzu/goproject/awsgolang/viewmodel"

	"github.com/futugyou/extensions"
)

func S3bucket(w http.ResponseWriter, r *http.Request) {
	if extensions.Cors(w, r) {
		return
	}

	if !verceltool.AuthForVercel(w, r) {
		return
	}

	operatetype := r.URL.Query().Get("type")
	switch operatetype {
	case "getBucket":
		getS3bucket(w, r)
	case "getItem":
		getS3bucketItem(w, r)
	case "download":
		downloadFile(w, r)
	case "url":
		getgetS3bucketItemUrl(w, r)
	}
}

func getS3bucket(w http.ResponseWriter, r *http.Request) {
	service := services.NewS3bucketService()
	paging := tools.GetPaging(r.URL.Query())
	bucketName := r.URL.Query().Get("bucketName")
	filter := model.S3BucketFilter{BucketName: bucketName}
	buckets := service.GetS3Buckets(paging, filter)
	body, _ := json.Marshal(buckets)
	w.Write(body)
	w.WriteHeader(200)
}

func getS3bucketItem(w http.ResponseWriter, r *http.Request) {
	service := services.NewS3bucketService()
	bucketName := r.URL.Query().Get("bucketName")
	accountId := r.URL.Query().Get("accountId")
	perfix := r.URL.Query().Get("perfix")
	del := r.URL.Query().Get("del")
	filter := model.S3BucketItemFilter{
		BucketName: bucketName,
		AccountId:  accountId,
		Perfix:     perfix,
		Del:        del,
	}

	items := service.GetS3BucketItems(filter)
	body, _ := json.Marshal(items)
	w.Write(body)
	w.WriteHeader(200)
}

func downloadFile(w http.ResponseWriter, r *http.Request) {
	service := services.NewS3bucketService()
	bucketName := r.URL.Query().Get("bucketName")
	accountId := r.URL.Query().Get("accountId")
	fileName := r.URL.Query().Get("fileName")
	filter := model.S3BucketFileFilter{
		BucketName: bucketName,
		AccountId:  accountId,
		FileName:   fileName,
	}

	file, err := service.GetS3BucketFile(filter)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	_, err = io.Copy(w, file.Body)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Description", "File Transfer")
	w.Header().Set("Content-Disposition", *file.ContentDisposition)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Content-Encoding", *file.ContentEncoding)
	w.Header().Set("Content-Length", fmt.Sprint(file.ContentLength))
	w.WriteHeader(200)
}

func getgetS3bucketItemUrl(w http.ResponseWriter, r *http.Request) {
	service := services.NewS3bucketService()
	bucketName := r.URL.Query().Get("bucketName")
	accountId := r.URL.Query().Get("accountId")
	fileName := r.URL.Query().Get("fileName")
	filter := model.S3BucketFileFilter{
		BucketName: bucketName,
		AccountId:  accountId,
		FileName:   fileName,
	}

	url := service.GetS3FileUrl(filter)
	body, _ := json.Marshal(url)
	w.Write(body)
	w.WriteHeader(200)
}
