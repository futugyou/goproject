module github.com/futugyousuzu/goproject/awsgolang

go 1.23.0

require (
	github.com/aws/aws-sdk-go-v2 v1.38.3
	github.com/aws/aws-sdk-go-v2/config v1.31.6
	github.com/aws/aws-sdk-go-v2/credentials v1.18.10
	github.com/aws/aws-sdk-go-v2/service/cloudwatch v1.49.2
	github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs v1.57.2
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.50.1
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.250.0
	github.com/aws/aws-sdk-go-v2/service/ecr v1.50.1
	github.com/aws/aws-sdk-go-v2/service/ecs v1.63.4
	github.com/aws/aws-sdk-go-v2/service/efs v1.40.3
	github.com/aws/aws-sdk-go-v2/service/iam v1.47.3
	github.com/aws/aws-sdk-go-v2/service/s3 v1.87.3
	github.com/aws/aws-sdk-go-v2/service/servicediscovery v1.39.5
	github.com/aws/aws-sdk-go-v2/service/ssm v1.64.2
	github.com/chidiwilliams/flatbson v0.3.0
	github.com/futugyou/extensions v0.0.0-00010101000000-000000000000
	github.com/futugyousuzu/identity/client v0.0.0
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
)

require github.com/aws/aws-sdk-go-v2/service/configservice v1.57.2

require github.com/aws/aws-sdk-go-v2/service/route53 v1.58.0

require github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.39.2

require github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2 v1.50.2

require github.com/aws/aws-sdk-go-v2/service/iot v1.69.1

require (
	github.com/aws/aws-sdk-go-v2/service/iotdataplane v1.32.2
	go.mongodb.org/mongo-driver v1.17.4
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/dlclark/regexp2 v1.11.2 // indirect
	github.com/redis/go-redis/v9 v9.7.3 // indirect
	golang.org/x/mod v0.19.0 // indirect
)

require (
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.7.1 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.18.6 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.6 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.6 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.3 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.4.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.8.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.11.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.19.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.29.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.34.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.38.2 // indirect
	github.com/aws/smithy-go v1.23.0 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.2.0 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/klauspost/compress v1.17.4 // indirect
	github.com/lestrrat-go/blackmagic v1.0.2 // indirect
	github.com/lestrrat-go/httpcc v1.0.1 // indirect
	github.com/lestrrat-go/httprc v1.0.5 // indirect
	github.com/lestrrat-go/iter v1.0.2 // indirect
	github.com/lestrrat-go/jwx/v2 v2.0.21 // indirect
	github.com/lestrrat-go/option v1.0.1 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/segmentio/asm v1.2.0 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	golang.org/x/crypto v0.36.0 // indirect
	golang.org/x/exp v0.0.0-20231219180239-dc181d75b848
	golang.org/x/oauth2 v0.27.0 // indirect
	golang.org/x/sync v0.12.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
)

replace github.com/futugyousuzu/identity/client v0.0.0 => github.com/futugyou/goproject/identity-client v0.0.0-20230713085205-834db99b0998

replace github.com/futugyou/extensions v0.0.0-00010101000000-000000000000 => ../extensions
