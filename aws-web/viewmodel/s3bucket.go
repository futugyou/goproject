package viewmodel

import "time"

type S3BucketViewModel struct {
	Id           string    `json:"id,omitempty"`
	AccountId    string    `json:"accountId"`
	AccountName  string    `json:"accountName"`
	Name         string    `json:"name"`
	Region       string    `json:"region"`
	IsPublic     bool      `json:"isPublic"`
	Policy       string    `json:"policy"`
	Permissions  []string  `json:"permissions"`
	CreationDate time.Time `json:"creationDate"`
}

type S3BucketItemViewModel struct {
	Id           string    `json:"id,omitempty"`
	BucketName   string    `json:"bucketName"`
	Key          string    `json:"key"`
	Size         int64     `json:"size"`
	CreationDate time.Time `json:"creationDate"`
	IsDirectory  bool      `json:"isDirectory"`
}

type S3BucketFilter struct {
	BucketName string `json:"bucketName,omitempty"`
}

type S3BucketItemFilter struct {
	BucketName string `json:"bucketName,omitempty"`
	AccountId  string `json:"accountId"`
	Perfix     string `json:"perfix,omitempty"`
}

type S3BucketFileFilter struct {
	BucketName string `json:"bucketName,omitempty"`
	AccountId  string `json:"accountId"`
	FileName   string `json:"fileName,omitempty"`
}

type S3BucketUrlViewModel struct {
	Url string `json:"url"`
}
