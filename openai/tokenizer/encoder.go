package tokenizer

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"golang.org/x/exp/slices"
)

var encoder map[string]int
var vocab map[*vacabItem]int = make(map[*vacabItem]int)
var bpeCache map[string]string = make(map[string]string)

const encodingRegex string = `'s|'t|'re|'ve|'m|'ll|'d| ?\p{L}+| ?\p{N}+| ?[^\s\p{L}\p{N}]+|\s+(?!\S)|\s+`

type vacabItem struct {
	x string
	y string
}

func init() {
	initializeBytesToUnicodeCache()
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

func initializeBytesToUnicodeCache() []string {
	result := make([]string, 0)
	list, list2 := []rune{}, []rune{}

	for i := '!'; i < '~'+1; i++ {
		list = append(list, i)
		list2 = append(list2, i)
	}

	for i := '¡'; i < '¬'+1; i++ {
		list = append(list, i)
		list2 = append(list2, i)
	}

	for i := '®'; i < 'ÿ'+1; i++ {
		list = append(list, i)
		list2 = append(list2, i)
	}

	n := 1
	for i := 0; i < 256; i++ {
		if !slices.Contains(list, rune(i)) {
			list = append(list, rune(i))
			list2 = append(list2, rune(256+n))
			n = n + 1
		}
	}

	for i := 0; i < len(list2); i++ {
		result = append(result, string(rune(list2[i])))
	}

	return result
}
