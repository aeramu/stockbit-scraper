package main

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/thoas/go-funk"
	"log"
	"strconv"
	"strings"
)

func scan(s string, t *ticker) error {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(s))
	if err != nil {
		return err
	}
	table := doc.Find("table").First()

	table.Find("tbody tr").Each(func(i int, row *goquery.Selection) {
		head := row.Find("span").AttrOr("data-lang-0", "")

		if funk.ContainsString(index, head) {
			row.Find("td.rowval").Each(func(i int, cell *goquery.Selection) {
				data := cell.AttrOr("data-raw", "")
				insert(t, head, data)
			})
		}
	})

	t = reverse(t)
	return nil
}

func insert(t *ticker, head string, data string) {
	if data == "" {
		data = "0"
	}
	dotIndex := strings.Index(data, ".")
	if dotIndex != -1 {
		data = data[:dotIndex]
	}
	n, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		log.Println(t.Ticker, err)
	}
	switch head {
	case index[0]:
		t.Revenue = append(t.Revenue, n)
	case index[1]:
		t.GrossProfit = append(t.GrossProfit, n)
	case index[2]:
		t.OperatingProfit = append(t.OperatingProfit, n)
	case index[4]:
		t.NetProfit = append(t.NetProfit, n)
	}
}

func reverse(t *ticker) *ticker{
	t.Revenue = funk.ReverseInt64(t.Revenue)
	t.GrossProfit = funk.ReverseInt64(t.GrossProfit)
	t.OperatingProfit = funk.ReverseInt64(t.OperatingProfit)
	t.NetProfit = funk.ReverseInt64(t.NetProfit)
	return t
}

var index = []string{"Total Pendapatan", "Laba Kotor", "Laba Usaha", "Laba Sebelum Pajak", "Laba Bersih Tahun Berjalan"}
