package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"
)

var (
	reponame = flag.String("reponame", "", "repository name, if empty, show all repository.")
	token    = os.Getenv("token")
)

const (
	EOF                = "\r\n"
	multiLineFileDelim = "_GitHubActionsFileCommandDelimeter_"
	multilineFileCmd   = "%s<<" + multiLineFileDelim + EOF + "%s" + EOF + multiLineFileDelim
)

func main() {
	flag.Parse()

	time := time.Now()
	SetOutput("time", time.Format("2006-01-02 15:04:05"))
	SetOutput("branch", "master")

	options := make([]Option, 0)
	options = append(options, Option{Own: "tom", Branch: "dev", Url: "http://baidu"}, Option{Own: "tony", Branch: "master", Url: "http://google"})

	jsonString, _ := json.Marshal(options)
	SetEnv("option", string(jsonString))
	SetEnv("doNextStep", "ok")
}

type Option struct {
	Own    string `json:"own"`
	Branch string `json:"branch"`
	Url    string `json:"url"`
}

func SetOutput(key, value string) {
	setGitFile(key, value, "GITHUB_OUTPUT")
}

func SetEnv(key, value string) {
	setGitFile(key, value, "GITHUB_ENV")
}

func setGitFile(key, value, command string) {
	msg := []byte(fmt.Sprintf(multilineFileCmd, key, value) + EOF)

	filepath := os.Getenv(command)
	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	if _, err := f.Write(msg); err != nil {
		fmt.Println(err)
		return
	}
}
