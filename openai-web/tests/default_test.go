package test

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"testing"

	_ "github.com/futugyousuzu/go-openai-web/routers"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	_, file, _, _ := runtime.Caller(0)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".."+string(filepath.Separator))))
	beego.TestBeegoInit(apppath)
}

// TestGet is a sample to run an endpoint test
func TestGet(t *testing.T) {
	r, _ := http.NewRequest("GET", "/api/v1/examples", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	logs.Info("testing", "TestGet", "Code[%d]\n%s", w.Code, w.Body.String())

	// This is just for testing purposes
	// The `app.conf` file doesn't contain a real database connection, so the HTTP response will inevitably be 500.
	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Status Code Should Be 500", func() {
			So(w.Code, ShouldEqual, 500)
		})
		Convey("The Result Should Be Empty", func() {
			So(w.Body.Len(), ShouldEqual, 0)
		})
	})
}
