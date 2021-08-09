package main

import (
	"log"
)

const maxConcurrent = 50
const authToken = ""
const mongodbURL = ""

func main() {
	if err := initDB(); err != nil {
		log.Println(err)
		return
	}

	tickers, err := fetchTickerList()
	if err != nil {
		log.Println(err)
		return
	}

	c := make(chan *ticker)
	count := 0
	concurrentCount := 0

	go func() {
		for _, v := range tickers {
			for concurrentCount > maxConcurrent {}
			concurrentCount++
			go func(v string) {
				s, err := fetchFinancial(v)
				if err != nil {
					log.Println(v, err)
					return
				}
				t := new(ticker)
				t.Ticker = v
				if err := scan(s, t); err != nil {
					log.Println(t.Ticker, err)
					return
				}
				c <- t
				concurrentCount--
				count++
			}(v)
		}
	}()

	for t := range c {
		if err := saveTicker(t); err != nil {
			log.Println(t.Ticker, err)
		}
		log.Println(count)
		if count == len(tickers) {
			close(c)
		}
	}
}

type ticker struct {
	Ticker          string  `bson:"symbol"`
	Revenue         []int64 `bson:"revenue"`
	GrossProfit     []int64 `bson:"gross_profit"`
	OperatingProfit []int64 `bson:"operating_profit"`
	NetProfit       []int64 `bson:"net_profit"`
}
