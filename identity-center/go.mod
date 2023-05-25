module github.com/futugyousuzu/identity

go 1.20

require (
	github.com/futugyousuzu/identity/mongo v0.0.0
	github.com/go-oauth2/oauth2/v4 v4.5.2
	github.com/go-session/session/v3 v3.2.1
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/google/uuid v1.3.0
	github.com/joho/godotenv v1.5.1
	go.mongodb.org/mongo-driver v1.11.6
	golang.org/x/crypto v0.9.0
)

require (
	github.com/bytedance/gopkg v0.0.0-20230512060433-7f5f1dee0b1e // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/klauspost/compress v1.16.5 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/tidwall/btree v1.6.0 // indirect
	github.com/tidwall/buntdb v1.3.0 // indirect
	github.com/tidwall/gjson v1.14.4 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20201027041543-1326539a0a0a // indirect
	golang.org/x/sync v0.2.0 // indirect
	golang.org/x/text v0.9.0 // indirect
)

replace github.com/futugyousuzu/identity/mongo v0.0.0 => github.com/futugyou/goproject/identity-mongo v0.0.0-20230523043116-2b60b63fab72
