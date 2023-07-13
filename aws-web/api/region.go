package api

import (
	"net/http"

	verceltool "github.com/futugyousuzu/goproject/awsgolang/vercel"
)

func Get(w http.ResponseWriter, r *http.Request) {
	if verceltool.CrosForVercel(w, r) {
		return
	}

	if !verceltool.AuthForVercel(w, r) {
		return
	}

	w.Write([]byte("ok"))
	w.WriteHeader(200)
}
