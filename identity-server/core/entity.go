package core

//go:generate gomockhandler -config=../gomockhandler.json  -destination ../mocks/mock_entity_test.go -package=core_test github.com/futugyousuzu/identity-server/core IEntity

type IEntity interface {
	GetType() string
}
