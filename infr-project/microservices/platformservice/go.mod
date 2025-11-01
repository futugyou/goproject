module github.com/futugyou/platformservice

go 1.24.0

require (
	github.com/futugyou/domaincore v1.0.0
	github.com/futugyou/domaincore/mongoimpl v1.0.0
	github.com/google/uuid v1.6.0
	go.mongodb.org/mongo-driver v1.17.6
)

require (
	github.com/chidiwilliams/flatbson v0.3.0 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/klauspost/compress v1.16.7 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	golang.org/x/crypto v0.26.0 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/text v0.17.0 // indirect
)

replace (
	github.com/futugyou/domaincore v1.0.0 => ../../../domaincore
	github.com/futugyou/domaincore/mongoimpl v1.0.0 => ../../../domaincore/mongoimpl
	github.com/futugyou/extensions v1.0.0 => ../../../extensions
)
