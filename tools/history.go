package tools

import (
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	BaseURL string
}

func HistoricalPrice(client Client) {
	url := "https://data.binance.vision/?prefix=data/futures/um/monthly/klines/ETHUSDT/1h/"
	fmt.Println("Get URL string: ", client.BaseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	conn := &http.Client{}
	resp, err := conn.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(body))
}
