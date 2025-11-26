FROM golang:latest AS build
ARG TARGETOS
ARG TARGETARCH

WORKDIR /app
COPY ./openai-web/go.mod ./openai-web/go.mod
COPY ./openai-web/go.sum ./openai-web/go.sum
RUN cd ./openai-web && go mod edit -replace github.com/futugyou/extensions=github.com/futugyou/goproject/extensions@master && go mod tidy && go mod download

COPY ./openai-web   ./openai-web
RUN cd ./openai-web && CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -o /openai ./main.go

FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=build /openai /openai
COPY ./openai-web/conf/app.conf conf/app.conf
COPY ./openai-web/.env .env
COPY ./openai-web/examples examples
COPY ./openai-web/swagger swagger
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/openai"]
