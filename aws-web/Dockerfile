FROM golang:latest

WORKDIR /app

COPY ./ ./

RUN go build -o /bin/app .

ENTRYPOINT ["app"]