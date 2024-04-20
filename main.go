package main

import (
	"bytes"
	"encoding/csv"
	"io"
	"net/http"
	"os"
	"sort"

	"github.com/PuerkitoBio/goquery"
)

func getToukijoCode() (map[string]string, error) {
	src := "https://www.touki-kyoutaku-online.moj.go.jp/toukinet/mock/SC01WS01.html"

	resp, err := http.Get(src)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	buff := bytes.NewBuffer(body)
	dom, err := goquery.NewDocumentFromReader(buff)
	if err != nil {
		return nil, err
	}
	codes := make(map[string]string)
	dom.Find("table tr").Each(func(i int, s *goquery.Selection) {
		if i > 0 {
			code := s.Find("td").Eq(0).Text()
			name := s.Find("td").Eq(1).Text()
			codes[code] = name
		}
	})

	// 以下の2つを削除
	// 0000: 全登記所
	delete(codes, "0000")

	// 登記所コード: 登記所名
	delete(codes, "登記所コード")

	return codes, nil
}

func main() {
	codes, err := getToukijoCode()
	if err != nil {
		panic(err)
	}
	var keys []string
	for key := range codes {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	// CSVとして出力
	file, err := os.Create("toukijo.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Write([]string{"code", "name"})
	for _, code := range keys {
		writer.Write([]string{code, codes[code]})
	}
	writer.Flush()

	// Markdownとして出力
	file, err = os.Create("README.md")
	if err != nil {
		panic(err)
	}

	file.WriteString("# 登記所コード一覧\n")

	header := "| code | name |\n| --- | --- |\n"
	file.WriteString(header)
	for _, code := range keys {
		file.WriteString("| " + code + " | " + codes[code] + " |\n")
	}

}
