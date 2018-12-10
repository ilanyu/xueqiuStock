package main

import (
	"flag"
	"strconv"
	"time"
)

type Cmd struct {
	MongoDBUrl            string
	MongoDBName           string
	GetStockListThreadNum int
	GetStockThreadNum     int
	StockBeginTime        string
	StockPeriod           string
	StockType             string
	StockCount            string
}

func pauseCmd() Cmd {
	var cmd Cmd
	flag.StringVar(&cmd.MongoDBUrl, "MongoDBUrl", "127.0.0.1", "MongoDB connect url")
	flag.StringVar(&cmd.MongoDBName, "MongoDBName", "xueqiuStock", "MongoDB database name")
	flag.IntVar(&cmd.GetStockListThreadNum, "GetStockListThreadNum", 10, "Get StockList Thread Num")
	flag.IntVar(&cmd.GetStockThreadNum, "GetStockThreadNum", 25, "Get Stock Thread Num")
	flag.StringVar(&cmd.StockBeginTime, "StockBeginTime", strconv.Itoa(int(time.Now().UnixNano()/1000000)), "Stock Begin Time")
	flag.StringVar(&cmd.StockPeriod, "StockPeriod", "day", "Stock Period")
	flag.StringVar(&cmd.StockType, "StockType", "before", "Stock Type")
	flag.StringVar(&cmd.StockCount, "StockCount", "-200", "Stock Count")
	flag.Parse()
	return cmd
}
