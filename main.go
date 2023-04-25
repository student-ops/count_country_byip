package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type IPInfo struct {
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	ISP         string  `json:"isp"`
	Org         string  `json:"org"`
	AS          string  `json:"as"`
	Query       string  `json:"query"`
}

func getIPInfo(ip string) (*IPInfo, error) {
	url := fmt.Sprintf("http://ip-api.com/json/%s", ip)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// レスポンスボディをバイト列に読み込む
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// レスポンスをパースしてIP情報を取得する
	var info IPInfo
	err = json.Unmarshal(body, &info)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

func main() {
	// ファイルを開く
	file, err := os.Open("data/aws_info.csv")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// CSVリーダーを作成し、レコードを読み込む
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV:", err)
		return
	}

	// IPアドレスのスライスを作成する
	var ipSlice []string
	for i, row := range records {
		if i == 0 {
			continue // ヘッダー行をスキップする
		}
		ipSlice = append(ipSlice, row[0])
	}
	accessCounts := make(map[string]int)

	for _, ip := range ipSlice {
		info, err := getIPInfo(ip)
		if err != nil {
			fmt.Println("Error getting IP info:", err)
			continue
		}

		if info.Status == "success" {
			fmt.Println(info)
			accessCounts[info.Country]++
		}
	}

	// 結果を出力する
	fmt.Println("Access counts by country:")
	for country, count := range accessCounts {
		fmt.Printf("%s: %d\n", country, count)
	}

}
