package entity

type AccountEntity struct {
	Id              string `bson:"_id,omitempty"`
	Alias           string `bson:"alias"`
	AccessKeyId     string `bson:"accessKeyId"`
	SecretAccessKey string `bson:"secretAccessKey"`
	Region          string `bson:"region"`
	CreatedAt       int    `bson:"createdAt"`
}

func (AccountEntity) GetType() string {
	return "accounts"
}
