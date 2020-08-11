package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	joinquantAPI = "https://dataapi.joinquant.com/apis"
)

var (
	token    = ""
	stockID  = []string{""}
	roeYears = []int{2015, 2016, 2017, 2018, 2019}
)

func joinquantRequest(body map[string]interface{}) (string, error) {

	bodyStr, err := json.Marshal(body)
	client := &http.Client{}
	req, err := http.NewRequest("POST", joinquantAPI, strings.NewReader(string(bodyStr)))
	resp, err := client.Do(req)
	if err != nil {

	}
	defer resp.Body.Close()

	res, err := ioutil.ReadAll(resp.Body)

	if err == nil {
		return string(res), err
	}

	fmt.Println(err)
	return "", err
}

// 获取指定股票的基本信息
func getStock(stocks []string, date string) []string {

	body := map[string]interface{}{
		"method":  "get_fundamentals",
		"token":   token,
		"table":   "valuation",
		"columns": "code,day,pb_ratio,pe_ratio,capitalization,market_cap",
		"code":    strings.Join(stocks, ","),
		"date":    date,
		//"count":   12,
	}

	res, _ := joinquantRequest(body)
	if res != "" {

		return strings.Fields(res)[1:]
	}

	return make([]string, 0)
}

// 获取股票ROE信息
func getROE(stocks []string, date string) []string {

	body := map[string]interface{}{
		"method":  "get_fundamentals",
		"token":   token,
		"table":   "indicator",
		"columns": "code,statDate,roe",
		"code":    strings.Join(stocks, ","),
		"date":    date,
	}

	res, _ := joinquantRequest(body)

	if res != "" {

		return strings.Fields(res)[1:]
	}

	return make([]string, 0)
}

func initToken() error {

	body := map[string]interface{}{
		"method": "get_token",
		"mob":    confData.APIAccount,
		"pwd":    confData.APIPwd,
	}

	var err error = nil

	token, err = joinquantRequest(body)
	if err != nil {

		log.Printf("fail to get token, err %s\n", err)
	}

	log.Printf("success to get token: %s\n", token)

	return err
}

// 获取所有股票信息
func allStocks() ([]string, error) {

	body := map[string]interface{}{

		"method": "get_all_securities",
		"token":  token,
		"code":   "stock",
		"date":   "2020-08-05",
	}

	res, err := joinquantRequest(body)

	if err == nil {

		return strings.Fields(res)[1:], err
	}
	return make([]string, 0), err
}

// 今日查询次数
func queryTimes() (int, int) {

	body := map[string]interface{}{

		"method": "get_query_count",
		"token":  token,
	}

	res, err := joinquantRequest(body)
	if err != nil {

		log.Printf("fail to queryTimes file::%s\n", err)
		return 0, 1000000

	}

	times, err := strconv.Atoi(res)
	if err == nil {
		return times, 1000000 - times
	}
	log.Printf("fail to convert %s to number, err %s\n", res, err)

	return 0, 1000000
}

// 过滤ROE
func filterRoe(min float32, max float32) {

	var (
		total float32
		avg   float32
	)
	count := 0
	cleanupCnt := 0
	for k, v := range roedatas.Data {

		total = 0
		count = 0
		for _, c := range v {

			if d, err := strconv.ParseFloat(c, 32); err == nil {

				total += float32(d)
				count++
			} else {

				log.Printf("%s,%s\n", err, k)
			}
		}
		avg = total / float32(count)

		if count == 0 || count != len(roeYears) || min > avg || avg > max {

			delete(roedatas.Data, k)
			cleanupCnt++
		}
	}
}

// 根据过滤出来的数据，再次获取其PePb数据
func getPePbData() {

	if roedatas.PbData == nil ||
		roedatas.PeData == nil {

		var roeStocks []string
		var data []string

		roedatas.PbData = make(map[string][]string)
		roedatas.PeData = make(map[string][]string)

		for k := range roedatas.Data {

			roeStocks = append(roeStocks, k)
		}

		if len(roeStocks) > 0 {

			for index, i := range roeYears {

				data = getStock(roeStocks, strconv.Itoa(i))
				for _, item := range data {

					d := strings.Split(item, ",")
					if roedatas.PeData[d[0]] == nil {
						roedatas.PeData[d[0]] = make([]string, len(roeYears)) // 和years对应
					}
					if roedatas.PbData[d[0]] == nil {
						roedatas.PbData[d[0]] = make([]string, len(roeYears)) // 和years对应
					}

					if f, err := strconv.ParseFloat(d[2], 64); err == nil {
						roedatas.PbData[d[0]][index] = fmt.Sprintf("%.2f", f)
					} else {
						roedatas.PbData[d[0]][index] = d[2]
					}
					if f, err := strconv.ParseFloat(d[3], 64); err == nil {
						roedatas.PeData[d[0]][index] = fmt.Sprintf("%.2f", f)
					} else {
						roedatas.PeData[d[0]][index] = d[2]
					}
				}
			}
		}
	}
}

func getStockData() {

	list, err := allStocks()
	if err != nil {

		log.Printf("fail to get all stocks from:%s", joinquantAPI)
		return
	}

	roedatas.Name = make(map[string]string)

	for _, stock := range list {

		if d := strings.Split(stock, ","); len(d) > 1 {

			roedatas.Name[d[0]] = d[1]

			stockID = append(stockID, d[0])
		}
	}

	var data []string

	// 初始化
	roedatas.Data = make(map[string][]string)
	roedatas.StatDate = make([]string, 0)
	roedatas.Updatetime = time.Now()
	roedatas.PbData = nil
	roedatas.PeData = nil

	for k, i := range roeYears {

		roedatas.StatDate = append(roedatas.StatDate, strconv.Itoa(i)+"年")

		data = getROE(stockID, strconv.Itoa(i))
		for _, item := range data {
			d := strings.Split(item, ",")
			if roedatas.Data[d[0]] == nil {
				roedatas.Data[d[0]] = []string{"0", "0", "0", "0", "0"} // 和years对应
			}
			if f, err := strconv.ParseFloat(d[2], 64); err == nil {
				roedatas.Data[d[0]][k] = fmt.Sprintf("%.2f", f)
			} else {
				roedatas.Data[d[0]][k] = "0.0"
			}
		}
	}

	filterRoe(confData.MinRoe, confData.MaxRoe)
}

func initShowRoeData() {

	var (
		total float32
		avg   float32
	)

	count := 0

	var data map[string][]string = make(map[string][]string, 0)

	for k, v := range roedatas.Data {

		total = 0
		count = 0
		for _, c := range v {

			if d, err := strconv.ParseFloat(c, 32); err == nil {
				if d < 0 {
					count = 0
					break
				}

				total += float32(d)
				count++
			} else {

				log.Printf("%s,%s\n", err, k)
			}
		}
		if count > 0 {

			avg = total / float32(count)

			if confData.ShowMinRoe < avg && avg < confData.ShowMaxRoe {

				data[k] = v
			}
		}
	}

	roeShowData = data
}
