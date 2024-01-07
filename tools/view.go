package tools

import (
	"encoding/csv"
	"github.com/go-echarts/go-echarts/charts"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func View(client Client) {
	kd := readCsvKline("ETHUSDT-1h-2020-01.csv", client)
	kline := charts.NewKLine()
	x := make([]string, 0)
	y := make([][4]float32, 0)
	for i := 0; i < len(kd); i++ {
		x = append(x, kd[i].date)
		y = append(y, kd[i].data)
	}

	kline.AddXAxis(x).AddYAxis("kline", y)
	kline.SetGlobalOptions(
		charts.TitleOpts{Title: "Kline simple"},
		charts.XAxisOpts{SplitNumber: 20},
		charts.YAxisOpts{Scale: true},
		charts.DataZoomOpts{XAxisIndex: []int{0}, Start: 50, End: 100},
	)

	kline.SetSeriesOptions(
		charts.ItemStyleOpts{
			Color:        "#FFFFFF",
			Color0:       "#000000",
			BorderColor:  "#000000",
			BorderColor0: "#000000",
		},
	)

	f, err := os.Create("klineData.html")
	if err != nil {
		log.Println(err)
	}
	err = kline.Render(f)
	if err != nil {
		log.Println(err)
	}

}

func readCsvKline(fileName string, client Client) []klineData {
	csvFile, err := os.Open(filepath.Join("data", client.Symbol, client.Interval, fileName))
	if err != nil {
		log.Println(err)
	}
	defer func(csvFile *os.File) {
		err := csvFile.Close()
		if err != nil {
			log.Println(err)
		}
	}(csvFile)

	reader := csv.NewReader(csvFile)
	records, err := reader.ReadAll()
	if err != nil {
		log.Println(err)
	}

	var data []klineData
	for _, v := range records {
		v1, err := strconv.ParseFloat(v[1], 32)
		if err != nil {
			log.Println(err)
		}
		v2, err := strconv.ParseFloat(v[2], 32)
		if err != nil {
			log.Println(err)
		}
		v3, err := strconv.ParseFloat(v[3], 32)
		if err != nil {
			log.Println(err)
		}
		v4, err := strconv.ParseFloat(v[4], 32)
		if err != nil {
			log.Println(err)
		}
		data = append(
			data,
			klineData{
				date: v[0],
				data: [4]float32{
					float32(v1),
					float32(v4),
					float32(v3),
					float32(v2),
				},
			})
	}
	return data
}

type klineData struct {
	date string
	data [4]float32 // [open, high, low, close] chart is [open close low high]
}
