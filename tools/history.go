package tools

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type Client struct {
	BaseURL string
}

func HistoricalPrice(client Client) {
	//url := "https://data.binance.vision/?prefix=data/futures/um/monthly/klines/ETHUSDT/1h/"
	//xml := "https://s3-ap-northeast-1.amazonaws.com/data.binance.vision?delimiter=/&prefix=data/futures/"
	fmt.Println("Get URL string: ", client.BaseURL)

	// 返回合約交易所有USDT交易對 symbol
	parser := xmlParser("data/futures/um/monthly/klines/ETHUSDT/")
	fmt.Println(parser)
}

func xmlParser(prefix string) map[string]string {
	result := make(map[string]string)
	xmlurl := "https://s3-ap-northeast-1.amazonaws.com/data.binance.vision?delimiter=/&prefix=" + prefix
	fmt.Println("Xml url is ", xmlurl)

	lastString := strings.Split(xmlurl, "/")[len(strings.Split(xmlurl, "/"))-2]

	fmt.Println(lastString)
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	// 將取得的xml解析成所需要的格式
	// Parse the obtained XML into the required format.
	content := new(ListBuckResult)
	err = xml.Unmarshal(body, content)

	// 為了方便可視化結果
	switch {
	case lastString == "klines":
		// 只記錄使用USDT的交易對
		// Only record trading pairs using USDT.
		lenP := len(prefix) // prefix lens
		for _, v := range content.CommonPrefixes {
			if v.Prefix[len(v.Prefix)-5:len(v.Prefix)-1] == "USDT" {
				result[v.Prefix[lenP:len(v.Prefix)-1]] = v.Prefix
			}
		}
		return result

	case lastString[len(lastString)-4:] == "USDT":
		for _, v := range content.CommonPrefixes {
			lastS := strings.LastIndex(v.Prefix, "/")
			fmt.Println(v.Prefix[lastS:])
			result["intervar"] += v.Prefix + "/"
		}
		return result

	default:
		// 未搜尋到匹配的類別
		//result["err"] = "No matching category found! "
		for _, v := range content.CommonPrefixes {
			result[v.Prefix] = v.Prefix
		}
		return result
	}

}

type ListBuckResult struct {
	Prefix         string           `xml:"Prefix"`
	CommonPrefixes []CommonPrefixes `xml:"CommonPrefixes"`
}

type CommonPrefixes struct {
	Prefix string `xml:"Prefix"`
}
