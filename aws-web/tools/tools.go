package tools

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"net"
	"net/url"
	"strconv"

	"github.com/futugyousuzu/goproject/awsgolang/core"
)

func IP4ToLong(ip string) uint32 {
	var long uint32
	binary.Read(bytes.NewBuffer(net.ParseIP(ip).To4()), binary.BigEndian, &long)
	return long
}

func LongToIP4(ipInt int64) string {
	// need to do two bit shifting and “0xff” masking
	b0 := strconv.FormatInt((ipInt>>24)&0xff, 10)
	b1 := strconv.FormatInt((ipInt>>16)&0xff, 10)
	b2 := strconv.FormatInt((ipInt>>8)&0xff, 10)
	b3 := strconv.FormatInt((ipInt & 0xff), 10)
	return b0 + "." + b1 + "." + b2 + "." + b3
}

func Sha1(text string) string {
	h := sha1.New()
	h.Write([]byte(text))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func String(n int32) string {
	buf := [11]byte{}
	pos := len(buf)
	i := int64(n)
	signed := i < 0
	if signed {
		i = -i
	}
	for {
		pos--
		buf[pos], i = '0'+byte(i%10), i/10
		if i == 0 {
			if signed {
				pos--
				buf[pos] = '-'
			}
			return string(buf[pos:])
		}
	}
}

func GetPaging(url url.Values) core.Paging {
	pageString := url.Get("page")
	limitString := url.Get("limit")

	page, _ := strconv.ParseInt(pageString, 10, 64)
	limit, _ := strconv.ParseInt(limitString, 10, 64)
	if page == 0 {
		page = 1
	}

	if limit == 0 {
		limit = 10
	}

	return core.Paging{Page: page, Limit: limit}
}
