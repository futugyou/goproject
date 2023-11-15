package entity

import "time"

type S3bucketEntity struct {
	Id           string    `bson:"_id,omitempty"`
	Name         string    `bson:"name"`
	Region       string    `bson:"region"`
	IsPublic     bool      `bson:"isPublic"`
	Policy       string    `bson:"policy"`
	Permissions  []string  `bson:"permissions"`
	CreationDate time.Time `bson:"creationDate"`
}

func (S3bucketEntity) GetType() string {
	return "s3bucket"
}

type S3bucketItemEntity struct {
	Id           string    `bson:"_id,omitempty"`
	BucketName   string    `bson:"bucketName"`
	Key          string    `bson:"key"`
	Size         int64     `bson:"size"`
	CreationDate time.Time `bson:"creationDate"`
}

func (S3bucketItemEntity) GetType() string {
	return "s3bucket-items"
}
