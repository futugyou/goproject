install bee tool
```
go install github.com/beego/bee/v2@develop
```

init web api project
```
bee api openai
```

swagger doc
```
bee generate docs
```

route
```
bee generate routers
```

```
bee run 
```

use docker command in parent folder
```
docker build -t openai -f ./openai-web/Dockerfile .
docker run  --name openai -p 8080:8080 -d openai

docker run -d --name redis-stack -p 6379:6379 -p 8001:8001 redis/redis-stack:latest

```