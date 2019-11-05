package main

import (
	"fmt"
	// "math/rand"
	// "net/http"
	// "net/url"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly"
	collyproxy "github.com/gocolly/colly/proxy"
	"github.com/kwuyoucloud/spider/internal/app"
	"github.com/kwuyoucloud/spider/pkg/database/mybolt"
	"github.com/kwuyoucloud/spider/pkg/handlehtml"
	"github.com/kwuyoucloud/spider/pkg/log"
	"github.com/kwuyoucloud/spider/pkg/proxy"
)

var visitDBName = "visit.db"
var visitedBucketName = "visitedlink"

var domainURL = "https://book.douban.com/"
var threadNum = 1000

// var domainURL = "http://weibo.com/"
var allowdomains = []string{"book.douban.com"}

var requestchan = make(chan struct{}, threadNum)
var proxies []string

func init() {
	// init thread number.
	for i := 0; i < threadNum-1; i++ {
		requestchan <- struct{}{}
	}

	// init proxies from a list, will achieve a dynamic api type later.
	proxies = proxy.GetProxyIPList()

	// init boltdb.
	_, err := os.Stat(visitDBName)
	if err != nil {
		if os.IsNotExist(err) {
			_, err := os.Create(visitDBName)
			if err != nil {
				log.ErrLog(fmt.Sprintf("Create visited db file[%s] with an error: %s", visitDBName, err.Error()))
			}

		}
	}
	// create bucket in boltdb database
	_ = mybolt.AddKeyValueonBolt(visitDBName, visitedBucketName, "1", "1")
}

func main() {
	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains
		colly.AllowedDomains(allowdomains...),
	)

	// action like chrome
	c.UserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Ubuntu Chromium/77.0.3865.90 Chrome/77.0.3865.90 Safari/537.36"

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Host", domainURL)
		r.Headers.Set("Connection", "keep-alive")
		r.Headers.Set("Accept", "*/*")
		r.Headers.Set("Origin", domainURL)
		r.Headers.Set("Referer", "https://book.douban.com/subject/26933399") //关键头 如果没有 则返回 错误
		r.Headers.Set("Accept-Encoding", "gzip, deflate")
		r.Headers.Set("Accept-Language", "zh-CN,zh;q=0.9")
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnResponse(func(r *colly.Response) {
		go func() {
			urlstr := r.Request.URL.String()
			fmt.Println("Visit url: ", urlstr)
			bookInfo, err := handlehtml.GetBookInformation(r.Body)
			if err != nil {
				fmt.Println("Get bookInfo from url with error: ", err)
			}
			// book name is not ""
			bookName := strings.Trim(bookInfo.Name, " ")
			if bookName != "" {
				err = app.InsertBookInfo(bookInfo)
				log.DebugLog(bookInfo.ToString())
			}

			if err != nil {
				fmt.Println("Insert into database with an error: ", err)
			}

			requestchan <- struct{}{}
			fmt.Println("requestchan <- struct{}{}")
		}()
	})

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		// fmt.Println("OnHtml")
		link := e.Attr("href")

		linktail := strings.TrimLeft(link, "https://book.douban.com/subject/")
		tailtemp := strings.Trim(linktail, "/")

		if strings.Count(tailtemp, "/") == 0 {
			/*
				if strings.HasPrefix(link, "https://book.douban.com/subject") && !strings.HasSuffix(link, "new_review") &&
					!strings.HasSuffix(link, "discussion") && !strings.HasSuffix(link, "doulists") &&
					!strings.HasSuffix(link, "wishes") && !strings.HasSuffix(link, "doing") {
			*/
			visitStr := mybolt.GetValueFromBolt(visitDBName, visitedBucketName, link)
			if visitStr != "" {
				// fmt.Println("This link has visited before, skip this linke: ", link)
				return
			}

			err := mybolt.AddKeyValueonBolt(visitDBName, visitedBucketName, link, "1")
			if err != nil {
				fmt.Println("Add link to visited list failed with an error: ", err)
				fmt.Println("Programma continue.")
			}

			// log.DebugLog(fmt.Sprintf("Link found: %q -> %s\n", e.Text, link))
			fmt.Println(fmt.Sprintf("Link found: %q -> %s\n", e.Text, link))
			// time.Sleep(time.Millisecond * 50)
			select {
			case <-requestchan:
				fmt.Println("<- requestchan")
			case <-time.After(time.Millisecond * 50):
				fmt.Println("<- requestchan timeout")
			}

			// visit url, the codes after c.Visit() will not run anyway.
			c.Visit(e.Request.AbsoluteURL(link))
		}
	})

	// Set proxy
	rp, err := collyproxy.RoundRobinProxySwitcher(proxies...)
	if err != nil {
		c.SetProxyFunc(rp)
	}

	// Error write into log file.
	c.OnError(func(_ *colly.Response, err error) {
		log.ErrLog("There is something wrong: ", err)
	})

	// Start scraping
	c.Visit(domainURL)

	c.OnScraped(func(r *colly.Response) {
		// log.DebugLog("Finished", r.Request.URL)
		fmt.Println("Finished", r.Request.URL)
	})
}
