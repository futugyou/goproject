module github.com/futugyousuzu/go-openai-web

go 1.19

require github.com/beego/beego/v2 v2.0.7

require (
	github.com/devfeel/mapper v0.7.10
	github.com/futugyousuzu/go-openai v0.0.0
	github.com/futugyousuzu/openai-tokenizer v0.0.0
	github.com/redis/go-redis/v9 v9.0.3
	github.com/smartystreets/goconvey v1.6.4
)

replace github.com/futugyousuzu/go-openai v0.0.0 => github.com/futugyou/goproject/openai v0.0.0-20230504105143-00ca81260dbf

replace github.com/futugyousuzu/openai-tokenizer v0.0.0 => github.com/futugyou/goproject/openai-tokenizer v0.0.0-20230504105143-00ca81260dbf

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/dlclark/regexp2 v1.9.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/gopherjs/gopherjs v0.0.0-20181017120253-0766667cb4d1 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/jtolds/gls v4.20.0+incompatible // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_golang v1.14.0 // indirect
	github.com/prometheus/client_model v0.3.0 // indirect
	github.com/prometheus/common v0.42.0 // indirect
	github.com/prometheus/procfs v0.9.0 // indirect
	github.com/shiena/ansicolor v0.0.0-20200904210342-c7312218db18 // indirect
	github.com/smartystreets/assertions v0.0.0-20180927180507-b2de0cb4f26d // indirect
	golang.org/x/crypto v0.8.0 // indirect
	golang.org/x/exp v0.0.0-20230321023759-10a507213a29 // indirect
	golang.org/x/net v0.9.0 // indirect
	golang.org/x/sys v0.7.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
