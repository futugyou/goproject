package imagesize

type ImageSize string

const Size256 ImageSize = "256x256"
const Size512 ImageSize = "512x512"
const Size1024 ImageSize = "1024x1024"

var SupportededImageSize = []ImageSize{
	Size256,
	Size512,
	Size1024,
}
