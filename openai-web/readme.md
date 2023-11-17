# install bee tool

```golang
go install github.com/beego/bee/v2@develop
```

init web api project

```golang
bee api openai
```

swagger doc

```golang
bee generate docs
```

route

```golang
bee generate routers
```

```golang
bee run 
```

use docker command in parent folder

```docker
docker build -t openai -f ./openai-web/Dockerfile .
docker run  --name openai -p 8080:8080 -d openai

docker run -d --name redis-stack -p 6379:6379 -p 8001:8001 redis/redis-stack:latest
docker exec -it redis-stack redis-cli

```
