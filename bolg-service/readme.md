# 第二章： 一个http应用blogservice

## gin web框架

```shell
go get -u github.com/gin-gonic/gin
```

<details>
<summary>gin初始化的代码片段</summary>

```golang
gin.SetMode(global.ServerSetting.RunModel) // DEBUG
 router := routers.NewRouter()
 s := &http.Server{
  Addr:   ":" + global.ServerSetting.HttpPort,
  Handler:router,
  ReadTimeout:global.ServerSetting.ReadTimeout,
  WriteTimeout:   global.ServerSetting.WriteTimeout,
  MaxHeaderBytes: 1 << 20,
 }
 go func() {
  if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
   log.Fatalf("s. listenandserve err: %v", err)
  }
 }()
```

</details>

<details>
<summary>配置中间件及路由</summary>

```golang
func NewRouter() *gin.Engine {
 r := gin.New()
 r.Use(middleware.Tracing())
 r.Use(middleware.AccessLog()) //原始的Logger()和Recovery()已经被替换了
 r.Use(middleware.Recovery())
 r.Use(middleware.RateLimiter(methodLimiters))
 r.Use(middleware.ContextTimeout(60 * time.Second))
 r.Use(middleware.Translations())

 r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
 r.GET("/auth", api.GetAuth)
 tag := v1.NewTage()
 article := v1.NewArticle()
 r.POST("/upload/file", api.UploadFile)
 r.StaticFS("/static", http.Dir(global.AppSetting.UploadSavePath))
 apiv1 := r.Group("/api/v1")
 {
  apiv1.Use(middleware.JWT())
  apiv1.POST("/tags", tag.Create)
  ...
 }
 return r
}
```

</details>

<details>
<summary>中间件写法举例，习惯写aspnetcore的人会觉得很亲切</summary>

```golang
func Recovery() gin.HandlerFunc {
 return func(c *gin.Context) {
  defer func() {
   if err := recover(); err != nil {
    global.Logger.WithCallersFrames().Errorf(c, "panic recover err : %v", err)
    app.NewResponse(c).ToErrorResponse(errcode.ServiceError)
    c.Abort()
   }
  }()
  c.Next()
 }
}
```

</details>

## viper 读配置文件

```shell
go get -u github.com/spf13/viper
go get -u golang.org/x/sys/...   // 下面这俩做热更新用
go get -u github.com/fsnotify/fsnotify
```

<details>
<summary>viper初始化的代码片段</summary>

```golang
type Setting struct {
 vp *viper.Viper
}

func NewSetting(configs ...string) (*Setting, error) {
 vp := viper.New()
 vp.SetConfigName("config") //文件名
 for _, config := range configs {
  if config != "" {
   vp.AddConfigPath(config) //查找路径
  }
 }
 vp.SetConfigType("yaml") //文件类型
 err := vp.ReadInConfig()
 if err != nil {
  return nil, err
 }
 s := &Setting{vp}
 s.WatchSettingChange()
 return s, nil
}

// 监视变更
func (s *Setting) WatchSettingChange() {
 go func() {
  s.vp.WatchConfig()
  s.vp.OnConfigChange(func(in fsnotify.Event) {
   _ = s.ReloadAllSection()
  })
 }()
}

// 保存各个sections 主要是为了做热更新而存在的
var sections = make(map[string]interface{})

// 读取单个section`k`的数据至对象`v`中
func (s *Setting) ReadSection(k string, v interface{}) error {
 err := s.vp.UnmarshalKey(k, v)
 if err != nil {
  return err
 }
 if _, ok := sections[k]; !ok {
  sections[k] = v
 }
 return nil
}

func (s *Setting) ReloadAllSection() error {
 for k, v := range sections {
  err := s.ReadSection(k, v)
  if err != nil {
   return err
  }
 }
 return nil
}
```

</details>
<details>
<summary>调用方式 比如main.go的init()中</summary>

```golang
func init() {
 err := setupSetting()
 ...
 }
 
func setupSetting() error {
 setting, err := setting.NewSetting(strings.Split(config, ",")...)
 if err != nil {
  return err
 }
 err = setting.ReadSection("Server", &global.ServerSetting)
 if err != nil {
  return err
 }
 ...
 return nil
}
```

</details>

## gorm 看名字就知道是orm了

```shell
go get -u github.com/jinzhu/gorm
```

<details>
<summary>gorm初始化的代码片段</summary>

