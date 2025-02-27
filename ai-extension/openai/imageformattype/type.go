package types

type ImageFormatType string

const Url ImageFormatType = "url"
const B64Json ImageFormatType = "b64_json"

var SupportedImageResponseFormat = []ImageFormatType{
	Url,
	B64Json,
}
