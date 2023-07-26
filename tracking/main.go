package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
)

const source = "api_example_cli"

var (
	server  string
	apiKey  string
	path    string
	archive string
)

func init() {
	flag.StringVar(&server, "server", "https://go.gruzi.ru", "server of gruzi.ru API")
	flag.StringVar(&apiKey, "key", "", "API key")
	flag.StringVar(&path, "path", "./tracking/template.xlsx", "path to excel files with tracking orders")
	flag.StringVar(&archive, "archive", "", "path to archive")
}

func main() {
	flag.Parse()
	if apiKey == "" {
		log.Println("API key is required")
		os.Exit(-1)
	}

	reader, err := NewExcelReader(path)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	data, err := reader.GetData()
	if err != nil {
		log.Println(err)
		os.Exit(2)
	}

	orders := make([]TrackingOrder, 0, len(data))
	for _, item := range data {
		orders = append(orders, item.ToOrder(source))
	}

	err = NewClient(server, apiKey).PostTrackingOrders(orders)
	if err != nil {
		log.Println(err)
		os.Exit(3)
	}

	if archive != "" {
		err = os.Rename(path, filepath.Join(archive, filepath.Base(path)))
		if err != nil {
			log.Println(err)
			os.Exit(4)
		}
	}

	log.Println("Done!")
}
