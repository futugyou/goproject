FROM golang:latest

WORKDIR /app

COPY ./ ./
RUN go mod edit -replace github.com/futugyou/alphavantage=github.com/futugyou/goproject/alphavantage@master
RUN go mod tidy
RUN go build -o /bin/app .

ENTRYPOINT ["app"]