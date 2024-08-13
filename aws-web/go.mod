module github.com/futugyousuzu/goproject/awsgolang

go 1.20

require (
	github.com/aws/aws-sdk-go-v2 v1.30.3
	github.com/aws/aws-sdk-go-v2/config v1.27.27
	github.com/aws/aws-sdk-go-v2/credentials v1.17.27
	github.com/aws/aws-sdk-go-v2/service/cloudwatch v1.40.3
	github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs v1.37.3
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.34.4
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.175.0
	github.com/aws/aws-sdk-go-v2/service/ecr v1.32.0
	github.com/aws/aws-sdk-go-v2/service/ecs v1.44.3
	github.com/aws/aws-sdk-go-v2/service/efs v1.31.3
	github.com/aws/aws-sdk-go-v2/service/iam v1.34.3
	github.com/aws/aws-sdk-go-v2/service/s3 v1.58.3
	github.com/aws/aws-sdk-go-v2/service/servicediscovery v1.31.3
	github.com/aws/aws-sdk-go-v2/service/ssm v1.52.4
	github.com/chidiwilliams/flatbson v0.3.0
	github.com/futugyou/extensions v0.0.0-00010101000000-000000000000
	github.com/futugyousuzu/identity/client v0.0.0
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
)

require github.com/aws/aws-sdk-go-v2/service/configservice v1.48.3

require github.com/aws/aws-sdk-go-v2/service/route53 v1.42.3

require github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.32.4

require github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2 v1.34.0

require github.com/aws/aws-sdk-go-v2/service/iot v1.55.3

require (
	github.com/aws/aws-sdk-go-v2/service/iotdataplane v1.24.3
	go.mongodb.org/mongo-driver v1.16.1
)

require (
	github.com/dlclark/regexp2 v1.11.2 // indirect
	golang.org/x/mod v0.19.0 // indirect
)

require (
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.6.3 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.11 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.15 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.15 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.0 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.3.15 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.11.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.3.17 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.9.16 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.11.17 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.17.15 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.22.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.26.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.30.3 // indirect
	github.com/aws/smithy-go v1.20.3 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.2.0 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
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
	github.com/youmark/pkcs8 v0.0.0-20201027041543-1326539a0a0a // indirect
	golang.org/x/crypto v0.25.0 // indirect
	golang.org/x/exp v0.0.0-20231219180239-dc181d75b848
	golang.org/x/net v0.24.0 // indirect
	golang.org/x/oauth2 v0.10.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/sys v0.22.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
)

replace github.com/futugyousuzu/identity/client v0.0.0 => github.com/futugyou/goproject/identity-client v0.0.0-20230713085205-834db99b0998

replace github.com/futugyou/extensions v0.0.0-00010101000000-000000000000 => ../extensions
