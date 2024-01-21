package main

import (
	_ "github.com/joho/godotenv/autoload"

	"os"
)

type PusherConfig struct {
	APP_ID      string
	APP_KEY     string
	APP_SECRET  string
	APP_CLUSTER string
	INSTANCE_ID string
	SECRET_KEY  string
}

func NewPusherConfig() *PusherConfig {
	return &PusherConfig{
		APP_ID:      os.Getenv("APP_ID"),
		APP_KEY:     os.Getenv("APP_KEY"),
		APP_SECRET:  os.Getenv("APP_SECRET"),
		APP_CLUSTER: os.Getenv("APP_CLUSTER"),
		INSTANCE_ID: os.Getenv("INSTANCE_ID"),
		SECRET_KEY:  os.Getenv("SECRET_KEY"),
	}
}
