package server

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/goproject/chatroom/global"
)

func homeHandleFunc(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFiles(global.RootDir + "/template/home.html")
	if err != nil {
		fmt.Fprint(w, "template error")
		return
	}

	err = tpl.Execute(w, nil)
	if err != nil {
		fmt.Fprint(w, "template execute error")
		return
	}
}
