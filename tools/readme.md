# 第六章： Go中的大杀器

## 简介

介绍了PProf，trace，godebug，gops，metrrics，prometheus等等库来进行性能监控等等的功能

## PProf

使用net/http/pprof可以能方便的采集web服务在运行时的数据，直接import十分简单

```golang
import (
_ "net/http/pprof"
)

// 看看pprof的init发现注入了很多handler
func init() {
http.HandleFunc("/debug/pprof/", Index)
http.HandleFunc("/debug/pprof/cmdline", Cmdline)
http.HandleFunc("/debug/pprof/profile", Profile)
http.HandleFunc("/debug/pprof/symbol", Symbol)
http.HandleFunc("/debug/pprof/trace", Trace)
}

```

可以通过浏览器或是交互式终端进行访问，我选择浏览器。-> IP地址/debug/pprof

浏览器有时效性，真要查问题还是用终端

```golang
go tool pprof ip/debug/pprof/profile?seconds=60
```

不同的路由对应的项，比如cpu/heap/goroutine等等，这就不贴了。

采集生成的profile文件也是可以用web方式查阅的，go tool pprof -http=:6000 profile，如果报graphviz的错说明要装组件。

通过Lookup进行采集,这种方式需要写code，支持6种类型，goroutine,threadcreate,heap,block,mutex

<details>
<summary> 这种方式需要写code，支持6种类型，goroutine,threadcreate,heap,block,mutex </summary>

```golang
package main

import (
"io"
"net/http"
_ "net/http/pprof"
"os"
"runtime"
"runtime/pprof"
)

func main() {
http.HandleFunc("/lookup/heap", func(w http.ResponseWriter, r *http.Request) {
_ = pprofLookup(LookupHeap, os.Stdout)
})

http.HandleFunc("/lookup/threadcreate", func(w http.ResponseWriter, r *http.Request) {
_ = pprofLookup(LookupThreadcreate, os.Stdout)
})

http.HandleFunc("/lookup/block", func(w http.ResponseWriter, r *http.Request) {
_ = pprofLookup(LookupBlock, os.Stdout)
})

http.HandleFunc("/lookup/goroutine", func(w http.ResponseWriter, r *http.Request) {
_ = pprofLookup(LookupGoroutine, os.Stdout)
})

_ = http.ListenAndServe("0.0.0.0:6060", nil)
}

type LookupType int8

const (
LookupGoroutine LookupType = iota
LookupThreadcreate
LookupHeap
LookupAllocs
LookupBlock
LookupMutex
)

func pprofLookup(lookupType LookupType, w io.Writer) error {
var err error
switch lookupType {
case LookupGoroutine:
p := pprof.Lookup("goroutine")
err = p.WriteTo(w, 2)
case LookupThreadcreate:
p := pprof.Lookup("threadcreate")
err = p.WriteTo(w, 2)
case LookupHeap:
p := pprof.Lookup("heap")
err = p.WriteTo(w, 2)
case LookupAllocs:
p := pprof.Lookup("allocs")
err = p.WriteTo(w, 2)
case LookupBlock:
p := pprof.Lookup("block")
err = p.WriteTo(w, 2)
case LookupMutex:
p := pprof.Lookup("mutex")
err = p.WriteTo(w, 2)
}
return err
}

func init() {
runtime.SetMutexProfileFraction(1)
runtime.SetBlockProfileRate(1)
}
```

</details>

## trace

<details>
<summary> 详细的使用方式 </summary>

```golang
package main

// --go run .\cmd\trace\main.go 2> trace.out

// go build .\cmd\trace\main.go
//  .\main.exe
// go tool trace trace.dat

import (
"context"
"fmt"
"os"
"runtime"
"runtime/trace"
"sync"
)

func main() {
// 为了看协程抢占，这里设置了一个cpu 跑
runtime.GOMAXPROCS(1)

f, _ := os.Create("trace.dat")
defer f.Close()

_ = trace.Start(f)
defer trace.Stop()

ctx, task := trace.NewTask(context.Background(), "sumTask")
defer task.End()

var wg sync.WaitGroup
wg.Add(10)
for i := 0; i < 10; i++ {
// 启动10个协程，只是做一个累加运算
go func(region string) {
defer wg.Done()

// 标记region
trace.WithRegion(ctx, region, func() {
var sum, k int64
for ; k < 1000000000; k++ {
sum += k
}
fmt.Println(region, sum)
})
}(fmt.Sprintf("region_%02d", i))
}
wg.Wait()
}
```

</details>

## godebug

这个同样太长不写了。环境变量可以通过vscode写入launch.json文件中，比如

```json
  "env": {
"GODEBUG":"scheddetail=1,schedtrace=1000"
} 
```

## 进程诊断工具 gops

```golang
go get -u github.com/google/gops 

package main

import (
"log"
"net/http"

"github.com/google/gops/agent"
)

func main() {
if err := agent.Listen(agent.Options{}); err != nil {
log.Fatal("agent listen err : %v", err)
}
http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
_, _ = w.Write([]byte("golang projecct"))
})
_ = http.ListenAndServe(":6060", http.DefaultServeMux)
}

```

gops help查看命令，也是有很多，不细写

## metrics 使用expvar标准库

<details>
<summary> code中自定义了类型，并且封装成了gin中间件，可以和gin联动了 </summary>

```golang
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

func main() {
router := NewRouter()
http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
appleCounter.Add(1)
_, _ = w.Write([]byte("go project"))
})

_ = http.ListenAndServe(":6060", router)
}

var (
appleCounter  *expvar.Int
GOMAXPROCSMetrics *expvar.Int
upTimeMetrice *upTimeVar
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
```

</details>
  
## Pronmetheus

Pronmetheus还是很出名的。四大指标类型Counter累计指标，Histogram一定时间范围内采样，
Gauge可任意变化的指标，Summary也是一定时间内采样，他有仨指标，分位数分布/样本值大小总和/样本总数。

```golang
go get -u github.com/prometheus/client_golang

package main

import (
"net/http"

"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
http.Handle("/metrics", promhttp.Handler())
http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
_, _ = w.Write([]byte("go project"))
})
_ = http.ListenAndServe(":6060", http.DefaultServeMux)
}

```

启动后访问 :6060/metrics

## 其他  包括附录

逃逸分析，有很多情况会造成逃逸，这个还是要经验的，初学者的我还是用命令最直接。

```golang
// 用-gcflags查看逃逸分析过程
go build -gcflags '-m -l' main.go

// 反编译命令查看
go tool compile -S main.go
```

Go modules

```golang
// go get后的模块会缓存在gopath/pkg/mod 和gopath/pkg/sumdb中，如果需要清理可以执行
go clean -modcache
```

为什么defer才能recover

panic结构是一个链表，defer结构中包含了一个对panic结构的引用，在gopanic(interface{})方法中，
会触发defer，如果没有defer则会直接跳出，就不会进行接来下的recover了。

还有一些defer了也无法recover的方法，比如fatalthrow,fatalpanic等，比如并发写入map时就会引起fatalthrow。

10种panic方法：数组切片越界，空指针调用，过早关闭HTTP响应体(resp.body.calose())，除零，
向关闭的chan发送消息，重复关闭chan，关闭未初始化的的chan，使用未初始化的map，跨goroutine处理panic，sync计数负数。

让golang更适应docker

```golang
// 这个库可以根据cgroup的挂载信息来修改GOMAXPROCS核数
import _ "go.uber.org/automaxprocs"
```
