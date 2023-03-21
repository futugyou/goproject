FROM golang:1.20.0

WORKDIR /app
COPY ./openai ./openai
COPY ./openai-tokenizer ./openai-tokenizer
COPY ./openai-web   ./openai-web

RUN cd ./openai-web && go mod tidy
RUN cd ./openai-web && go build -o /openai

EXPOSE 8080

CMD [ "/openai" ]