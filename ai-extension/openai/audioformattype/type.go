package types

type AudioFormatType string

const Json AudioFormatType = "json"
const Text AudioFormatType = "text"
const Srt AudioFormatType = "srt"
const VerboseJson AudioFormatType = "verbose_json"
const Vtt AudioFormatType = "vtt"

var SupportededResponseFormatType = []AudioFormatType{
	Json,
	Text,
	Srt,
	VerboseJson,
	Vtt,
}
