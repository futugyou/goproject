package code0392

import "fmt"

func Exection() {
	s := "abe"
	t := "ahbgdc"
	exection(s, t)
}

func exection(s, t string) {
	i := 0
	j := 0
	for {
		if s[i] == t[j] {
			i++
		}
		j++
		if i >= len(s) || j >= len(t) {
			break
		}
	}
	fmt.Println(len(s) == i)
}
