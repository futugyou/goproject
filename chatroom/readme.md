# 第四章： websocket服务

## 简介

我写笔记的目的是为了记录所用到的一些组件，但是从这章开始到结束组件不像前面多了。

ws可在单个tcp连接上建立全双工通讯，允许服务端主动向客户端输出。保持连接状态，是一种有状态的应用层协议。
建立方式通过可以http代理来握手，使用HTTP Upgrade头进行协议升级。

## 库选择

书上写了多种库，包括nhooyr.io/websocket和gorilla/wesocket等，后者成熟功能多，
前者更符合go的习惯并且有并发，完整关闭，可以编译为wasm等好处。

```golang
go get -u nhooyr.io/websocket
```

## 基本用法

<details>
<summary> server.go </summary>

```golang
package main

import (
 "context"
 "fmt"
 "log"
 "net/http"
 "time"

 "nhooyr.io/websocket"
 "nhooyr.io/websocket/wsjson"
)

func main() {
 http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintln(w, "HTTP Hello")
 })

 http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
  conn, err := websocket.Accept(w, r, nil)
  if err != nil {
   log.Println(err)
   return
  }
  defer conn.Close(websocket.StatusInternalError, "internal error")
  ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
  defer cancel()

  var v interface{}
  err = wsjson.Read(ctx, conn, &v)
  if err != nil {
   log.Println(err)
   return
  }
  log.Printf("client resvered: %v\n", v)

  err = wsjson.Write(ctx, conn, "Hello websocket client.")
  if err != nil {
   log.Println(err)
   return
  }

  conn.Close(websocket.StatusNormalClosure, "")
 })

 log.Fatal(http.ListenAndServe(":2021", nil))
}

```

</details>

<details>
<summary> client.go </summary>

```golang
package main

import (
 "context"
 "fmt"
 "time"

 "nhooyr.io/websocket"
 "nhooyr.io/websocket/wsjson"
)

func main() {
 ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
 defer cancel()

 c, _, err := websocket.Dial(ctx, "ws://localhost:2021/ws", nil)
 if err != nil {
  panic(err)
 }
 defer c.Close(websocket.StatusInternalError, "internal error")
 err = wsjson.Write(ctx, c, "Hello websocket server")
 if err != nil {
  panic(err)
 }

 var v interface{}
 err = wsjson.Read(ctx, c, &v)
 if err != nil {
  panic(err)
 }
 fmt.Printf("resvered server's response : %v\n", v)
 c.Close(websocket.StatusNormalClosure, "")
}
```

</details>

## 其他

github和博客园都支持不了md的图

```flowchart TD
sequenceDiagram
client->>home/: HTTP / 请求首页
home/->>client: 聊天室页面
client->>webscoket: 通过ws协议请求： /ws
webscoket->>webscoket: 协议转换
webscoket->>client: 响应
webscoket->broadcast: channel通信
webscoket->client: 通过ws长连接发送消息
```

代码没有特别的，不仔细贴了，基本就是channel的操作，完整的还是去github。

<details>
<summary> 一段单例的代码，双检锁基本都长这样，
不过没有看到‘易变’。golang有没有这个东西我不记得了 </summary>

```golang
package singleton

import "sync"

type singleton2 struct {
 count int
}

var (
 instance2 *singleton2
 mutex sync.Mutex
)

func New() *singleton2 {
 if instance2 == nil {
  mutex.Lock()
  if instance2 == nil {
   instance2 = new(singleton2)
  }
  mutex.Unlock()
 }
 return instance2
}
func (s *singleton2) Add() int {
 s.count++
 return s.count
}

```

</details>

在敏感词处理环节介绍了DFA和贝叶斯分类算法，提到了两个库

```golang
go get -u github.com/antlinker/go-dirtytilter
go get -u github.com/jbrukh/bayesian
```

用户识别用的token生成方式和之前章节一样，离线消息处理使用了内置的container/ring环形链表
