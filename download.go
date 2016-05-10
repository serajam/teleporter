package main

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"os"
	"github.com/serajam/teleporter/utils"
	"fmt"
)

func main() {
	prepareDir()

	mix := utils.Mix{}
	mix.UrlOriginal = os.Args[1]

	doc := getDoc(mix.UrlOriginal)
	parseData(doc, &mix)

	mix.MixPrint()
	checkFileName(&mix)

	utils.Download(&mix)
}

func checkFileName(m *utils.Mix) {
	_, err := os.Stat(fmt.Sprintf("%s%s.mp3", utils.DownloadDir, m.FileName()))
	if err == nil {
		log.Fatal(fmt.Sprintf("File with name \"%s\" exist", m.FileName()))
	}
}

func parseData(doc *goquery.Document, m *utils.Mix) {
	title := doc.Find("h1.cloudcast-title").Text()
	details, isExists := doc.Find("span.play-button").Attr("m-url")
	if isExists == false {
		log.Fatal("Source not found")
	}

	m.Title = title
	m.Details = details
}

func prepareDir() {
	_, err := os.Stat(utils.DownloadDir)
	if err != nil {
		err = os.Mkdir(utils.DownloadDir, 0700)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func getDoc(url string) *goquery.Document {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}
	return doc
}
