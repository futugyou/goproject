package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

func f1(t int) {
	c1 := context.Background()
	c1, cancel := context.WithCancel(c1)
	defer cancel()
	go func() {
		time.Sleep(4 * time.Second)
		cancel()
	}()
	select {
	case <-c1.Done():
		fmt.Println("f1()", c1.Err())
		return
	case r := <-time.After(time.Duration(t) * time.Second):
		fmt.Println("f1():", r)
	}
	return
}

func f2(t int) {
	c1 := context.Background()
	c1, cancel := context.WithTimeout(c1, time.Duration(t)*time.Second)
	defer cancel()

	go func() {
		time.Sleep(4 * time.Second)
		cancel()
	}()

	select {
	case <-c1.Done():
		fmt.Println("f2():", c1.Err())
		return
	case r := <-time.After(time.Duration(t) * time.Second):
		fmt.Println("f2():", r)
	}
}

func ff3(t int) {
	c3 := context.Background()
	deadline := time.Now().Add(time.Duration(2*t) * time.Second)
	c3, cancel := context.WithDeadline(c3, deadline)
	defer cancel()

	go func() {
		time.Sleep(4 * time.Second)
		cancel()
	}()

	select {
	case <-c3.Done():
		fmt.Println("f3():", c3.Err())
		return
	case r := <-time.After(time.Duration(t) * time.Second):
		fmt.Println("f3():", r)
	}
	return
}

var (
	myurl string = "http://www.baidu.com"
	delay int    = 5
	w     sync.WaitGroup
)

type myData struct {
	r   *http.Response
	err error
}

func connect(c context.Context) error {
	defer w.Done()
	data := make(chan myData, 1)
	tr := &http.Transport{}
	httpClient := &http.Client{Transport: tr}
	req, _ := http.NewRequest("GET", myurl, nil)

	go func() {
		response, err := httpClient.Do(req)
		if err != nil {
			fmt.Println(err)
			data <- myData{nil, err}
			return
		} else {
			pack := myData{response, err}
			data <- pack
		}
	}()

	select {
	case <-c.Done():
		tr.CancelRequest(req)
		<-data
		fmt.Println("canceled")
		return c.Err()
	case ok := <-data:
		err := ok.err
		resp := ok.r
		if err != nil {
			fmt.Println("error :", err)
			return err
		}
		defer resp.Body.Close()

		relHTTPData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("error ", err)
			return err
		}
		fmt.Println("response: %s\n", string(relHTTPData))
	}
	return nil
}

func main() {
	c := context.Background()
	c, cancel := context.WithTimeout(c, time.Duration(4)*time.Second)
	defer cancel()

	fmt.Printf("connecting to %s \n", myurl)
	w.Add(1)
	go connect(c)
	w.Wait()
	fmt.Println("exiting...")
}
