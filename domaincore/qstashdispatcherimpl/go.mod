module github.com/futugyou/qstashdispatcherimpl

go 1.24.0

require (
	github.com/futugyou/domaincore v1.0.0
	github.com/futugyou/qstash v1.0.0
)

require github.com/golang-jwt/jwt/v4 v4.5.1 // indirect

replace github.com/futugyou/qstash v1.0.0 => ../../qstash_sdk

replace github.com/futugyou/domaincore v1.0.0 => ../
