package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-echarts/go-echarts/charts"
)

// 基础数据
var (
	sname           string = "股票名字"
	snameItems      []string
	peItems         []string
	pbItems         []string
	mkcapItems      []string
	capItems        []string
	roedatas        roeData
	confData        config
	roeShowData     map[string][]string
	radarData       map[string][][]float32
	stockIndicators = []charts.IndicatorOpts{
		{Name: "Roe", Max: 300},
		{Name: "Pe", Max: 250},
		{Name: "Pb", Max: 300},
	}
)

func stockROEbar() *charts.Bar {

	title := fmt.Sprintf("%d - %d 年，ROE[%.2f -%.2f] 区间统计图", roeYears[0], roeYears[len(roeYears)-1], confData.ShowMinRoe, confData.ShowMaxRoe)

	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.InitOpts{PageTitle: "ROE统计图", Theme: charts.ThemeType.Shine, Width: "100%"},
		charts.TitleOpts{Title: title, Right: "35%"},
		charts.LegendOpts{Bottom: "0%"},
		//charts.DataZoomOpts{XAxisIndex: []int{0}, Start: 50, End: 100},
		charts.ToolboxOpts{Show: true})

	bar.AddXAxis(roedatas.StatDate)
	for k, v := range roeShowData {
		if roedatas.Name[k] != "" {
			bar.AddYAxis(roedatas.Name[k], v)
		} else {
			bar.AddYAxis(k, v)
		}
	}

	//bar.SetGlobalOptions(charts.YAxisOpts{SplitLine: charts.SplitLineOpts{Show: true}})
	bar.SetSeriesOptions(charts.LabelTextOpts{Show: true})

	return bar
}

func stockPebar() *charts.Bar {

	title := fmt.Sprintf("%d - %d 年，Pe统计图", roeYears[0], roeYears[len(roeYears)-1])

	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.InitOpts{PageTitle: "ROE统计图", Theme: charts.ThemeType.Macarons, Width: "100%"},
		charts.TitleOpts{Title: title, Right: "35%"},
		charts.LegendOpts{Bottom: "0%"},
		//charts.DataZoomOpts{XAxisIndex: []int{0}, Start: 50, End: 100},
		charts.ToolboxOpts{Show: true})

	bar.AddXAxis(roedatas.StatDate)
	for k, v := range roedatas.PeData {

		if roeShowData[k] != nil {

			if roedatas.Name[k] != "" {
				bar.AddYAxis(roedatas.Name[k], v)
			} else {
				bar.AddYAxis(k, v)
			}
		}
	}

	//bar.SetGlobalOptions(charts.YAxisOpts{SplitLine: charts.SplitLineOpts{Show: true}})
	bar.SetSeriesOptions(charts.LabelTextOpts{Show: true})

	return bar
}

func getAvg(data []string) float32 {

	var f float32
	for _, v := range data {
		if d, err := strconv.ParseFloat(v, 32); err == nil {
			f += float32(d)
		}
	}
	if len(data) > 0 {

		return f / float32(len(data))
	}

	return 0.0
}

func getRadarData() {

	if radarData == nil {

		radarData = make(map[string][][]float32, 0)

		radarInfo := map[string]float32{"Roe": 0.0, "Pe": 0.0, "Pb": 0.0}

		// ROE、PE、PB
		for k := range roeShowData {

			roe := getAvg(roeShowData[k])    //ROE 均值
			pe := getAvg(roedatas.PeData[k]) //PE 均值
			pb := getAvg(roedatas.PbData[k]) //PB 均值

			radarData[k] = [][]float32{{roe, pe, pb}}

			if radarInfo["Roe"] < roe {
				radarInfo["Roe"] = roe
			}
			if radarInfo["Pb"] < pb {
				radarInfo["Pb"] = pb
			}
			if radarInfo["Pe"] < pe {
				radarInfo["Pe"] = pe
			}
		}

		radarInfo["Roe"] += 10.0
		radarInfo["Pb"] += 10.0
		radarInfo["Pe"] += 10.0

		// ROE、PE越小越好，则取其反数
		for k := range radarData {

			radarData[k][0][1] = radarInfo["Pe"] - radarData[k][0][1]
			radarData[k][0][2] = radarInfo["Pb"] - radarData[k][0][2]
		}

		// 更新
		for k := range stockIndicators {
			stockIndicators[k].Max = radarInfo[stockIndicators[k].Name]
		}
	}
}

