module github.com/futugyousuzu/goproject/awsgolang

go 1.24.0

require (
	github.com/aws/aws-sdk-go-v2 v1.41.6
	github.com/aws/aws-sdk-go-v2/config v1.32.16
	github.com/aws/aws-sdk-go-v2/credentials v1.19.15
	github.com/aws/aws-sdk-go-v2/service/cloudwatch v1.56.2
	github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs v1.71.0
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.57.2
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.299.0
	github.com/aws/aws-sdk-go-v2/service/ecr v1.57.1
	github.com/aws/aws-sdk-go-v2/service/ecs v1.79.0
	github.com/aws/aws-sdk-go-v2/service/efs v1.41.15
	github.com/aws/aws-sdk-go-v2/service/iam v1.53.8
	github.com/aws/aws-sdk-go-v2/service/s3 v1.100.0
	github.com/aws/aws-sdk-go-v2/service/servicediscovery v1.39.27
	github.com/aws/aws-sdk-go-v2/service/ssm v1.68.5
	github.com/chidiwilliams/flatbson v0.3.0
	github.com/futugyou/extensions v0.0.0-00010101000000-000000000000
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
)

require github.com/aws/aws-sdk-go-v2/service/configservice v1.62.2

require github.com/aws/aws-sdk-go-v2/service/route53 v1.62.6

require github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.41.6

require github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2 v1.54.11

require github.com/aws/aws-sdk-go-v2/service/iot v1.72.7

require (
	github.com/aws/aws-sdk-go-v2/service/iotdataplane v1.32.23
	go.mongodb.org/mongo-driver v1.17.9
)

require (
	github.com/aws/aws-sdk-go-v2/service/signin v1.0.10 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/dlclark/regexp2 v1.11.5 // indirect
	github.com/redis/go-redis/v9 v9.17.1 // indirect
	golang.org/x/mod v0.30.0 // indirect
)

require (
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.7.9 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.18.22 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.22 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.22 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.4.23 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.8 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.9.14 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.11.22 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.22 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.19.22 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.30.16 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.35.20 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.42.0 // indirect
	github.com/aws/smithy-go v1.25.0 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.4.0 // indirect
	github.com/goccy/go-json v0.10.5 // indirect
	github.com/golang/snappy v1.0.0 // indirect
	github.com/klauspost/compress v1.18.1 // indirect
	github.com/lestrrat-go/blackmagic v1.0.4 // indirect
	github.com/lestrrat-go/httpcc v1.0.1 // indirect
	github.com/lestrrat-go/httprc v1.0.6 // indirect
	github.com/lestrrat-go/iter v1.0.2 // indirect
	github.com/lestrrat-go/jwx/v2 v2.1.6 // indirect
	github.com/lestrrat-go/option v1.0.1 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/segmentio/asm v1.2.1 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.2.0 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	golang.org/x/crypto v0.45.0 // indirect
	golang.org/x/exp v0.0.0-20251125195548-87e1e737ad39
	golang.org/x/oauth2 v0.33.0 // indirect
	golang.org/x/sync v0.18.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
	golang.org/x/text v0.31.0 // indirect
)

replace github.com/futugyou/extensions v0.0.0-00010101000000-000000000000 => ../extensions
