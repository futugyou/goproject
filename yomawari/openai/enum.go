package openai

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
type ChatRole string

const ChatRoleSystem ChatRole = "system"
const ChatRoleUser ChatRole = "user"
const ChatRoleAssistant ChatRole = "assistant"

var SupportedChatRoles = []ChatRole{ChatRoleSystem, ChatRoleUser, ChatRoleAssistant}
type ImageFormatType string

const Url ImageFormatType = "url"
const B64Json ImageFormatType = "b64_json"

var SupportedImageResponseFormat = []ImageFormatType{
	Url,
	B64Json,
}
type ImageSize string

const Size256 ImageSize = "256x256"
const Size512 ImageSize = "512x512"
const Size1024 ImageSize = "1024x1024"

var SupportededImageSize = []ImageSize{
	Size256,
	Size512,
	Size1024,
}
