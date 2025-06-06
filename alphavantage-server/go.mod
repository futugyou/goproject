module github.com/futugyou/alphavantage-server

go 1.23.0

require (
	github.com/futugyou/alphavantage v0.0.0-00010101000000-000000000000
	github.com/futugyou/extensions v0.0.0-00010101000000-000000000000
	github.com/joho/godotenv v1.5.1
)

require github.com/chidiwilliams/flatbson v0.3.0

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/dlclark/regexp2 v1.11.2 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/redis/go-redis/v9 v9.7.3 // indirect
	golang.org/x/mod v0.19.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
)

require (
	github.com/golang/snappy v0.0.4 // indirect
	github.com/jszwec/csvutil v1.10.0 // indirect
	github.com/klauspost/compress v1.16.7 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	go.mongodb.org/mongo-driver v1.17.4
	golang.org/x/crypto v0.35.0 // indirect
	golang.org/x/sync v0.11.0 // indirect
	golang.org/x/text v0.22.0 // indirect
)

replace github.com/futugyou/alphavantage v0.0.0-00010101000000-000000000000 => ../alphavantage

replace github.com/futugyou/extensions v0.0.0-00010101000000-000000000000 => ../extensions
