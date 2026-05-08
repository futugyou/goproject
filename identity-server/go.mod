module github.com/futugyousuzu/identity-server

go 1.26.0

require (
	github.com/chidiwilliams/flatbson v0.3.0
	github.com/futugyou/domaincore v1.0.0
	github.com/futugyou/domaincore/mongoimpl v1.0.0
	github.com/futugyou/extensions v1.0.0
	github.com/go-oauth2/oauth2/v4 v4.5.2
	github.com/joho/godotenv v1.5.1
	github.com/sendgrid/sendgrid-go v3.16.1+incompatible
	go.mongodb.org/mongo-driver v1.17.6
	go.uber.org/mock v0.2.0
	golang.org/x/crypto v0.49.0
	golang.org/x/oauth2 v0.33.0
)

require (
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/ProjectZKM/Ziren/crates/go-runtime/zkvm_runtime v0.0.0-20251001021608-1fe7b43fc4d6 // indirect
	github.com/StackExchange/wmi v1.2.1 // indirect
	github.com/bits-and-blooms/bitset v1.20.0 // indirect
	github.com/cayleygraph/quad v1.3.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/consensys/gnark-crypto v0.18.1 // indirect
	github.com/crate-crypto/go-eth-kzg v1.5.0 // indirect
	github.com/deckarep/golang-set/v2 v2.6.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/dlclark/regexp2 v1.11.5 // indirect
	github.com/ethereum/c-kzg-4844/v2 v2.1.6 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/goccy/go-json v0.10.5 // indirect
	github.com/golang/snappy v1.0.0 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/holiman/uint256 v1.3.2 // indirect
	github.com/klauspost/compress v1.18.1 // indirect
	github.com/lestrrat-go/blackmagic v1.0.4 // indirect
	github.com/lestrrat-go/dsig v1.3.0 // indirect
	github.com/lestrrat-go/httpcc v1.0.1 // indirect
	github.com/lestrrat-go/httprc v1.0.6 // indirect
	github.com/lestrrat-go/httprc/v3 v3.0.5 // indirect
	github.com/lestrrat-go/iter v1.0.2 // indirect
	github.com/lestrrat-go/option v1.0.1 // indirect
	github.com/lestrrat-go/option/v2 v2.0.0 // indirect
	github.com/lestrrat-go/option/v3 v3.0.0-alpha1 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/mr-tron/base58 v1.3.0 // indirect
	github.com/multiformats/go-base32 v0.1.0 // indirect
	github.com/multiformats/go-base36 v0.2.0 // indirect
	github.com/pquerna/cachecontrol v0.2.0 // indirect
	github.com/redis/go-redis/v9 v9.17.1 // indirect
	github.com/segmentio/asm v1.2.1 // indirect
	github.com/sendgrid/rest v2.6.9+incompatible // indirect
	github.com/shirou/gopsutil v3.21.4-0.20210419000835-c7a38de76ee5+incompatible // indirect
	github.com/supranational/blst v0.3.16 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	github.com/tklauser/go-sysconf v0.3.12 // indirect
	github.com/tklauser/numcpus v0.6.1 // indirect
	github.com/valyala/fastjson v1.6.10 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.2.0 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20180127040702-4e3ac2762d5f // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel v1.40.0 // indirect
	go.opentelemetry.io/otel/metric v1.40.0 // indirect
	go.opentelemetry.io/otel/trace v1.40.0 // indirect
	golang.org/x/mod v0.35.0 // indirect
	golang.org/x/sync v0.20.0 // indirect
	golang.org/x/sys v0.42.0 // indirect
	golang.org/x/text v0.35.0 // indirect
)

require (
	github.com/btcsuite/btcd/btcutil v1.1.6
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.4.0
	github.com/ethereum/go-ethereum v1.17.2
	github.com/google/uuid v1.6.0
	github.com/jwx-go/jwkfetch/v4 v4.0.1
	github.com/lestrrat-go/jwx/v2 v2.1.6
	github.com/lestrrat-go/jwx/v4 v4.0.0
	github.com/multiformats/go-multibase v0.3.0
	github.com/piprate/json-gold v0.8.0
	github.com/tidwall/gjson v1.18.0
	github.com/xeipuuv/gojsonschema v1.2.0
)

replace (
	github.com/futugyou/domaincore v1.0.0 => ../domaincore
	github.com/futugyou/domaincore/mongoimpl v1.0.0 => ../domaincore/mongoimpl
	github.com/futugyou/extensions v1.0.0 => ../extensions
)
