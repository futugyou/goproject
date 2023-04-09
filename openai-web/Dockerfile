FROM golang:1.20.0 AS build
ARG TARGETOS
ARG TARGETARCH

WORKDIR /app
COPY ./openai ./openai
COPY ./openai-tokenizer ./openai-tokenizer
COPY ./openai-web/go.mod ./openai-web/go.mod
COPY ./openai-web/go.sum ./openai-web/go.sum
RUN cd ./openai-web && go mod download


COPY ./openai-web   ./openai-web
RUN cd ./openai-web && CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -o /openai ./main.go

FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=build /openai /openai
COPY ./openai-web/conf/app.conf conf/app.conf
COPY ./openai-web/examples examples
COPY ./openai-web/swagger swagger
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/openai"]