```golang
// mysql的
func NewDBEngine(dbsetting *setting.DatabaseSettingS) (*gorm.DB, error) {
 db, err := gorm.Open(dbsetting.DBType,
  fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=Local",
   dbsetting.Username,
   dbsetting.Password,
   dbsetting.Host,
   dbsetting.DBName,
   dbsetting.Charset,
   dbsetting.ParseTime,
  ))
 if err != nil {
  return nil, err
 }
 if global.ServerSetting.RunModel == "debug" {
  db.LogMode(true)
 }
 db.SingularTable(true)
 // 自动迁移
 if db.HasTable(&Tag{}) {
  db.AutoMigrate(&Tag{})
 } else {
  db.CreateTable(&Tag{})
 }
 ...

// 增删改的回调 用于处理公共字段
 db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
 db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
 db.Callback().Delete().Replace("gorm:delete", deleteCallback)

 db.DB().SetMaxIdleConns(dbsetting.MaxIdleConns)
 db.DB().SetMaxOpenConns(dbsetting.MaxOpenConns)
 // 这个是tracing相关 来源自github.com/eddycjy/opentracing-gorm
 otgorm.AddGormCallbacks(db)
 return db, nil
}
```

</details>

<details>
<summary>回调处理公共字段 有些麻烦</summary>

```golang
func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
 if _, ok := scope.Get("gorm:update_column"); !ok {
  _ = scope.SetColumn("ModifiedOn", time.Now().Unix())
 }
}

func updateTimeStampForCreateCallback(scope *gorm.Scope) {
 if !scope.HasError() {
  nowTime := time.Now().Unix()
  if createTimeField, ok := scope.FieldByName("CreatedOn"); ok {
   if createTimeField.IsBlank {
    _ = createTimeField.Set(nowTime)
   }
  }
  if modifyTimeField, ok := scope.FieldByName("ModifiedOn"); ok {
   if modifyTimeField.IsBlank {
    _ = modifyTimeField.Set(nowTime)
   }
  }
 }
}

func deleteCallback(scope *gorm.Scope) {
 if !scope.HasError() {
  var extraoption string
  if str, ok := scope.Get("gorm:delete_option"); ok {
   extraoption = fmt.Sprint(str)
  }
  deleteOnField, hasDeletedonField := scope.FieldByName("DeletedOn")
  isDelFiled, hasIsDelField := scope.FieldByName("IsDel")
  if !scope.Search.Unscoped && hasDeletedonField && hasIsDelField {
   now := time.Now().Unix()
   scope.Raw(fmt.Sprintf(
    "update %v set %v=%v ,%v=%v%v%v",
    scope.QuotedTableName(),
    scope.Quote(deleteOnField.DBName),
    scope.AddToVars(now),
    scope.Quote(isDelFiled.DBName),
    scope.AddToVars(1),
    addExtraSpaceIfExist(scope.CombinedConditionSql()),
    addExtraSpaceIfExist(extraoption),
   )).Exec()
  } else {
   scope.Raw(fmt.Sprintf(
    "delete from %v%v%v",
    scope.QuotedTableName(),
    addExtraSpaceIfExist(scope.CombinedConditionSql()),
    addExtraSpaceIfExist(extraoption),
   )).Exec()
  }
 }
}

func addExtraSpaceIfExist(str string) string {
 if str != "" {
  return " " + str
 }
 return ""
}
```

</details>

<details>
<summary>curd来一套</summary>

```golang
func (t Tag) Count(db *gorm.DB) (int, error) {
 var count int
 if t.Name != "" {
  db = db.Where("name = ?", t.Name)
 }
 db = db.Where("state = ?", t.State)
 if err := db.Model(&t).Where("is_del = ?", 0).Count(&count).Error; err != nil {
  return 0, err
 }
 return count, nil
}

func (t Tag) Get(db *gorm.DB) (Tag, error) {
 var tag Tag
 db = db.Where("id = ? and is_del = ? and state = ?", t.ID, 0, t.State)
 if err := db.Model(&t).Where("is_del = ?", 0).Find(&tag).Error; err != nil {
  return tag, err
 }
 return tag, nil
}

func (t Tag) List(db *gorm.DB, pageoffset, pagesize int) ([]*Tag, error) {
 var tags []*Tag
 var err error
 if pageoffset >= 0 && pagesize > 0 {
  db = db.Offset(pageoffset).Limit(pagesize)
 }
 if t.Name != "" {
  db = db.Where("name = ?", t.Name)
 }
 db = db.Where("state = ?", t.State)
 if err = db.Where("is_del = ?", 0).Find(&tags).Error; err != nil {
  return nil, err
 }
 return tags, nil
}

func (t Tag) Create(db *gorm.DB) error {
 return db.Create(&t).Error
}

func (t Tag) Update(db *gorm.DB, value interface{}) error {
 return db.Model(&t).Where("id = ? and is_del = ?", t.ID, t.IsDel).Updates(value).Error
}

func (t Tag) Delete(db *gorm.DB) error {
 return db.Where("id = ? and is_del = ?", t.Model.ID, 0).Delete(&t).Error
}

```

