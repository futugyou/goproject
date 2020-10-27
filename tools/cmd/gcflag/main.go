package main

type User struct {
	ID     int64
	Name   string
	Avatar string
}

func GetUserInfo() *User {
	return &User{
		ID:     12308877,
		Name:   "ekjhhjk",
		Avatar: "https://github.com/microsoft/Terminal",
	}
}

// go build -gcflags '-m -l' .\cmd\gcflag\main.go
// go tool compile -S .\cmd\gcflag\main.go

func main() {
	_ = GetUserInfo()
}
