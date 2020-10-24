package server

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/goproject/chatroom/logic"
)

func RegisterHandle() {
	inferRootDir()

	go logic.Broadcaster.Start()

	http.HandleFunc("/", homeHandleFunc)
	http.HandleFunc("/ws", WebSocketHandleFunc)
}

var rootDir string

func inferRootDir() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	var infer func(d string) string
	infer = func(d string) string {
		if exists(d + "/template") {
			return d

		}
		return infer(filepath.Dir(d))
	}
	rootDir = infer(cwd)
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return nil == err || os.IsExist(err)
}
