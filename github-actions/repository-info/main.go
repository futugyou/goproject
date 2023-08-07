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

	client := NewGithubClient(token)
	repoSvc := NewRepositoryService(client.Client)
	repos := repoSvc.GetRepository(*reponame)

	SetOutput("time", time.Now().Format("2006-01-02 15:04:05"))

	jsonString, _ := json.Marshal(repos)
	SetEnv("result", string(jsonString))
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
