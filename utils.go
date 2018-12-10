package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
)

func GetCookie() []*http.Cookie {
	request, err := http.NewRequest("GET", "https://xueqiu.com/", nil)
	if err != nil {
		panic(err)
	}
	client := http.DefaultClient
	for {
		response, err := client.Do(request)
		if err != nil {
			log.Println(err)
		}
		return response.Cookies()
	}
}

func GetContents(url string, cookies []*http.Cookie, ch chan []byte, lock chan int) {
	recover()
	var id int
	id = <-lock
	defer func() {
		lock <- id
	}()
	client := http.DefaultClient
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
	}
	for _, cookie := range cookies {
		request.AddCookie(cookie)
	}
	for {
		response, err := client.Do(request)
		if err != nil {
			log.Println(err)
			continue
		}
		defer response.Body.Close()
		buf := make([]byte, 8)
		var buffer bytes.Buffer
		if _, err := io.CopyBuffer(&buffer, response.Body, buf); err != nil {
			log.Println(err)
			continue
		}
		ch <- buffer.Bytes()
		return
	}
}
