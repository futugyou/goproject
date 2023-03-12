package tokenizer

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

var encoder map[string]int
var vocab map[*vacabItem]int = make(map[*vacabItem]int)

type vacabItem struct {
	x string
	y string
}

func init() {
	createEncoder()
	createVocab()
}

func createEncoder() {
	en, err := os.ReadFile("./tokenizer/encoder.json")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	json.Unmarshal(en, &encoder)
}

func createVocab() {
	en, err := os.ReadFile("./tokenizer/vocab.bpe")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	raws := strings.Split(string(en), "\n")[1:]
	raws = filter(raws, vocabFilter)

	for i := 0; i < len(raws); i++ {
		t := strings.Split(raws[i], " ")
		item := &vacabItem{x: t[0], y: t[1]}
		vocab[item] = i
	}
}

var vocabFilter = func(s string) bool {
	return len(strings.TrimSpace(s)) > 0
}

func filter(raws []string, filter func(string) bool) (ret []string) {
	for i := 0; i < len(raws); i++ {
		if filter(raws[i]) {
			ret = append(ret, raws[i])
		}
	}

	return
}
