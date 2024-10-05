package entity

type AccountEntity struct {
	Id              string `bson:"_id,omitempty"`
	Alias           string `bson:"alias"`
	AccessKeyId     string `bson:"accessKeyId"`
	SecretAccessKey string `bson:"secretAccessKey"`
	Region          string `bson:"region"`
	Valid           bool   `bson:"valid"`
	CreatedAt       int64  `bson:"createdAt"`
}

func (AccountEntity) GetType() string {
	return "accounts"
}