</details>

使用orm的目的在于完全屏蔽数据库细节，解放生产力。但书中举例的很多操是在拼sql及硬编码字段名，尤其是where条件里面和用于处理公共字段的回调方法，
这完全不是orm的理念。是这本书的作者使用方式不对么？

## lumberjack 写日志的

```shell
go get -u gopkg.in/natefinch/lumberjack.v2
```

<details>
<summary>作者对log做了大量的封装，使用了很多runtime的方法来获取数据，
具体还是看源码吧，我只记下这个组件最基本的用法</summary>

```golang
func setupLogger() error {
 global.Logger = logger.NewLogger(
  &lumberjack.Logger{
   Filename:  global.AppSetting.LogSavaPath + "/" + global.AppSetting.LogFileName + global.AppSetting.LogFileExt,
   MaxSize:   600,
   MaxAge:10,
   LocalTime: true,
  },
  "",
  log.LstdFlags,
 ).WithCaller(2)
 return nil
}

type Level int8
type Fields map[string]interface{}  

type Logger struct {
 newLogger *log.Logger
 ctx   context.Context
 level Level
 fieldsFields
 callers   []string
}

func NewLogger(w io.Writer, prefix string, flag int) *Logger {
 l := log.New(w, prefix, flag)
 return &Logger{newLogger: l}
}

func (l *Logger) WithCaller(skip int) *Logger {
 ll := l.clone()
 pc, file, line, ok := runtime.Caller(skip)
 if ok {
  f := runtime.FuncForPC(pc)
  ll.callers = []string{fmt.Sprintf("%s: %d %s", file, line, f.Name())}
 }
 return ll
}
```

</details>

## swagger

```shell
go get -u github.com/swaggo/gin-swagger 
go get -u github.com/swaggo/swag 
go get -u github.com/alecthomas/template
swag -v
```

<details>
<summary>手写注释很麻烦 写完后swag init就能生成，路由在前面的gin环节已经配置过了</summary>

```golang
// @Summary get tags
// @Produce json
// @Param name query string false "tag name" maxlength(100)
// @Param state query int false "state" Enums(0,1) defaulrt(1)
// @Param page query int false "page"
// @Param page_size query int false "page_size"
// @Success 200 {object} model.TagSwagger "success"
// @Failure 400 {object} errcode.Error "error"
// @Failure 500 {object} errcode.Error "error"
// @Router /api/v1/tags [get]
func (a Tag) List(c *gin.Context) {
...
}
```

</details>

## validator 接口验证器

```shell
go get -u github.com/go-playground/validator/v10 
go get -u github.com/go-playground/locales //多语言包
go get -u github.com/go-playground/universal-translator 翻译器
```

gin中默认的验证就是用的这个组件，现在对其作一定的定制，并添加到中间件中，参照gin环节

<details>
<summary>先写个中间件，主要做国际化和验证器注册</summary>

```golang
func Translations() gin.HandlerFunc {
 return func(c *gin.Context) {
  uni := ut.New(en.New(), zh.New(), zh_Hant_TW.New())
  locale := c.GetHeader("locale")
  trans, _ := uni.GetTranslator(locale)
  v, ok := binding.Validator.Engine().(*Validator.Validate)
  if ok {
   switch locale {
   case "zh":
    _ = zh_translations.RegisterDefaultTranslations(v, trans)
    break
   case "en":
    _ = en_translations.RegisterDefaultTranslations(v, trans)
    break
   default:
    _ = zh_translations.RegisterDefaultTranslations(v, trans)
    break
   }
   c.Set("trans", trans)
  }
  c.Next()
 }
}

```

</details>

<details>
<summary>再对gin的ShouldBind封装个，用中间件中的"trans"对错误国际化</summary>

```golang
func BindAndValid(c *gin.Context, v interface{}) (bool, ValidErrors) {
 var errs ValidErrors
 err := c.ShouldBind(v)
 if err != nil {
  t := c.Value("trans")
  trans, _ := t.(ut.Translator)
  verrs, ok := err.(val.ValidationErrors)
  if !ok {
   return false, nil
  }
  for key, value := range verrs.Translate(trans) {
   errs = append(errs, &ValidError{
    Key: key,
    Message: value,
   })
  }
  return false, errs
 }
 return true, nil
}

```

