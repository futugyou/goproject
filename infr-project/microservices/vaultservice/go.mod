module github.com/futugyou/vaultservice

go 1.24.0

require (
	github.com/futugyou/domaincore v1.0.0
	github.com/futugyou/domaincore/mongoimpl v1.0.0
	github.com/futugyou/extensions v1.0.0
	github.com/google/uuid v1.6.0
	go.mongodb.org/mongo-driver v1.17.6
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/chidiwilliams/flatbson v0.3.0 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/dlclark/regexp2 v1.11.2 // indirect
	github.com/goccy/go-json v0.10.5 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/golang/snappy v1.0.0 // indirect
	github.com/klauspost/compress v1.18.1 // indirect
	github.com/lestrrat-go/blackmagic v1.0.1 // indirect
	github.com/lestrrat-go/httpcc v1.0.1 // indirect
	github.com/lestrrat-go/httprc v1.0.4 // indirect
	github.com/lestrrat-go/iter v1.0.2 // indirect
	github.com/lestrrat-go/jwx/v2 v2.0.11 // indirect
	github.com/lestrrat-go/option v1.0.1 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/redis/go-redis/v9 v9.7.0 // indirect
	github.com/segmentio/asm v1.2.0 // indirect
	github.com/stretchr/testify v1.11.1 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	golang.org/x/crypto v0.41.0 // indirect
	golang.org/x/mod v0.27.0 // indirect
	golang.org/x/net v0.43.0 // indirect
	golang.org/x/oauth2 v0.10.0 // indirect
	golang.org/x/sync v0.16.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
	golang.org/x/text v0.28.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.36.10 // indirect
)

replace (
	github.com/futugyou/domaincore v1.0.0 => ../../../domaincore
	github.com/futugyou/domaincore/mongoimpl v1.0.0 => ../../../domaincore/mongoimpl
	github.com/futugyou/domaincore/qstashdispatcherimpl v1.0.0 => ../../../domaincore/qstashdispatcherimpl
	github.com/futugyou/extensions v1.0.0 => ../../../extensions
	github.com/futugyou/qstash v1.0.0 => ../../../qstash_sdk
)
