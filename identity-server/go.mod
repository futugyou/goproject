module github.com/futugyousuzu/identity-server

go 1.24.0

require (
	github.com/chidiwilliams/flatbson v0.3.0
	github.com/futugyou/domaincore v1.0.0
	github.com/futugyou/extensions v1.0.0
	github.com/go-oauth2/oauth2/v4 v4.5.2
	github.com/joho/godotenv v1.5.1
	github.com/lestrrat-go/jwx/v2 v2.1.6
	github.com/sendgrid/sendgrid-go v3.16.1+incompatible
	go.mongodb.org/mongo-driver v1.17.6
	go.uber.org/mock v0.2.0
	golang.org/x/crypto v0.45.0
	golang.org/x/oauth2 v0.33.0
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.4.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/dlclark/regexp2 v1.11.5 // indirect
	github.com/goccy/go-json v0.10.5 // indirect
	github.com/golang/snappy v1.0.0 // indirect
	github.com/klauspost/compress v1.18.1 // indirect
	github.com/lestrrat-go/blackmagic v1.0.4 // indirect
	github.com/lestrrat-go/httpcc v1.0.1 // indirect
	github.com/lestrrat-go/httprc v1.0.6 // indirect
	github.com/lestrrat-go/iter v1.0.2 // indirect
	github.com/lestrrat-go/option v1.0.1 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/redis/go-redis/v9 v9.17.1 // indirect
	github.com/segmentio/asm v1.2.1 // indirect
	github.com/sendgrid/rest v2.6.9+incompatible // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.2.0 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	golang.org/x/mod v0.30.0 // indirect
	golang.org/x/sync v0.18.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
	golang.org/x/text v0.31.0 // indirect
)

require github.com/google/uuid v1.6.0

replace (
	github.com/futugyou/domaincore v1.0.0 => ../domaincore
	github.com/futugyou/domaincore/mongoimpl v1.0.0 => ../domaincore/mongoimpl
	github.com/futugyou/extensions v1.0.0 => ../extensions
)