</details>

下面是使用方式

<details>
<summary>制定对象规则</summary>

```golang
type UpdateTagRequest struct {
 Id uint32 `form:"id" binding:"required,gte=1"`
 Name   string `form:"name" binding:"required,min=3,max=100"`
 ModifiedBy string `form:"modified_by" binding:"required,min=3,max=100"`
 State  uint8  `form:"state,default=1" binding:"oneof=0 1"`
}
```

</details>

<details>
<summary>具体的方法调用</summary>

```golang
func (a Tag) Update(c *gin.Context) {
 params := service.UpdateTagRequest{Id: convert.StrTo(c.Param("id")).MustUInt32()}
 response := app.NewResponse(c)
 valid, errs := app.BindAndValid(c, &params)
 if !valid {
  // global.Logger.Errorf("app bindandvalid err: %v", errs)
  response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
  return
 }
 ...
}
```

</details>

## jwt

```shell
go get -u github.com/dgrijalva/jwt-go 
```

<details>
<summary>生成一个jwt token</summary>

```golang
func GenerateToken(appKey, appSecret string) (string, error) {
 nowTime := time.Now()
 expireTime := nowTime.Add(global.JWTSetting.Expire)
 claims := Claims{
  AppKey:util.EncodeMD5(appKey),
  AppSecret: util.EncodeMD5(appSecret),
  StandardClaims: jwt.StandardClaims{
   ExpiresAt: expireTime.Unix(),
   Issuer:global.JWTSetting.Issuer,
  },
 }

 tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
 token, err := tokenClaims.SignedString([]byte(GetJWTSecret()))
 return token, err
}
```

</details>

<details>
<summary>解析一个jwt token</summary>

```golang
func ParseToken(token string) (*Claims, error) {
 tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
  return []byte(GetJWTSecret()), nil
 })
 if tokenClaims != nil {
  if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
   return claims, nil
  }
 }
 return nil, err
}

```

</details>

<details>
<summary>一个jwt中间件</summary>

```golang
func JWT() gin.HandlerFunc {
 return func(c *gin.Context) {
  var (
   token string
   ecode = errcode.Success
  )
  if s, exist := c.GetQuery("token"); exist {
   token = s
  } else {
   token = c.GetHeader("token")
  }
  if token == "" {
   ecode = errcode.InvalidParams
  } else {
   _, err := app.ParseToken(token)
   if err != nil {
    switch err.(*jwt.ValidationError).Errors {
    case jwt.ValidationErrorExpired:
     ecode = errcode.UnauthorizedTokenTimeout
    default:
     ecode = errcode.UnauthorizedTokenError
    }
   }
  }
  if ecode != errcode.Success {
   response := app.NewResponse(c)
   response.ToErrorResponse(ecode)
   c.Abort()
   return
  }
  c.Next()
 }
}
```

</details>

## ratelimit 限流

```shell
go get -u github.com/juju/ratelimit 
```

<details>
<summary>封装一个基本的限流器</summary>

```golang
type LimiterIface interface {
 Key(c *gin.Context) string
 GetBucket(key string) (*ratelimit.Bucket, bool)
 AddBuckets(rules ...LimiterBucketRule) LimiterIface
}

type LimiterBucketRule struct {
 Key  string
 FillInterval time.Duration
 Capacity int64
 Quantum  int64
}
type Limiter struct {
 LimiterBuckets map[string]*ratelimit.Bucket
}
```

</details>

<details>
<summary>再次封装成一个对特定路由的限流器</summary>

```golang
type MethodLimiter struct {
 *Limiter
}

func NewMethodLimiter() LimiterIface {
 return MethodLimiter{
  Limiter: &Limiter{LimiterBuckets: make(map[string]*ratelimit.Bucket)},
 }
}

func (l MethodLimiter) Key(c *gin.Context) string {
 uri := c.Request.RequestURI
 index := strings.Index(uri, "?")
 if index == -1 {
  return uri
 }
 return uri[:index]
}

func (l MethodLimiter) GetBucket(key string) (*ratelimit.Bucket, bool) {
 bucket, ok := l.LimiterBuckets[key]
 return bucket, ok
}

func (l MethodLimiter) AddBuckets(rules ...LimiterBucketRule) LimiterIface {
 for _, rule := range rules {
  if _, ok := l.LimiterBuckets[rule.Key]; !ok {
   l.LimiterBuckets[rule.Key] = ratelimit.NewBucketWithQuantum(rule.FillInterval, rule.Capacity, rule.Quantum)
  }
 }
 return l
}
```

