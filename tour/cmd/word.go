package cmd

import (
	"github/go-project/tour/internal/word"
	"log"
	"strings"

	"github.com/spf13/cobra"
)

var str string
var mode int8
var desc = strings.Join([]string{
	"this is frist line ",
	"1 to upper",
	"2 tolower",
	"3 underline to upper camelcase",
	"4 underline to lower camelcase",
	"5 camelcase to underline",
}, "\n")

var wordCmd = &cobra.Command{
	Use:   "word",
	Short: "change word",
	Long:  desc,
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		var content string
		switch mode {
		case ModeUpper:
			content = word.ToUpper(str)
		case ModeLower:
			content = word.ToLower(str)
		case ModeUnderscoreToUpperCamelcase:
			content = word.UnderscoreToUpperCamelCase(str)
		case ModeUnderscoreToLowerCamelcase:
			content = word.UnderscoreToLowerCamelCase(str)
		case ModeCamelcaseToUnderscore:
			content = word.CamelCaseToUnderscore(str)
		default:
			log.Fatalf("not support now")
		}
		log.Printf("output is %s", content)
	},
}

const (
	// ModeUpper is upper
	ModeUpper = iota + 1
	// ModeLower is lower
	ModeLower
	// ModeUnderscoreToUpperCamelcase ModeUnderscoreToUpperCamelcase
	ModeUnderscoreToUpperCamelcase
	// ModeUnderscoreToLowerCamelcase ModeUnderscoreToLowerCamelcase
	ModeUnderscoreToLowerCamelcase
	// ModeCamelcaseToUnderscore ModeCamelcaseToUnderscore
	ModeCamelcaseToUnderscore
)

func init() {
	wordCmd.Flags().StringVarP(&str, "str", "s", "", "please input word !")
	wordCmd.Flags().Int8VarP(&mode, "mode", "m", 0, "please intout change mode !")
}
