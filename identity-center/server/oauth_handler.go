package server

import (
	"context"
	"fmt"
	"net/http"

	session "github.com/go-session/session/v3"

	"github.com/futugyousuzu/identity/user"
)

func UserAuthorizeHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	store, err := session.Start(r.Context(), w, r)
	if err != nil {
		fmt.Println("something error, sid: " + store.SessionID() + ", info: " + err.Error())
		return
	}

	uid, ok := store.Get("LoggedInUserID")
	if !ok {
		if r.Form == nil {
			r.ParseForm()
		}

		store.Set("ReturnUri", r.Form)
		store.Save()

		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusFound)
		return
	}

	// loginid := uid.(string)
	// userstore := NewUserStore()
	// user, _ := userstore.GetLoginInfo(r.Context(), loginid)
	// userID = user.UserID
	userID = uid.(string)
	store.Delete("LoggedInUserID")
	store.Save()
	return
}

func AuthorizeScopeHandler(w http.ResponseWriter, r *http.Request) (scope string, err error) {
	scope = r.FormValue("scope")
	if len(scope) == 0 {
		scope = "offline, openid, profile"
	} else {
		scope += ", offline, openid, profile"
	}

	return
}

func PasswordAuthorizationHandler(ctx context.Context, clientID, username, password string) (userID string, err error) {
	store := user.NewUserStore()
	user, err := store.Login(ctx, username, password)
	if err == nil {
		userID = user.UserID
	}
	return
}