</details>

<details>
<summary>gin的中间件</summary>

```golang
var methodLimiters = limiter.NewMethodLimiter().AddBuckets(limiter.LimiterBucketRule{
 Key:  "/auth",
 FillInterval: time.Second,
 Capacity: 10,
 Quantum:  10,
})

r.Use(middleware.RateLimiter(methodLimiters))

func RateLimiter(l limiter.LimiterIface) gin.HandlerFunc {
 return func(c *gin.Context) {
  key := l.Key(c)
  if bucket, ok := l.GetBucket(key); ok {
   count := bucket.TakeAvailable(1)
   if count == 0 {
    response := app.NewResponse(c)
    response.ToErrorResponse(errcode.TooManyRequests)
    c.Abort()
    return
   }
  }
  c.Next()
 }
}
```

</details>

## opentracing and jaeger

```shell
go get -u github.com/opentracing/opentracing-go
go get -u github.com/uber/jaeger-client-go
go get -u github.com/eddycjy/opentracing-gorm  //gorm的trace在gorm那里讲过了就一行code
```

jaeger安装就不写了，直接docker

<details>
<summary>写一个tracer，opentracing.SetGlobalTracer(tracer)
这句我当时忘写了，直到做grpc那章追踪连不起来才发现这个错误</summary>

```golang
func NewJaegerTracer(servicename, agentHostPort string) (opentracing.Tracer, io.Closer, error) {
 cfg := &config.Configuration{
  ServiceName: servicename,
  Sampler: &config.SamplerConfig{
   Type:  "const",
   Param: 1,
  },
  Reporter: &config.ReporterConfig{
   LogSpans:true,
   BufferFlushInterval: 1 * time.Second,
   LocalAgentHostPort:  agentHostPort,
  },
 }
 tracer, closer, err := cfg.NewTracer()
 if err != nil {
  return nil, nil, err
 }
 opentracing.SetGlobalTracer(tracer)
 return tracer, closer, nil
}

// 在main.go里面初始化它
func setupTracing() error {
 tacer, _, err := tracer.NewJaegerTracer("blog_service", "127.0.0.1:6831")
 if err != nil {
  return err
 }
 global.Tracer = tacer
 return nil
}
```

</details>

<details>
<summary>写一个中间件，这段开始也写错，也是到grpc才发现的</summary>

```golang
func Tracing() gin.HandlerFunc {
 return func(c *gin.Context) {

  var newCtx context.Context
  var span opentracing.Span
  spanCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
  if err != nil {
   span, newCtx = opentracing.StartSpanFromContextWithTracer(c.Request.Context(), global.Tracer, c.Request.URL.Path)
  } else {
   span, newCtx = opentracing.StartSpanFromContextWithTracer(
    c.Request.Context(),
    global.Tracer,
    c.Request.URL.Path,
    opentracing.ChildOf(spanCtx),
    opentracing.Tag{Key: string(ext.Component), Value: "HTTP"},
   )
  }
  defer span.Finish()

  var tracid string
  var spanid string
  var spanContext = span.Context()
  switch spanContext.(type) {
  case jaeger.SpanContext:
   tracid = spanContext.(jaeger.SpanContext).TraceID().String()
   spanid = spanContext.(jaeger.SpanContext).SpanID().String()
  }
  c.Set("X-Trace-ID", tracid)
  c.Set("X-Span-ID", spanid)
  c.Request = c.Request.WithContext(newCtx)
  c.Next()
 }
}
```

</details>

<details>
<summary>对原有的log做改进，要记录下X-Trace-ID X-Span-ID</summary>

```golang
func (l *Logger) WithTrace() *Logger {
 ginCtx, ok := l.ctx.(*gin.Context)
 if ok {
  return l.WithFields(Fields{
   "trace_id": ginCtx.MustGet("X-Trace-ID"),
   "span_id":  ginCtx.MustGet("X-Span-ID"),
  })
 }
 return l
} 

// 原有的所有方法都要再加上.WithTrace()
func (l *Logger) Panicf(ctx context.Context, format string, v ...interface{}) {
 l.WithLevel(LevelPanic).WithContext(ctx).WithTrace().Output(fmt.Sprintf(format, v...))
}
```

</details>
