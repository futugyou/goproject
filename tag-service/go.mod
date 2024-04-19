module github.com/goproject/tag-service

go 1.15

require (
	github.com/HdrHistogram/hdrhistogram-go v1.1.2 // indirect
	github.com/elazarl/go-bindata-assetfs v1.0.1
	github.com/golang/protobuf v1.5.3
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.2
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/opentracing/opentracing-go v1.2.0
	github.com/uber/jaeger-client-go v2.25.0+incompatible
	github.com/uber/jaeger-lib v2.4.0+incompatible // indirect
	go.uber.org/atomic v1.7.0 // indirect
	golang.org/x/net v0.23.0
	google.golang.org/genproto v0.0.0-20230410155749-daa745c078e1
	google.golang.org/grpc v1.56.3
	google.golang.org/protobuf v1.33.0
)

//replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
