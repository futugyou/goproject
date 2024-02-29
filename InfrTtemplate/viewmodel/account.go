package viewmodel

import "time"

type UserAccountRequest struct {
	Id              string    `json:"id"`
	Alias           string    `json:"alias"`
	AccessKeyId     string    `json:"accessKeyId"`
	SecretAccessKey string    `json:"secretAccessKey"`
	Region          string    `json:"region"`
	CreatedAt       time.Time `json:"createdAt"`
}

type UserAccountResponse struct {
	Alias           string    `json:"alias"`
	AccessKeyId     string    `json:"accessKeyId"`
	SecretAccessKey string    `json:"secretAccessKey"`
	Region          string    `json:"region"`
	CreatedAt       time.Time `json:"createdAt"`
}
