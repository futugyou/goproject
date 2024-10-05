package viewmodel

import "time"

type UserAccount struct {
	Id              string    `json:"id"`
	Alias           string    `json:"alias"`
	AccessKeyId     string    `json:"accessKeyId"`
	SecretAccessKey string    `json:"secretAccessKey"`
	Region          string    `json:"region"`
	Valid           bool      `bson:"valid"`
	CreatedAt       time.Time `json:"createdAt"`
}
