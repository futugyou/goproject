package code0028

import (
	"fmt"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	haystack := "mississippi"
	needle := "issip"
	r := strStr(haystack, needle)
	fmt.Println(r)
	r = kmp(haystack, needle)
	fmt.Println(r)
}

func strStr(haystack string, needle string) int {
	return common.RabinKarp(haystack, needle)
}

func kmp(haystack string, needle string) int {
	k := common.NewKmp(needle)
	return k.Search(haystack)
}
