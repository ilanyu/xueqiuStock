package main

import (
	"encoding/json"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

func Stock() {
	klineDatas := getKlineDatas()
	log.Println("已获取到全部股票")
	saveKlineDatas(klineDatas)
	log.Println("已将全部股票存储到数据库中")
}

func saveKlineDatas(klineDatas []KlineData) {
	session, err := mgo.Dial(cmd.MongoDBUrl)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(cmd.MongoDBName).C("stock_" + cmd.StockBeginTime + "_" + cmd.StockPeriod)
	index := mgo.Index{
		Key:        []string{"symbol", "timestamp"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err = c.EnsureIndex(index)
	if err != nil {
		log.Println(err)
	}
	for e := range klineDatas {
		ms := bson.M{"symbol": klineDatas[e].Symbol}
		for i := 0; i < len(klineDatas[e].Item); i++ {
			for j := 0; j < len(klineDatas[e].Column); j++ {
				ms[klineDatas[e].Column[j]] = klineDatas[e].Item[i][j]
			}
			n, err := c.Find(bson.M{"symbol": klineDatas[e].Symbol, "timestamp": ms["timestamp"]}).Count()
			if err != nil {
				log.Println(err)
				continue
			}
			if n > 0 {
				err := c.Update(bson.M{"symbol": klineDatas[e].Symbol, "timestamp": ms["timestamp"]}, ms)
				if err != nil {
					log.Println(err)
				}
			} else {
				err := c.Insert(ms)
				if err != nil {
					log.Println(err)
				}
			}
		}
		ms = nil
		klineDatas[e] = KlineData{}
	}
}

func getKlineDatas() []KlineData {
	stockList := GetStockList()
	cookies := GetCookie()
	lock := make(chan int, cmd.GetStockThreadNum)
	for i := 0; i < cmd.GetStockThreadNum; i++ {
		lock <- i
	}
	ch := make(chan []byte, len(stockList)*2)
	for i := 0; i < len(stockList); i++ {
		go GetContents("https://stock.xueqiu.com/v5/stock/chart/kline.json?symbol="+stockList[i].Symbol+"&begin="+cmd.StockBeginTime+"&period="+cmd.StockPeriod+"&type="+cmd.StockType+"&count="+cmd.StockCount+"&indicator=kline,ma,macd,kdj,boll,rsi,wr,bias,cci,psy", cookies, ch, lock)
	}
	var klineDatas []KlineData
	for i := 0; i < len(stockList); i++ {
		var klineDataResp KlineDataResp
		err := json.Unmarshal(<-ch, &klineDataResp)
		if err != nil {
			log.Println(err)
		}
		klineDatas = append(klineDatas, klineDataResp.Data)
	}
	stockList = nil
	return klineDatas
}

func GetStockList() []StockT {
	session, err := mgo.Dial(cmd.MongoDBUrl)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(cmd.MongoDBName).C("stock_list_" + cmd.StockBeginTime + "_" + cmd.StockPeriod)
	var stocks []StockT
	err = c.Find(bson.M{}).All(&stocks)
	if err != nil {
		panic(err)
	}
	return stocks
}
