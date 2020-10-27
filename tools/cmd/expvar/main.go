package main

import (
	"expvar"
	_ "expvar"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

//http://localhost:6060/debug/vars
func main() {
	router := NewRouter()
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		appleCounter.Add(1)
		_, _ = w.Write([]byte("go project"))
	})

	_ = http.ListenAndServe(":6060", router)
}

var (
	appleCounter      *expvar.Int
	GOMAXPROCSMetrics *expvar.Int
	upTimeMetrice     *upTimeVar
)

type upTimeVar struct {
	value time.Time
}

func (v *upTimeVar) Set(date time.Time) {
	v.value = date
}

func (v *upTimeVar) Add(duration time.Duration) {
	v.value = v.value.Add(duration)
}

func (v *upTimeVar) String() string {
	return v.value.Format(time.UnixDate)
}

func init() {
	upTimeMetrice = &upTimeVar{value: time.Now().Local()}
	expvar.Publish("uptime", upTimeMetrice)
	appleCounter = expvar.NewInt("apple")
	GOMAXPROCSMetrics = expvar.NewInt("GOMAXPROCS")
	GOMAXPROCSMetrics.Set(int64(runtime.NumCPU()))
}

func Expvar(c *gin.Context) {
	c.Writer.Header().Set("content-type", "application/json; charset=utf-8")
	first := true
	report := func(key string, value interface{}) {
		if !first {
			fmt.Fprintf(c.Writer, ",\n")
		}
		first = false
		if str, ok := value.(string); ok {
			fmt.Fprintf(c.Writer, "%q: %q", key, str)
		} else {
			fmt.Fprintf(c.Writer, "%q: %v", key, value)
		}
	}

	fmt.Fprintf(c.Writer, "{\n")
	expvar.Do(func(kv expvar.KeyValue) {
		report(kv.Key, kv.Value)
	})
	fmt.Fprintf(c.Writer, "\n}\n")
}

func NewRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.GET("/debug/vars", Expvar)
	return r
}
