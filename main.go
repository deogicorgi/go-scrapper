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

	var jobs []extractJob

	totalPages := getPages()

	for i := 0; i < totalPages; i++ {
		extratedJobs := getPage(i)
		jobs = append(jobs, extratedJobs...)
	}

	writeJobs(jobs)

	duration := time.Since(start)

	log.Printf("Job scrapper total execute time : %s", duration)

}

func writeJobs(jobs []extractJob) {
	file, err := os.Create("jobs.csv")
	checkErr(err)

	w := csv.NewWriter(file)

	headers := []string{"Title", "companyName", "address"}

	wErr := w.Write(headers)
	checkErr(wErr)

	for _, job := range jobs {
		jwErr := w.Write([]string{job.title, job.companyName, job.address})
		checkErr(jwErr)
	}

	defer w.Flush()
}

func getPage(page int) []extractJob {
	start := time.Now()
	var jobs []extractJob
	pageUrl := baseUrl + "&start=" + strconv.Itoa(page*50)
	log.Println("Request URL : " + pageUrl)

	res, err := http.Get(pageUrl)
	checkErr(err)
	checkStatusCode(res)

	doc, err := goquery.NewDocumentFromReader(res.Body)

	doc.Find(".jobCard_mainContent").Each(func(i int, card *goquery.Selection) {
		job := extractJobs(card)
		jobs = append(jobs, job)
	})

	defer res.Body.Close()
	duration := time.Since(start)

	log.Printf("Response and extract job completion time. : %s", duration)
	return jobs
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

func extractJobs(card *goquery.Selection) extractJob {
	jobTitle := clearString(card.Find(".jobTitle").Text())
	companyName := clearString(card.Find(".companyName").Text())
	address := clearString(card.Find(".companyLocation").Text())

	return extractJob{title: jobTitle, address: address, companyName: companyName}
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