func stockRadar() *charts.Radar {

	title := fmt.Sprintf("%d - %d 年，优质公司雷达图", roeYears[0], roeYears[len(roeYears)-1])
	radar := charts.NewRadar()
	radar.SetGlobalOptions(
		charts.InitOpts{PageTitle: "优质公司雷达图", Theme: charts.ThemeType.Macarons, Width: "100%"},
		charts.TitleOpts{
			Title: title,
			Right: "center"},

		charts.LegendOpts{Bottom: "0%"},
		charts.RadarComponentOpts{
			Indicator: stockIndicators,
			SplitLine: charts.SplitLineOpts{Show: true},
			SplitArea: charts.SplitAreaOpts{Show: true},
		},
	)

	for k, v := range radarData {

		color := (int(seed.Int63()) % 0xffffff)
		if roedatas.Name[k] != "" {
			radar.Add(roedatas.Name[k], v, charts.ItemStyleOpts{Color: fmt.Sprintf("#%X", color)})
		} else {
			radar.Add(k, v, charts.ItemStyleOpts{Color: fmt.Sprintf("#%X", color)})
		}

	}

	return radar
}

func stockPbbar() *charts.Bar {

	title := fmt.Sprintf("%d - %d 年，Pb统计图", roeYears[0], roeYears[len(roeYears)-1])

	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.InitOpts{PageTitle: "ROE统计图", Theme: charts.ThemeType.Roma, Width: "100%"},
		charts.TitleOpts{Title: title, Right: "35%"},
		charts.LegendOpts{Bottom: "0%"},
		//charts.DataZoomOpts{XAxisIndex: []int{0}, Start: 50, End: 100},
		charts.ToolboxOpts{Show: true})

	bar.AddXAxis(roedatas.StatDate)
	for k, v := range roedatas.PbData {

		if roeShowData[k] != nil {

			if roedatas.Name[k] != "" {
				bar.AddYAxis(roedatas.Name[k], v)
			} else {
				bar.AddYAxis(k, v)
			}
		}
	}

	//bar.SetGlobalOptions(charts.YAxisOpts{SplitLine: charts.SplitLineOpts{Show: true}})
	bar.SetSeriesOptions(charts.LabelTextOpts{Show: true})

	return bar
}

func queryTimesPie() *charts.Pie {

	ltimes, qtimes := queryTimes()

	times := make(map[string]interface{})

	times["Left Times"] = ltimes
	times["query Times"] = qtimes

	pie := charts.NewPie()
	pie.SetGlobalOptions(
		charts.TitleOpts{Title: "当日查询情况"},
		charts.ToolboxOpts{Show: true},
		charts.InitOpts{PageTitle: "优质股票统计图"},
	)
	pie.Add("pie", times,
		charts.LabelTextOpts{Show: true, Formatter: "{b}: {c}"},
		//charts.PieOpts{Radius: []string{percentLeft, percentQuery}},
		charts.PieOpts{Radius: []string{"45%", "75%"}},
	)

	return pie
}

func barHandler(w http.ResponseWriter, _ *http.Request) {

	page := charts.NewPage(orderRouters("bar")...)

	page.Add(
		stockROEbar(),
		stockPebar(),
		stockPbbar(),
		stockRadar(),
	)

	f, err := os.Create(getRenderPath("bar.html"))
	if err != nil {

		log.Println(err)
	}

	page.Render(w, f)
}
