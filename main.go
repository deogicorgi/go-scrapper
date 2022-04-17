package main

import (
	"encoding/csv"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var baseUrl = "https://kr.indeed.com/jobs?q=python&limit=50"

type extractJob struct {
	title       string
	companyName string
	address     string
}

func main() {
	start := time.Now()

	mc := make(chan []extractJob)

	var jobs []extractJob

	totalPages := getPages()

	for i := 0; i < totalPages; i++ {
		go getPage(i, mc)
	}

	for i := 0; i < totalPages; i++ {
		jobs = append(jobs, <-mc...)
	}

	writeJobs(jobs)

	log.Printf("Job scrapper total execute time : %s", time.Since(start))

}

func writeJobs(jobs []extractJob) {

	wc := make(chan bool)

	file, err := os.Create("jobs.csv")
	checkErr(err)

	w := csv.NewWriter(file)

	headers := []string{"Title", "companyName", "address"}

	wErr := w.Write(headers)
	checkErr(wErr)

	for _, job := range jobs {
		go writeJob(w, []string{job.title, job.companyName, job.address}, wc)
	}

	for i := 0; i < len(jobs); i++ {
		<-wc
	}

	defer w.Flush()
}

func writeJob(w *csv.Writer, job []string, wc chan bool) {
	jwErr := w.Write(job)
	checkErr(jwErr)
	wc <- true
}

func getPage(page int, mc chan<- []extractJob) {
	start := time.Now()
	var jobs []extractJob

	c := make(chan extractJob)

	pageUrl := baseUrl + "&start=" + strconv.Itoa(page*50)
	log.Println("Request URL : " + pageUrl)

	res, err := http.Get(pageUrl)
	checkErr(err)
	checkStatusCode(res)

	doc, err := goquery.NewDocumentFromReader(res.Body)

	cards := doc.Find(".jobCard_mainContent")

	extractStart := time.Now()
	cards.Each(func(i int, card *goquery.Selection) {
		go extractJobs(card, c)
	})
	log.Printf("extract job completion time. : %s", time.Since(extractStart))

	for i := 0; i < cards.Length(); i++ {
		job := <-c
		jobs = append(jobs, job)
	}

	log.Printf("Page process completion time. : %s", time.Since(start))

	mc <- jobs
	defer res.Body.Close()
}

func getPages() int {
	res, err := http.Get(baseUrl)
	pages := 0

	checkErr(err)
	checkStatusCode(res)

	doc, err := goquery.NewDocumentFromReader(res.Body)

	doc.Find(".pagination").Each(func(i int, selection *goquery.Selection) {
		pages = selection.Find("a").Length()
	})

	defer res.Body.Close()
	return pages
}

func extractJobs(card *goquery.Selection, c chan<- extractJob) {
	jobTitle := clearString(card.Find(".jobTitle").Text())
	companyName := clearString(card.Find(".companyName").Text())
	address := clearString(card.Find(".companyLocation").Text())

	c <- extractJob{title: jobTitle, address: address, companyName: companyName}
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkStatusCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln(errors.New("status code error"))
	}
}

func clearString(str string) string {

	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}
