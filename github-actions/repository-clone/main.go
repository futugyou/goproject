package main

import (
	"flag"
	"log"
	"os"
)

var (
	source_owner  = flag.String("owner", "", "source repository owner.")
	source_branch = flag.String("branch", "", "source repository branch.")
	source_name   = flag.String("name", "", "source repository name.")

	dest_owner  = flag.String("dest_owner", "", "destination repository owner.")
	dest_branch = flag.String("dest_branch", "", "destination repository branch.")
	dest_name   = flag.String("dest_name", "", "destination repository name.")

	source_token = os.Getenv("source_token")
	dest_token   = os.Getenv("dest_token")

	step = flag.String("step", "1", "step.")
)

func main() {
	flag.Parse()

	if len(*source_owner) == 0 ||
		len(*source_branch) == 0 ||
		len(*source_name) == 0 ||
		len(*dest_owner) == 0 ||
		len(*dest_branch) == 0 ||
		len(*dest_name) == 0 {
		log.Println("not enough parameters.")
		return
	}

	info := &CloneInfo{
		SourceOwner:  *source_owner,
		SourceBranch: *source_branch,
		SourceName:   *source_name,
		SourceToken:  source_token,
		DestOwner:    *dest_owner,
		DestBranch:   *dest_branch,
		DestName:     *dest_name,
		DestToken:    dest_token,
	}

	CloneRepoPipeline(info)
}

func CloneRepoPipeline(info *CloneInfo) {
	git := NewGitCommand(info)
	git.SetConfig()

	if *step == "1" {
		CreateRepository(info)
		if err := git.CloneDest(); err != nil {
			return
		}

		if err := git.RmoveDest(); err != nil {
			return
		}

		git.CloneSource()
	} else {
		git.GitAdd()
		needCommit := git.GitStatus()
		if needCommit {
			if err := git.GitCommit(); err != nil {
				return
			}

			git.GitPush()
		}
	}
}

func CreateRepository(info *CloneInfo) error {
	client := NewGithubClient(info.DestToken)
	repoSvc := NewRepositoryService(client.Client)
	return repoSvc.CreateRepository(info)
}
