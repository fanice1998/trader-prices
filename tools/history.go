package tools

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Client struct {
	BaseURL  string
	Symbol   string
	Interval string
}

func HistoricalPrice(client Client) {
	//url := "https://data.binance.vision/?prefix=data/futures/um/monthly/klines/ETHUSDT/1h/"
	fmt.Println("Get URL string: ", client.BaseURL)

	// 返回合約交易所有USDT交易對 symbol
	// 取得xml 標籤結果
	var parser []string
	prefix := "data/futures/um/monthly/klines"
	if client.Interval != "" && client.Symbol != "" {
		// 取得歷史數據路徑
		parser = xmlParser(fmt.Sprintf("%s/%s/%s/", prefix, client.Symbol, client.Interval))
	} else if client.Interval == "" && client.Symbol != "" {
		// 取得時間軸分類
		parser = xmlParser(fmt.Sprintf("%s/%s/", prefix, client.Symbol))
	} else if client.Symbol == "" {
		// 取得symbols
		parser = xmlParser(prefix + "/")
	} else {
		// 發生錯誤
		parser = append(parser, "Error message: client class configuration error.")
	}

	// 開始下載檔案
	fmt.Println("Download file...")
	// 確認目錄存在
	path := filepath.Join("data", client.Symbol, client.Interval)
	checkDir(path)
	for _, v := range parser {
		if !strings.Contains(v, "CHECKSUM") {
			// 下載檔案
			fileName := downloadData(v, client)
			path := filepath.Join("data", client.Symbol, client.Interval)
			err := unzipSource(filepath.Join(path, fileName), path)
			if err != nil {
				log.Println(err)
			}
			break
		}
	}
}

func unzipSource(zipFilePath, extracTo string) error {
	r, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return err
	}
	defer func(r *zip.ReadCloser) {
		err := r.Close()
		if err != nil {

		}
	}(r)

	err = os.MkdirAll(extracTo, os.ModeDir)
	if err != nil {
		return err
	}

	for _, file := range r.File {
		rc, err := file.Open()
		if err != nil {
			return err
		}
		defer func(rc io.ReadCloser) {
			err := rc.Close()
			if err != nil {

			}
		}(rc)

		tragetFilePath := filepath.Join(extracTo, file.Name)

		if file.FileInfo().IsDir() {
			err := os.MkdirAll(tragetFilePath, os.ModeDir)
			if err != nil {
				return err
			}
		} else {
			err := os.MkdirAll(extracTo, os.ModeDir)
			if err != nil {
				return err
			}
			w, err := os.Create(tragetFilePath)
			if err != nil {
				return err
			}
			defer func(w *os.File) {
				err := w.Close()
				if err != nil {

				}
			}(w)

			_, err = io.Copy(w, rc)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func downloadData(prefix string, client Client) string {

	// 請求目標URL後下載檔案
	url := "https://data.binance.vision/"
	// 下載網址
	downloadUrl := url + prefix
	fmt.Println(downloadUrl)
	req, err := http.NewRequest("GET", downloadUrl, nil)
	if err != nil {
		log.Println(err)
	}
	conn := &http.Client{}
	resp, err := conn.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	// 寫入檔案
	splitResult := strings.Split(prefix, "/")
	fileName := splitResult[len(splitResult)-1:][0]
	path := filepath.Join("data", client.Symbol, client.Interval, fileName)
	file, err := os.Create(path)
	if err != nil {
		log.Println("Write file fail. Err msg: ", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Println(err)
		}
	}(file)

	_, err = file.Write(body)
	if err != nil {
		fmt.Println("寫入時發生問題. Err: ", err)
	}

	return fileName
}

func creatDir(name string) {
	// 取得當前路徑
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	// 建立符合系統的目錄位址
	pathFile := filepath.Join(path, name)
	err = os.MkdirAll(pathFile, 0755)
	if err != nil {
		log.Println(err)
	}
}

// 檢查目錄是否存在
func checkDir(file string) {
	// 取得當前路徑
	currentDir, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	_, err = os.Stat(currentDir + "/" + file)
	if os.IsNotExist(err) {
		log.Println("Err: Directory does not exist.", currentDir+file)
		log.Println("Create Directory... ", file)
		creatDir(file)
	} else if err != nil {
		log.Println("Error occurred while checking the directory. Error message: ", err)
	} else {
		fmt.Println(filepath.Join(currentDir, file), "已存在")
	}

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
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)

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
			//log.Println("content.Contents", v)
		}
	} else if len(content.Prefix) != 0 {
		// 取得檔案名稱
		for _, v := range content.CommonPrefixes {
			result = append(result, v.Prefix)
			//log.Println("content.Prefix", v)
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
