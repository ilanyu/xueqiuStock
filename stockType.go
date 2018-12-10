package main

type StockListT struct {
	Count   CountT   `json:"count"`
	Success string   `json:"success"`
	Stocks  []StockT `json:"stocks"`
}

type CountT struct {
	Count int64 `json:"count"`
}

type StockT struct {
	Symbol        string `json:"symbol"`
	Code          string `json:"code"`
	Name          string `json:"name"`
	Current       string `json:"current"`
	Percent       string `json:"percent"`
	Change        string `json:"change"`
	High          string `json:"high"`
	Low           string `json:"low"`
	High52w       string `json:"high52w"`
	Low52w        string `json:"low52w"`
	Marketcapital string `json:"marketcapital"`
	Amount        string `json:"amount"`
	Type          string `json:"type"`
	Pettm         string `json:"pettm"`
	Volume        string `json:"volume"`
	Hasexist      string `json:"hasexist"`
}

type KlineDataResp struct {
	Data             KlineData `json:"data"`
	ErrorCode        int       `json:"error_code"`
	ErrorDescription string    `json:"error_description"`
}

type KlineData struct {
	Symbol string      `json:"symbol"`
	Column []string    `json:"column"`
	Item   [][]float64 `json:"item"`
}
