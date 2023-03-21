FROM golang:1.20.0 AS build
WORKDIR /app
COPY ./openai ./openai
COPY ./openai-tokenizer ./openai-tokenizer
COPY ./openai-web   ./openai-web
RUN cd ./openai-web && go mod tidy
RUN cd ./openai-web && go build -o /openai

FROM gcr.io/distroless/base-debian11
WORKDIR /
COPY --from=build /openai /openai
COPY ./openai-web/conf/app.conf conf/app.conf
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/openai"]
