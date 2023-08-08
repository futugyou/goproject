FROM golang:1.18.3-alpine3.16
RUN apk add git 
WORKDIR /app

COPY ./ ./

RUN go build -o /bin/app .

ENTRYPOINT ["app"]