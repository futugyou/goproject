package main

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptrace"
	"net/http/httputil"
	"strings"
	"time"
)

func Timeout(network, host string) (net.Conn, error) {
	conn, err := net.DialTimeout(network, host, 15*time.Second)
	if err != nil {
		return nil, err
	}
	conn.SetDeadline(time.Now().Add(15 * time.Second))
	return conn, nil
}

func main() {
	client := http.Client{
		Timeout: 15 * time.Second,
		//Transport: &http.Transport{Dial: Timeout},
	}
	url := "http://www.baidu.com"
	req, _ := http.NewRequest("GET", url, nil)
	trace := &httptrace.ClientTrace{
		GotFirstResponseByte: func() {
			fmt.Println("first response byte")
		},
		GotConn: func(connInfo httptrace.GotConnInfo) {
			fmt.Printf("got conn: %+v\n", connInfo)
		},
		DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
			fmt.Printf("got dns info : %+v\n", dnsInfo)
		},
		ConnectStart: func(network, addr string) {
			fmt.Println("dial start")
		},
		ConnectDone: func(network, addr string, err error) {
			fmt.Println("dial done")
		},
		WroteHeaders: func() {
			fmt.Println("wrote headers")
		},
	}

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	_, err := http.DefaultTransport.RoundTrip(req)

	if err != nil {
		fmt.Println(err)
		return
	}

	response, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("status code:", response.Status)
	header, _ := httputil.DumpResponse(response, false)
	fmt.Println(string(header))

	contentType := response.Header.Get("Content-Type")
	charset := strings.SplitAfter(contentType, "charset=")
	if len(charset) > 1 {
		fmt.Println("charset set:", charset[1])
	}
	if response.ContentLength == -1 {
		fmt.Println("contentlength is unknown")
	} else {
		fmt.Println("contentlength :", response.ContentLength)
	}
	length := 0
	var buffer [1024]byte
	r := response.Body
	for {
		n, err := r.Read(buffer[0:])
		if err != nil {
			fmt.Println(err)
			break
		}
		length += n
	}
	fmt.Println("cale response data length:", length)
}
