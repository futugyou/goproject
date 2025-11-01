module github.com/futugyou/platformservice

go 1.24.0

require (
	github.com/futugyou/domaincore v1.0.0
	github.com/google/uuid v1.6.0
	go.mongodb.org/mongo-driver v1.17.6
)

require github.com/google/go-cmp v0.7.0 // indirect

replace (
	github.com/futugyou/domaincore v1.0.0 => ../../../domaincore
	github.com/futugyou/domaincore/mongoimpl v1.0.0 => ../../../domaincore/mongoimpl
	github.com/futugyou/extensions v1.0.0 => ../../../extensions
)
