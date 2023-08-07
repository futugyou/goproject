package main

import (
	"fmt"
	"os"
)

const (
	EOF = "\r\n"
)

type ActionSet struct {
	Delimeter string
}

func NewDefaultActionSet() *ActionSet {
	var delimeter = "_GitHubActionsFileCommandDelimeter_"
	return &ActionSet{
		Delimeter: "%s<<" + delimeter + EOF + "%s" + EOF + delimeter,
	}
}

func NewActionSet(delimeter string) *ActionSet {
	return &ActionSet{
		Delimeter: "%s<<" + delimeter + EOF + "%s" + EOF + delimeter,
	}
}

func (s *ActionSet) SetOutput(key, value string) {
	s.setGitFile(key, value, "GITHUB_OUTPUT")
}

func (s *ActionSet) SetEnv(key, value string) {
	s.setGitFile(key, value, "GITHUB_ENV")
}

func (s *ActionSet) setGitFile(key, value, command string) {
	msg := []byte(fmt.Sprintf(s.Delimeter, key, value) + EOF)

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
