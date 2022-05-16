package code0990

import (
	"fmt"
	"strings"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	equations := []string{"a==b", "b==a"}
	r := exection(equations)
	fmt.Println(r)
}

func exection(equations []string) bool {
	uf := common.NewUnionFind(26)
	for _, v := range equations {
		if !strings.Contains(v, "!") {
			a := int(v[0] - 'a')
			b := int(v[3] - 'a')
			uf.Union(a, b)
		}
	}
	for _, v := range equations {
		if strings.Contains(v, "!") {
			a := int(v[0] - 'a')
			b := int(v[3] - 'a')
			if uf.Connected(a, b) {
				return false
			}
		}
	}
	return true
}
