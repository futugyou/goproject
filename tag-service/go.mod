module github.com/goproject/tag-service

go 1.15

require (
	github.com/HdrHistogram/hdrhistogram-go v1.1.2 // indirect
	github.com/elazarl/go-bindata-assetfs v1.0.1
	github.com/golang/protobuf v1.4.3
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.2
	github.com/grpc-ecosystem/grpc-gateway v1.15.2
	github.com/opentracing/opentracing-go v1.2.0
	github.com/pkg/errors v0.9.1 // indirect
	github.com/stretchr/objx v0.2.0 // indirect
	github.com/uber/jaeger-client-go v2.25.0+incompatible
	github.com/uber/jaeger-lib v2.4.0+incompatible // indirect
	go.uber.org/atomic v1.7.0 // indirect
	golang.org/x/net v0.7.0
	google.golang.org/genproto v0.0.0-20201019141844-1ed22bb0c154
	google.golang.org/grpc v1.33.1
	google.golang.org/protobuf v1.25.0
)

//replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
