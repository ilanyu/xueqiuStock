package main

import (
	"bytes"
	"encoding/json"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"
)

func StockList() {
	stocks := getAllStock()
	log.Println("已在线获取到股票列表")
	saveStock(stocks)
	log.Println("已将股票列表保存到数据库")
}

func GetPageCount(url string, cookies []*http.Cookie) int64 {
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
		var stockList StockListT
		err = json.Unmarshal(buffer.Bytes(), &stockList)
		if err != nil {
			log.Println(err)
			continue
		}
		return int64(math.Ceil(float64(stockList.Count.Count) / 100))
	}
}

func saveStock(stocks []StockT) {
	session, err := mgo.Dial(cmd.MongoDBUrl)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(cmd.MongoDBName).C("stock_list_" + cmd.StockBeginTime + "_" + cmd.StockPeriod)
	index := mgo.Index{
		Key:        []string{"code"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err = c.EnsureIndex(index)
	if err != nil {
		log.Println(err)
	}
	for e := range stocks {
		n, err := c.Find(bson.M{"code": stocks[e].Code}).Count()
		if err != nil {
			log.Println(err)
			continue
		}
		if n > 0 {
			err := c.Update(bson.M{"code": stocks[e].Code}, stocks[e])
			if err != nil {
				log.Println(err)
			}
		} else {
			err = c.Insert(stocks[e])
			if err != nil {
				log.Println(err)
			}
		}
		stocks[e] = StockT{}
	}
}

func getAllStock() []StockT {
	cookies := GetCookie()
	lock := make(chan int, cmd.GetStockListThreadNum)
	for i := 0; i < cmd.GetStockListThreadNum; i++ {
		lock <- i
	}
	pageNumCount := GetPageCount("https://xueqiu.com/stock/cata/stocklist.json?page=1&size=100&order=desc&orderby=percent&type=11&_="+strconv.Itoa(int(time.Now().UnixNano()/1000000)), cookies)
	ch := make(chan []byte, pageNumCount*2)
	for i := 1; i <= int(pageNumCount); i++ {
		go GetContents("https://xueqiu.com/stock/cata/stocklist.json?page="+strconv.Itoa(i)+"&size=100&order=desc&orderby=percent&type=11&_="+strconv.Itoa(int(time.Now().UnixNano()/1000000)), cookies, ch, lock)
	}
	var stocks []StockT
	for i := 1; i <= int(pageNumCount); i++ {
		var stockList StockListT
		err := json.Unmarshal(<-ch, &stockList)
		if err != nil {
			log.Println(err)
		}
		for e := range stockList.Stocks {
			stocks = append(stocks, stockList.Stocks[e])
		}
	}
	return stocks
}
