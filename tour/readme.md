# Using cobra to implement a simple command line tool

## dynamodb command

```golang
go run . dynamo generate [flags]
```

## mongodb command

```golang
go run . mongo generate [flags]
```

## openapi command

```golang
go run . openapi generate [flags]
go run . openapi swag2openapi [flags]
```

## mysql command

```golang
go run . sql struct [flags]
```

## time command (test)

```golang
go run . time now
go run . time calc [flags]
```

## word command (test)

```golang
go run . word
```

## use in go:generate with env

```golang
//go:generate go install github.com/joho/godotenv/cmd/godotenv@latest
//go:generate godotenv -f ./.env go run ../tour/main.go word
```
