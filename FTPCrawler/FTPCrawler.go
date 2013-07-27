package main

import (
	"fmt"
	"github.com/robfig/config"
	. "github.com/shenwei356/FTPCrawler/FTPCrawlerLib"
	"log"
	"os"
	"runtime"
	"time"
)

var (
	site     *Site
	logger   *log.Logger
	loggerFH *os.File
)

func init() {
	logFile := "FTPCrawler.log"
	loggerFH, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0766)
	if err != nil {
		recover()
		fmt.Println(err)
	}
	logger = log.New(loggerFH, "", log.LstdFlags)

	conf, err := config.ReadDefault("FTPCrawler.conf")
	if err != nil {
		recover()
		logger.Fatalln(err)
	}

	name := "DEFAULT"
	host, _ := conf.String(name, "host")
	port, _ := conf.String(name, "port")
	user, _ := conf.String(name, "user")
	passwd, _ := conf.String(name, "passwd")
	path, _ := conf.String(name, "path")
	charset, _ := conf.String(name, "charset")

	site = NewSite(host, port, user, passwd, path, charset, logger)

	var filterTime time.Time
	crawlToday, _ := conf.Bool(name, "crawlToday")
	if crawlToday {
		today := time.Now()
		filterTime = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 1, time.Local)
	} else {
		year := 0
		month := 0
		day := 0
		hour := 0
		min := 0

		year, _ = conf.Int(name, "year")
		month, _ = conf.Int(name, "month")
		day, _ = conf.Int(name, "day")
		hour, _ = conf.Int(name, "day")
		min, _ = conf.Int(name, "min")

		filterTime = time.Date(year, monthTransform(month), day, hour, min, 0, 1, time.Local)
	}
	site.FilterTime = filterTime
}

func main() {
	runtime.GOMAXPROCS(2)
	defer func() {
		if loggerFH != nil {
			loggerFH.Close()
		}
	}()

	for {
		CrawlFTP(site)
		fmt.Println("Finish a round\n\n\n\n")
		time.Sleep(time.Second * 0)
	}

}

func CrawlFTP(site *Site) {
	logger.Println("===========================================================")

	// login
	err := site.Login()
	HandleErr(err)
	logger.Printf("Login: %s@%s\n", site.User, site.Host)

	defer func() {
		err := site.Logout()
		logger.Println("Logout")
		HandleErr(err)
	}()

	// crawl
	crawlStartTime := time.Now()
	logger.Printf("Crawl %s\n", site.Dec.ConvertString(site.Path))

	err = site.DownloadWhileCrawling()
	HandleErr(err)
	crawlTime := time.Now().Sub(crawlStartTime)
	logger.Printf("Crawl using %s\n", crawlTime)

	site.Queue.Wg.Wait()

	downloadTime := time.Now().Sub(crawlStartTime)
	logger.Printf("Download using %s\n", downloadTime)

	logger.Printf("%d files downloaded.\n", site.Queue.Sum)

}

func HandleErr(err error) {
	if err != nil {
		recover()
		logger.Fatalln(err)
	}
}

func monthTransform(m int) (m2 time.Month) {
	month := map[int]time.Month{
		1:  time.January,
		2:  time.February,
		3:  time.March,
		4:  time.April,
		5:  time.May,
		6:  time.June,
		7:  time.July,
		8:  time.August,
		9:  time.September,
		10: time.October,
		11: time.November,
		12: time.December}
	if m >= 1 && m <= 12 {
		m2 = month[m]
	} else {
		m2 = time.January
	}
	return m2
}
