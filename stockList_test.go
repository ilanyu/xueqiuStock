package main

import (
	"encoding/json"
	"testing"
	"time"
)

func TestGetCookie(t *testing.T) {
	cookies := GetCookie()
	for _, k := range cookies {
		t.Logf("%+v\n", k)
	}
}

func TestName(t *testing.T) {
	print(time.Now().UnixNano() / 1000000)
}

func TestGetContents(t *testing.T) {
	cookies := GetCookie()
	ch := make(chan []byte, 10)
	lock := make(chan int, 10)
	lock <- 1
	go GetContents("https://xueqiu.com/stock/cata/stocklist.json?page=1&size=100&order=desc&orderby=percent&type=11%2C12&_=1544168052880", cookies, ch, lock)
	stockList := &StockListT{}
	err := json.Unmarshal(<-ch, &stockList)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	if stockList.Success != "true" {
		t.Fail()
	}
}
