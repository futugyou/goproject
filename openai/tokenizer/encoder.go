package tokenizer

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"golang.org/x/exp/slices"
)

var encoder map[string]int
var vocab map[vacabItem]int = make(map[vacabItem]int)
var bpeCache map[string]string = make(map[string]string)
var bytesToUnicodeCache map[int]string = make(map[int]string)

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
		item := vacabItem{x: t[0], y: t[1]}
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

func getPairs(word []string) (result map[string][]string) {
	result = make(map[string][]string)
	prevChar := word[0]
	for i := 1; i < len(word); i++ {
		curr := word[i]
		if _, ok := result[prevChar]; !ok {
			result[prevChar] = make([]string, 0)
		}
		result[prevChar] = append(result[prevChar], curr)
		prevChar = curr
	}
	return
}

func indexArrayWithOffset(array []string, first string, offset int) int {
	if len(array) < offset {
		return -1
	}
	arr := array[offset:]
	for i := 0; i < len(arr); i++ {
		if arr[i] == first {
			return i + offset
		}
	}
	return -1
}

func indexWithOffset(s, substr string, offset int) int {
	if len(s) < offset {
		return -1
	}
	if idx := strings.Index(s[offset:], substr); idx >= 0 {
		return offset + idx
	}
	return -1
}

func getMapKeys(m map[int]vacabItem) []int {
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func initializeBytesToUnicodeCache() {
	result := make(map[int]string)
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
		result[int(list[i])] = string(rune(list2[i]))
	}

	bytesToUnicodeCache = result
}
