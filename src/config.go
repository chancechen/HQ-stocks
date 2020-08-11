package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func readStockFile(file string) (b []byte) {

	nildata := make([]byte, 0)

	fp, err := os.OpenFile(file, os.O_RDONLY, 0777)

	if err != nil {

		log.Printf("fail to open file::%s\n", err)
		return nildata
	}

	defer fp.Close()

	data, err := ioutil.ReadAll(fp)

	if err != nil {

		log.Printf("fail to read file::%s\n", err)
		return nildata
	}

	return data
}

func readConfigFile() {

	if err := json.Unmarshal(readStockFile(confjson), &confData); err != nil {

		log.Printf("fail to unmarshal:%s\n", err)
		return
	}

	log.Printf("success to load %s,def ROE[%.2f-:%.2f],Forceload:%t show ROE[%.2f - %.2f]\n",
		confjson,
		confData.MinRoe,
		confData.MaxRoe,
		confData.ForceLoad,
		confData.ShowMinRoe,
		confData.ShowMaxRoe)
}

func writeStockFile() {

	fp, err := os.OpenFile(stockjson, os.O_RDWR|os.O_CREATE, 0777)

	if err != nil {
		log.Printf("fail to open file::%s\n", err)
		return
	}

	defer fp.Close()

	if err = json.NewEncoder(fp).Encode(&roedatas); err != nil {

		log.Printf("fail to write file::%s\n", err)
		return
	}

	log.Printf("success to update %s\n", stockjson)
}

// 获取本地数据
func readLoacalData() {

	if err := json.Unmarshal(readStockFile(stockjson), &roedatas); err != nil {

		log.Printf("fail to unmarshal:%s\n", err)
		return
	}
}

// 是否是最新的数据
func isNewData() bool {

	return time.Now().Format("2006-01-02") == roedatas.Updatetime.Format("2006-01-02")
}

func loadData() {

	readConfigFile() // 加载配置文件
	readLoacalData() // 加载本地文件

	if confData.ForceLoad { //如果需要强制加载第三方数据平台数据

		initToken()
		getStockData()
		getPePbData()

		writeStockFile()

		log.Printf("success to load stock data frome : %s\n", joinquantAPI)
	}

	log.Printf("success to load config and stock data !\n") // 成功加载数据

	// 根据配置设置主机信息
	if confData.Host != "" {

		host = confData.Host
	}

	// 根据配置信息过滤展示的股票信息
	initShowRoeData()
	getPePbData()
	getRadarData()
	getPePbData()
}
