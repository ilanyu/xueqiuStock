package main

import "log"

var cmd Cmd

func main() {
	cmd = pauseCmd()
	log.Println("正在爬取股票列表")
	StockList()
	log.Println("正在爬取具体股票")
	Stock()
}
