package code0204

import "fmt"

func Exection() {
	n := 10
	exection(n)
}

func exection(n int) {
	count := 0
	isprime := make([]bool, n)
	for i := 0; i < n; i++ {
		isprime[i] = true
	}
	for i := 2; i*i < n; i++ {
		if isprime[i] {
			for j := i * i; j < n; j = j + i {
				isprime[j] = false
			}
		}
	}
	for i := 2; i < n; i++ {
		if isprime[i] {
			count++
		}
	}
	fmt.Println(count)
}
