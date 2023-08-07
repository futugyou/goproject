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
	dest_folder = flag.String("dest_folder", "", "destination repository folder.")

	source_token = os.Getenv("source_token")
	dest_token   = os.Getenv("dest_token")
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

	_ = &CloneInfo{
		SourceOwner:  *source_owner,
		SourceBranch: *source_branch,
		SourceName:   *source_name,
		SourceToken:  source_token,
		DestOwner:    *dest_owner,
		DestBranch:   *dest_branch,
		DestName:     *dest_name,
		DestToken:    dest_token,
		DestFolder:   *dest_folder,
	}
}
