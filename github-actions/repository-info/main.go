package main

import (
	"encoding/json"
	"flag"
	"os"
	"time"
)

var (
	reponame = flag.String("reponame", "", "repository name, if empty, show all repository.")
	token    = os.Getenv("token")
)

func main() {
	flag.Parse()

	client := NewGithubClient(token)
	repoSvc := NewRepositoryService(client.Client)
	repos := repoSvc.GetRepository(*reponame)

	set := NewDefaultActionSet()
	set.SetOutput("time", time.Now().Format("2006-01-02 15:04:05"))

	jsonString, _ := json.Marshal(repos)
	set.SetEnv("result", string(jsonString))
}
