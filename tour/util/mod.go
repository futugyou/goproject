package util

import (
	"fmt"
	"os"

	"golang.org/x/mod/modfile"
)

func GetModuleName() string {
	goModBytes, err := os.ReadFile("go.mod")
	if err != nil {
		fmt.Println("Error reading go.mod:", err)
		return ""
	}
	return modfile.ModulePath(goModBytes)
}
