# tip

change .env.template to .env

Go: Install/Update Tools

## debug

```golang
go install github.com/go-delve/delve/cmd/dlv@latest
```

## markdown check

```nodejs
npm install -g markdownlint-cli
markdownlint  '**/*.md'
```

## go test

```cmd
go get -u github.com/cweill/gotests/...
gotests -all -w ./*
go install gotest.tools/gotestsum@latest
gotestsum --junitfile ./tmp/test-reports/aws-web-unit-tests.xml
```

## doc

1. [awsconfig](./doc/01.awsconfig.md)
