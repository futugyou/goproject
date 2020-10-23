package server

import (
	"fmt"
	"html/template"
	"net/http"
)

func homeHandleFunc(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFiles(rootDir + "/template/home.html")
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
