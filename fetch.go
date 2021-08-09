package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

const screenerURL = "https://api.stockbit.com/v2.4/screener"
const authorizationKey = "Authorization"
const totalPage = 31

func fetchTickerList() ([]string, error) {
	var res []string
	c := make(chan []string)
	count := 0
	for i := 1; i <= totalPage; i++ {
		go func(i int) {
			list, err := fetchTickerListPage(i)
			if err != nil {
				log.Println(err)
				return
			}
			c <- list
			count++
		}(i)
	}
	for v := range c {
		res = append(res, v...)
		if count == totalPage {
			close(c)
		}
	}
	sort.Strings(res)
	return res, nil
}

func fetchTickerListPage(page int) ([]string, error){
	data := url.Values{}
	data.Set("universe", "{\"scope\":\"IHSG\",\"scopeID\":\"\",\"name\":\"IHSG\"}")
	data.Set("filters", "[{\"type\":\"basic\",\"item1\":2892,\"item2\":0,\"operator\":\">\",\"multiplier\":\"\"}]")
	data.Set("page", fmt.Sprintf("%d", page))

	req, err := http.NewRequest(http.MethodPost, screenerURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add(authorizationKey, authToken)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var model ScreenerResponse
	if err := json.Unmarshal(b, &model); err != nil {
		return nil, err
	}

	return model.TickerList(), nil
}

func fetchFinancial(ticker string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, financialURL(ticker), nil)
	if err != nil {
		return "", err
	}
	req.Header.Add(authorizationKey, authToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var model FinancialResponse
	if err := json.Unmarshal(b, &model); err != nil {
		return "", err
	}

	return model.Data.HTMLReport, nil
}

func financialURL(ticker string) string {
	const url = "https://api.stockbit.com/v2.4/company/financial?symbol=%s&datatype=reported&reporttype=is&statement=ANNUALLY"
	return fmt.Sprintf(url, ticker)
}

type FinancialResponse struct {
	Message string `json:"message"`
	Data struct {
		HTMLReport string `json:"htmlReport"`
	} `json:"data"`
}

type ScreenerResponse struct {
	Data struct{
		Calcs []struct{
			Company struct{
				Symbol string `json:"symbol"`
			} `json:"company"`
		} `json:"calcs"`
	} `json:"data"`
}

func (s ScreenerResponse) TickerList() []string {
	var res []string
	for _, v := range s.Data.Calcs {
		res = append(res, v.Company.Symbol)
	}
	return res
}