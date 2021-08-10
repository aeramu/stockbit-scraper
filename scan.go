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

	doc.Find("tr").Each(func(i int, row *goquery.Selection) {
		head := row.Find("span").AttrOr("data-lang-0", "")

		if funk.ContainsString(index, head) {
			row.Find("td").Each(func(i int, cell *goquery.Selection) {
				if i == 0 {
					return
				}
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
	switch head {
	case index[0]:
		t.Revenue = append(t.Revenue, Int64(data))
	case index[1]:
		t.GrossProfit = append(t.GrossProfit, Int64(data))
	case index[2]:
		t.OperatingProfit = append(t.OperatingProfit, Int64(data))
	case index[4]:
		t.NetProfit = append(t.NetProfit, Int64(data))
	case index[5]:
		t.EPS = append(t.EPS, Float32(data))
	case index[6]:
		t.PER = append(t.PER, Float32(data))
	}
}

func Int64(s string) int64 {
	dotIndex := strings.Index(s, ".")
	if dotIndex != -1 {
		s = s[:dotIndex]
	}
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		log.Println("[Parse Int]", err, s)
	}
	return n
}

func Float32(s string) float32 {
	f, err := strconv.ParseFloat(s, 32)
	if err != nil {
		log.Println("[Parse Float]", err, s)
	}
	return float32(f)
}

func reverse(t *ticker) *ticker{
	t.Revenue = funk.ReverseInt64(t.Revenue)
	t.GrossProfit = funk.ReverseInt64(t.GrossProfit)
	t.OperatingProfit = funk.ReverseInt64(t.OperatingProfit)
	t.NetProfit = funk.ReverseInt64(t.NetProfit)
	t.EPS = funk.ReverseFloat32(t.EPS)
	t.PER = funk.ReverseFloat32(t.PER)
	return t
}

var index = []string{
	"Total Pendapatan",
	"Laba Kotor",
	"Laba Usaha",
	"Laba Sebelum Pajak",
	"Laba Bersih Tahun Berjalan",
	"EPS (Annual)",
	"PE Ratio (Annual)",
}
