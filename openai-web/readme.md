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
```