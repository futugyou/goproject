package add

var datas []string

func Add(str string) int {
	data := []byte(str)
	datas := append(datas, string(data))
	return len(datas)
}
