package handlehtml

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/kwuyoucloud/spider/internal/app"
	"github.com/kwuyoucloud/spider/pkg/log"
	"testing"
)

var domainURL = "https://book.douban.com/subject/27169753/"

func TestHandleHtml(t *testing.T) {
	// Instantiate default collector
	c := colly.NewCollector()
	// action like chrome
	c.UserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Ubuntu Chromium/77.0.3865.90 Chrome/77.0.3865.90 Safari/537.36"

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Host", domainURL)
		r.Headers.Set("Connection", "keep-alive")
		r.Headers.Set("Accept", "*/*")
		r.Headers.Set("Origin", domainURL)
		// r.Headers.Set("Referer", "http://www.sse.com.cn/assortment/stock/list/share/") //关键头 如果没有 则返回 错误
		r.Headers.Set("Accept-Encoding", "gzip, deflate")
		r.Headers.Set("Accept-Language", "zh-CN,zh;q=0.9")
		log.DebugLog("Visiting", r.URL.String())
	})

	c.OnResponse(func(r *colly.Response) {
		log.DebugLog("Visit url: ", r.Request.URL)
		// Deal the response html.
		// TODO:
		bookInfo, err := GetBookInformation(r.Body)
		if err != nil {
			log.ErrLog("get book information with an error: ", err)
		}
		app.InsertBookInfo(bookInfo)

		fmt.Println((*bookInfo).ToString())
	})

	log.DebugLog("Start.")
	// Start scraping
	c.Visit(domainURL)

	c.OnScraped(func(r *colly.Response) {
		log.DebugLog("Finished", r.Request.URL)
	})
}
