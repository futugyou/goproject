```shell
go mod edit -replace github.com/futugyou/alphavantage=../alphavantage

go mod edit -replace github.com/futugyou/alphavantage=github.com/futugyou/goproject/alphavantage@master
go mod tidy
```

### use replace in go.work or go.mod

### this service will called by github action at UTC 23:30 every day, it will get data from www.alphavantage.co.

### Because I am lazy, I use vercel to obtain the data. See 'api' folder

## doc

![1](./doc/images/arch.drawio.png)
