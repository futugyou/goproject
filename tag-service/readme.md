# 第三章： grpc服务

## 简介

装protobuf插件,[protobuf](https://github.com/protocolbuffers/protobuf/releases)

```golang
 go get -u github.com/golang/protobuf/protoc-gen-go
 protoc --go_out=plugins=grpc:. ./proto/*.proto
```

简单说下grpc四种模式都是client发起，server响应。

1. 一元的没啥
2. server流式: client读到EOF时结束
3. client流式: client的Send(r)发完后调用CloseAndRecv()等待server的返回，server端Recv()，
直到EOF后调用SendAndClose(r)通知client结束

    ```TL
    sequenceDiagram
    client->>server: Send(r) 
    client->>server: Send(r) 
    client->>server: Send(r) 
    server->>server: Recv()
    server->>client: SendAndClose(r)
    client->>client: CloseAndRecv()
    ```

4. 双向流式: 双方都是流式的，server没有特别的方法，client结束时调用CloseSend()

详细信息可以翻官网还有pb.go文件

装protobuf插件  

```golang
 go get -u github.com/golang/protobuf/protoc-gen-go
```

## 写个grpc服务

<details>
<summary> 写个protobuf文件tag.proto </summary>

```protobuf
syntax = "proto3";

package proto;

import "google/api/annotations.proto";

service TagService {
rpc GetTagList (GetTagListRequest) returns (GetTagListReply) {
option (google.api.http) = {
get: "/api/v1/tags"
};
};
}

message  GetTagListRequest {
string name = 1;
uint32 state = 2;
}

message GetTagListReply {
repeated Tag list = 1;
Pager pager = 2;
}

message Tag {
int64  id = 1;
string name = 2;
uint32 state = 3;
}

message Pager {
int64 page = 1;
int64 page_size = 2;
int64 totle_rows = 3;
}
```

</details>

service中间option的部分是后面grpcgateway用的

<details>
<summary> 看下执行protoc命令后生成的tag.pb.go中的和server相关的code </summary>

```golang
 
type TagServiceServer interface {
GetTagList(context.Context, *GetTagListRequest) (*GetTagListReply, error)
}

// UnimplementedTagServiceServer can be embedded to have forward compatible implementations.
type UnimplementedTagServiceServer struct {
}

func (*UnimplementedTagServiceServer) GetTagList(context.Context, *GetTagListRequest) (*GetTagListReply, error) {
return nil, status.Errorf(codes.Unimplemented, "method GetTagList not implemented")
}

func RegisterTagServiceServer(s *grpc.Server, srv TagServiceServer) {
s.RegisterService(&_TagService_serviceDesc, srv)
}

```

</details>

<details>
<summary> 我们需要定一个struct 实现下 TagServiceServer这个接口</summary>

```golang
type TagServer struct {
}

func NewTagServer() *TagServer {
return &TagServer{}
}

func (t *TagServer) GetTagList(ctx context.Context, r *pb.GetTagListRequest) (*pb.GetTagListReply, error) {
...
}

```

</details>

<details>
<summary> main.go 里面注册下 </summary>

```golang
func main {
s := grpc.NewServer()
pb.RegisterTagServiceServer(s, server.NewTagServer())
reflection.Register(s) // 这句是使用下面的grpcurl所必需的
lis,err := net.Listen("tcp", ":8001")
...
err = s.Serve(lis)
...
}
```

</details>

## 弄个工具grpcurl验证下

```golang
// 安装grpcurl命令行工具
go get -u github.com/fullstorydev/grpcurl

// list命令可以显示有多少grpc服务可以调用，然后可以通过具体服务名进行调用
grpcurl -plaintext localhost:8001 list
grpcurl -plaintext localhost:8001 proto.TagService.GetTagList
```

## 写个client

<details>
<summary> 先看看tag.pb.go里面的client相关的code </summary>

```golang
type TagServiceClient interface {
GetTagList(ctx context.Context, in *GetTagListRequest, opts ...grpc.CallOption) (*GetTagListReply, error)
}

type tagServiceClient struct {
cc grpc.ClientConnInterface
}

func NewTagServiceClient(cc grpc.ClientConnInterface) TagServiceClient {
return &tagServiceClient{cc}
}

func (c *tagServiceClient) GetTagList(ctx context.Context, in *GetTagListRequest, opts ...grpc.CallOption) (*GetTagListReply, error) {
out := new(GetTagListReply)
err := c.cc.Invoke(ctx, "/proto.TagService/GetTagList", in, out, opts...)
if err != nil {
return nil, err
}
return out, nil
}
```

</details>

<details>
<summary> 写个client.go </summary>

```golang
func main() {
ctx := context.Background()
conn, err := grpc.DialContext(ctx, "localhost:8001", nil)
defer conn.Close()
client := pb.NewTagServiceClient(conn)
resp, err := client.GetTagList(newCtx, &pb.GetTagListRequest{Name: "golang"})
}
```

</details>

还是.net的grpc好些，server和client生成的代码不会都混在一起

## 使用cmux同时支持Http/grpc

cmux 可以做到在同一个端口监听grpc和http。实际上就是在同一tcplistener上进行多路复用

```shell
go get -u github.com/soheilhy/cmux
```

<details>
<summary> 1. 初始化tcp listener  2.content-type为application/grpc</summary>

```golang

func main() {
l, err := RunTCPServer(port)
if err != nil {
log.Fatalf("run tcp server err: %v", err)
}
m := cmux.New(l)
grpcL := m.MatchWithWriters(
cmux.HTTP2MatchHeaderFieldPrefixSendSettings(
"content-type",
"application/grpc"))
httpL := m.Match(cmux.HTTP1Fast())

grpcS := RunGrpcServer()
httpS := RunHttpServer(port)

go grpcS.Serve(grpcL)
go httpS.Serve(httpL)
err = m.Serve()
if err != nil {
log.Fatalf("run serve err: %v", err)
}
}

func RunTCPServer(port string) (net.Listener, error) {
return net.Listen("tcp", ":"+port)
}

func RunGrpcServer() *grpc.Server {
s := grpc.NewServer()
pb.RegisterTagServiceServer(s, server.NewTagServer())
reflection.Register(s)
return s
}

func RunHttpServer(port string) *http.Server {
serveMux := http.NewServeMux()
serveMux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
_, _ = w.Write([]byte(`pong`))
})
return &http.Server{
Addr:":" + port,
Handler: serveMux,
}
}
```

</details>

## 使用grpc-gateway同时支持Http/grpc

gateway生成反向代理，可将restful api转为grpc，根据protobuf中的 google.api.http。具体看上面的protobuf文件

```shell
go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway

// 这次命令很长， 因为路径长，还有外加了swagger，下面会写
protoc -IC:\Users\TerraformRs\Documents\GitHub\
protobuf-3.13.0\include -I. -IC:\Users\TerraformRs\go -IC:\Users\TerraformRs\go\pkg\mod\github.com\
grpc-ecosystem\grpc-gateway@v1.15.2\third_party\
googleapis --go_out=plugins=grpc:. --swagger_out=logtostderr=true:.  --grpc-gateway_out=logtostderr=true:. ./proto/*.proto

// 这次执行完后除了原本的tag.pb.go外还会生成一个tag.pb.gw.go文件
```

<details>
<summary> server端code，可以看到主要添加了对gateway的注册，重要的就是RegisterTagServiceHandlerFromEndpoint这个方法。
具体的流程可以看tag.pb.gw.go这个文件。 </summary>

```golang
func main() {
err := RunServer(port)
if err != nil {
log.Fatalf("run tcp server err: %v", err)
}
}

func RunServer(port string) error {
httpMux := runHttpServer()
gatewayMux := runGrpcGatewayServer()
grpcS := runGrpcServer()
httpMux.Handle("/", gatewayMux)
return http.ListenAndServe(":"+port, grpcHandlerFunc(grpcS, httpMux))
}

func runHttpServer() *http.ServeMux {
serverMux := http.NewServeMux()
serverMux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
_, _ = w.Write([]byte(`pong`))
})
return serverMux
}

func runGrpcGatewayServer() *runtime.ServeMux {
endpoint := "0.0.0.0:" + port
runtime.HTTPError = grpcGatewayError
gwmux := runtime.NewServeMux()
dopts := []grpc.DialOption{grpc.WithInsecure()}
_ = pb.RegisterTagServiceHandlerFromEndpoint(context.Background(), gwmux, endpoint, dopts)
return gwmux
}

func runGrpcServer() *grpc.Server {
opts := []grpc.ServerOption{
grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
...
)),
}
s := grpc.NewServer(opts...)
pb.RegisterTagServiceServer(s, server.NewTagServer())
reflection.Register(s)
return s
}
```

</details>

## swagger

这次使用swagger和上个项目不同，需要从官网下一些文件，并打包。

```shell
go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
go get -u github.com/go-bindata/go-bindata/...

// 然后打包成data.go文件
go-bindata --nocompress -pkg swagger -o pkg/swagger/data.go third_party/swagger-ui/...

// 还要结合go-bindata-assetfs 
go get -u github.com/elazarl/go-bindata-assetfs 
```

<details>
<summary> 修改main.go的runHttpServer方法，注册swagger的路由 </summary>

```golang
func runHttpServer() *http.ServeMux {
serverMux := http.NewServeMux()
serverMux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
_, _ = w.Write([]byte(`pong`))
})
prefix := "/swagger-ui/"
fileServer := http.FileServer(&assetfs.AssetFS{
Asset:swagger.Asset,
AssetDir: swagger.AssetDir,
Prefix:   "third_party/swagger-ui",
})
serverMux.Handle(prefix, http.StripPrefix(prefix, fileServer))
serverMux.HandleFunc("/swagger/", func(w http.ResponseWriter, r *http.Request) {
if !strings.HasSuffix(r.URL.Path, "swagger.json") {
http.NotFound(w, r)
return
}
p := strings.TrimPrefix(r.URL.Path, "/swagger/")
p = path.Join("proto", p)
http.ServeFile(w, r, p)
})
return serverMux
}
```

</details>

protoc命令前面已经写了，执行后会生成比如tag.swagger.json这样大的文件。运行后先访问/swagger-ui，
然后指定/swagger/tag.swagger.json Explore

## 拦截器

client和server都有内置的拦截器，还分为一元和流式的。为了达到同时使用多个拦截器，需要使用go-grpc-middleware

```shell
go get -u github.com/grpc-ecosystem/go-grpc-middleware 
```

<details>
<summary> 先写两个拦截器</summary>

```golang
// 这是server的UnaryServerInfo 一元拦截器，这个参数换成StreamServerInfo就是流式的
func AccessLog(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
requestLog := "access request log: method: %s, begin_time: %d, request: %v"
beginTime := time.Now().Local().Unix()
log.Printf(requestLog, info.FullMethod, beginTime, req)
resp, err := handler(ctx, req)
responseLog := "access response log:method: %s,begin_time: %d, end_time: %d,response: %v"
endTime := time.Now().Local().Unix()
log.Printf(responseLog, info.FullMethod, beginTime, endTime, resp)
return resp, err
}

// 这是一个client的流式拦截器
func StreamContextTimeout() grpc.StreamClientInterceptor {
return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
ctx, cancel := defaultContextTimeout(ctx)
if cancel != nil {
defer cancel()
}
return streamer(ctx, desc, cc, method, opts...)
}
}

func defaultContextTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
var cancel context.CancelFunc
if _, ok := ctx.Deadline(); !ok {
defaultTimeout := 60 * time.Second
ctx, cancel = context.WithTimeout(ctx, defaultTimeout)
}
return ctx, cancel
}
```

</details>

<details>
<summary> 注册拦截器</summary>

```golang
// server的
func runGrpcServer() *grpc.Server {
opts := []grpc.ServerOption{
grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
middleware.AccessLog,
)),
}
s := grpc.NewServer(opts...)
pb.RegisterTagServiceServer(s, server.NewTagServer())
reflection.Register(s)
return s
}

// client的
opts = append(opts, grpc.WithStreamInterceptor(
grpc_middleware.ChainStreamClient(
middleware.StreamContextTimeout(),
),
))
return grpc.DialContext(ctx, target, opts...)
```

</details>

## tracing

tracing和jaeger可以和上个项目联动

```shell
go get -u github.com/opentracing/opentracing-go 
go get -u github.com/uber/jaeger-client-go

```

grpc可以借助metadata来完成tracing 。server拦截器可从metadata提取信息，并追加到context。
client拦截器从context中提取信息，并作为metadata加入到grpc调用中

<details>
<summary> 封装一个metadata </summary>

```golang
type MetadataTextMap struct {
metadata.MD
}

func (m MetadataTextMap) Set(key, val string) {
key = strings.ToLower(key)
m.MD[key] = append(m.MD[key], key)
}

func (m MetadataTextMap) ForeachKey(handler func(key, val string) error) error {
for k, vs := range m.MD {
for _, v := range vs {
if err := handler(k, v); err != nil {
return err
}
}
}
return nil
}

```

</details>

<details>
<summary> server拦截器 </summary>

```golang
func ServerTracing(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
md, ok := metadata.FromIncomingContext(ctx)
if !ok {
md = metadata.New(nil)
}

parentSpanContext, _ := global.Tracer.Extract(opentracing.TextMap, metatext.MetadataTextMap{md})
spanOpts := []opentracing.StartSpanOption{
opentracing.Tag{Key: string(ext.Component), Value: "grpc"},
ext.SpanKindRPCServer,
ext.RPCServerOption(parentSpanContext),
}

span := global.Tracer.StartSpan(info.FullMethod, spanOpts...)
defer span.Finish()
ctx = opentracing.ContextWithSpan(ctx, span)

return handler(ctx, req)
}

```

</details>

<details>
<summary> client拦截器 </summary>

```golang
func ClientTracing() grpc.UnaryClientInterceptor {
return func(ctx context.Context, method string, req, resp interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
var parentCtx opentracing.SpanContext
var spanOpts []opentracing.StartSpanOption
var parentSpan = opentracing.SpanFromContext(ctx)
if parentSpan != nil {
parentCtx = parentSpan.Context()
spanOpts = append(spanOpts, opentracing.ChildOf(parentCtx))
}
spanOpts = append(spanOpts, []opentracing.StartSpanOption{
opentracing.Tag{Key: string(ext.Component), Value: "grpc"},
ext.SpanKindRPCClient,
}...)
span := global.Tracer.StartSpan(method, spanOpts...)
defer span.Finish()

md, ok := metadata.FromOutgoingContext(ctx)
if !ok {
md = metadata.New(nil)
}

_ = global.Tracer.Inject(span.Context(), opentracing.TextMap, metatext.MetadataTextMap{md})
newCtx := opentracing.ContextWithSpan(metadata.NewOutgoingContext(ctx, md), span)
return invoker(newCtx, method, req, resp, cc, opts...)
}
}

```

</details>

<details>
<summary> Http追踪，这个方法前面没提到，是用来调用上一个blog项目的的相关接口 </summary>

```golang
func (a *API) httpGet(ctx context.Context, path string) ([]byte, error) {
url := fmt.Sprintf("%s/%s", a.URL, path)
req, err := http.NewRequest("GET", url, nil) 
if err != nil {
return nil, err
}
span, newCtx := opentracing.StartSpanFromContext(
ctx,
"HTTP GET: "+a.URL,
opentracing.Tag{Key: string(ext.Component), Value: "HTTP"},
)
span.SetTag("url", url)

_ = opentracing.GlobalTracer().Inject(
span.Context(),
opentracing.HTTPHeaders,
opentracing.HTTPHeadersCarrier(req.Header),
)
req = req.WithContext(newCtx)
client := http.Client{Timeout: time.Second * 60}
resp, err := client.Do(req)
if err != nil {
return nil, err
}

defer resp.Body.Close()
defer span.Finish()
body, err := ioutil.ReadAll(resp.Body)
if err != nil {
return nil, err
}
return body, nil
}
```

</details>
