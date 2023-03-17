package tokenizer

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/futugyousuzu/openai-tokenizer/common"

	"embed"

	"golang.org/x/exp/slices"
)

var encoder map[string]int = make(map[string]int)
var decoder map[int]string = make(map[int]string)
var vocab map[vacabItem]int = make(map[vacabItem]int)
var bpeCache map[string]string = make(map[string]string)
var byte_encoder map[int]string = make(map[int]string)
var byte_decoder map[string]int = make(map[string]int)

const encodingRegex string = `'s|'t|'re|'ve|'m|'ll|'d| ?\p{L}+| ?\p{N}+| ?[^\s\p{L}\p{N}]+|\s+(?!\S)|\s+`

type vacabItem struct {
	x string
	y string
}

//go:embed encoder.json
var encoderfs embed.FS

//go:embed vocab.bpe
var vocabfs embed.FS

func init() {
	initiByteEncoderAndDecoder()
	createEncoderAndDecoder()
	createVocab()
}

func Test() (map[int]string, map[string]int) {
	return byte_encoder, byte_decoder
}

func createEncoderAndDecoder() {
	en, err := encoderfs.ReadFile("encoder.json")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	json.Unmarshal(en, &encoder)
	for k, v := range encoder {
		decoder[v] = k
	}
}

func createVocab() {
	en, err := vocabfs.ReadFile("vocab.bpe")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	raws := strings.Split(string(en), "\n")[1:]
	raws = common.ArrayFilter(raws, vocabFilter)

	for i := 0; i < len(raws); i++ {
		t := strings.Split(raws[i], " ")
		item := vacabItem{x: t[0], y: t[1]}
		vocab[item] = i
	}
}

var vocabFilter = func(s string) bool {
	return len(strings.TrimSpace(s)) > 0
}

func Encode(text string) []int {
	if len(text) == 0 {
		return []int{}
	}

	matches := common.Regexp2FindAllString(encodingRegex, text)
	bpeTokens := make([]int, 0)
	for i := 0; i < len(matches); i++ {
		match := matches[i]
		tokenBytes := []byte(match)
		ts := make([]string, 0)
		for j := 0; j < len(tokenBytes); j++ {
			ts = append(ts, byte_encoder[int(tokenBytes[j])])
		}

		token := strings.Join(ts, "")
		for _, v := range strings.Split(bytePairEncoding(token), " ") {
			bpeTokens = append(bpeTokens, encoder[v])
		}
	}

	return bpeTokens
}

func Decode(tokens []int) string {
	t := make([]string, 0)
	for i := 0; i < len(tokens); i++ {
		t = append(t, decoder[tokens[i]])
	}

	text := strings.Join(t, "")
	textSplit := strings.Split(text, "")
	byteArr := make([]byte, 0)
	for i := 0; i < len(textSplit); i++ {
		byteArr = append(byteArr, byte(byte_decoder[textSplit[i]]))
	}

	return string(byteArr)
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

func bytePairEncoding(token string) string {
	if v, ok := bpeCache[token]; ok {
		return v
	}

	word := strings.Split(token, "")
	pairs := getPairs(word)
	if len(pairs) == 0 {
		bpeCache[token] = token
		return token
	}

	for {
		minPairs := make(map[int]vacabItem)
		for key, pair := range pairs {
			for _, p := range pair {
				v := vacabItem{x: key, y: p}
				if rank, ok := vocab[v]; ok {
					minPairs[rank] = v
				} else {
					minPairs[100000000000] = v
				}
			}
		}

		mapkey := common.GetMapKeys(minPairs)
		sort.Ints(mapkey)
		minKey := mapkey[0]
		var biGram = minPairs[minKey]

		if _, ok := vocab[biGram]; !ok {
			break
		}

		first := biGram.x
		second := biGram.y

		newWord := make([]string, 0)
		i := 0

		for {
			if i >= len(word) {
				break
			}
			j := common.IndexArrayWithOffset(word, first, i)

			if j == -1 {
				newWord = append(newWord, word[i:]...)
				break
			}

			newWord = append(newWord, word[i:j]...)
			i = j

			if word[i] == first && i < len(word)-1 && word[i+1] == second {
				newWord = append(newWord, first+second)
				i += 2
			} else {
				newWord = append(newWord, word[i])
				i += 1
			}
		}

		word = newWord

		if len(word) == 1 {
			break
		}

		pairs = getPairs(word)
	}

	result := strings.Join(word, " ")
	bpeCache[token] = result
	return result
}

func initiByteEncoderAndDecoder() {
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

	n := 0
	for i := 0; i < 256; i++ {
		if !slices.Contains(list, rune(i)) {
			list = append(list, rune(i))
			list2 = append(list2, rune(256+n))
			n = n + 1
		}
	}

	for i := 0; i < len(list2); i++ {
		byte_encoder[int(list[i])] = string(rune(list2[i]))
		byte_decoder[string(rune(list2[i]))] = int(list[i])
	}
}
