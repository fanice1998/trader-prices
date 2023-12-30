package tools

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Client struct {
	BaseURL  string
	symbol   string
	interval string
}

func HistoricalPrice(client Client) {
	//url := "https://data.binance.vision/?prefix=data/futures/um/monthly/klines/ETHUSDT/1h/"
	fmt.Println("Get URL string: ", client.BaseURL)

	// 返回合約交易所有USDT交易對 symbol
	// 取得xml 標籤結果
	var parser []string
	prefix := "data/futures/um/monthly/klines"
	if client.interval != "" && client.symbol != "" {
		// 取得歷史數據路徑
		parser = xmlParser(fmt.Sprintf("%s/%s/%s/", prefix, client.symbol, client.interval))
	} else if client.interval == "" && client.symbol != "" {
		// 取得時間軸分類
		parser = xmlParser(fmt.Sprintf("%s/%s/", prefix, client.symbol))
	} else if client.symbol == "" {
		// 取得symbols
		parser = xmlParser(prefix + "/")
	} else {
		// 發生錯誤
		parser = append(parser, "Error message: client class configuration error.")
	}
	fmt.Println(parser)
}

func xmlParser(prefix string) []string {
	result := make([]string, 0, 150)
	xmlurl := "https://s3-ap-northeast-1.amazonaws.com/data.binance.vision?delimiter=/&prefix=" + prefix
	fmt.Println("Xml url is ", xmlurl)

	// 請求url
	req, err := http.NewRequest("GET", xmlurl, nil)
	if err != nil {
		log.Println(err)
	}
	conn := &http.Client{}
	resp, err := conn.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	// 解析xml響應
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	// 將取得的xml解析成所需要的格式
	// Parse the obtained XML into the required format.
	content := new(ListBuckResult)
	err = xml.Unmarshal(body, content)

	// 取得xml標籤內容
	if len(content.Contents) != 0 {
		// 取得路徑
		for _, v := range content.Contents {
			result = append(result, v.Key)
			log.Println("content.Contents", v)
		}
	} else if len(content.Prefix) != 0 {
		// 取得檔案名稱
		for _, v := range content.CommonPrefixes {
			result = append(result, v.Prefix)
			log.Println("content.Prefix", v)
		}
	} else {
		// 未搜尋到匹配的項目
		result = append(result, "Error message: not matching category found.")
	}
	return result

}

type ListBuckResult struct {
	// xml 格式
	Prefix         string           `xml:"Prefix"`
	CommonPrefixes []CommonPrefixes `xml:"CommonPrefixes"`
	Contents       []Contents       `xml:"Contents"`
}

type CommonPrefixes struct {
	Prefix string `xml:"Prefix"`
}

type Contents struct {
	Key string `xml:"Key"`
}
