package handlehtml

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"github.com/kwuyoucloud/spider/internal/app"
	"strconv"
	"strings"
)

// GetBookInformation return BookInfo from html
func GetBookInformation(htmlBody []byte) (*app.BookInfo, error) {
	bookinfo := &app.BookInfo{}
	bookinfo.PageNum = 0
	bookinfo.Price = 0.0
	bookinfo.DoubanScore = 0.0
	bookinfo.CommentNum = 0

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(htmlBody))
	if err != nil {
		return nil, err
	}

	// Book name
	bookinfo.Name = doc.Find("#wrapper > h1 > span").Text()
	infoDiv := doc.Find("#info")

	infoDiv.Find("span > a").Each(func(i int, s *goquery.Selection) {
		switch i {
		case 0:
			bookinfo.Author = s.Text()
		case 1:
			bookinfo.Translator = s.Text()
		}
	})

	bookinfoStr := infoDiv.Text()
	bookinfoStr = strings.Replace(bookinfoStr, "\n", "|", -1)

	// for replace two landspace to one landspace
	for strings.Index(bookinfoStr, " ") != -1 {
		bookinfoStr = strings.Replace(bookinfoStr, " ", "", -1)
	}
	for strings.Index(bookinfoStr, "||") != -1 {
		bookinfoStr = strings.Replace(bookinfoStr, "||", "|", -1)
	}
	// fmt.Println("bookinfostring: ", bookinfoStr)

	bookinfos := strings.Split(bookinfoStr, "|")
	for _, v := range bookinfos {
		v = strings.Trim(v, " ")
		if strings.HasPrefix(v, "出版社") {
			bookinfo.Press = strings.Split(v, ":")[1]
		}
		if strings.HasPrefix(v, "页数") {
			page, err := strconv.ParseInt(strings.Split(v, ":")[1], 0, 12)
			if err == nil {
				bookinfo.PageNum = int(page)
			}
		}
		if strings.HasPrefix(v, "定价") {
			price, err := strconv.ParseFloat(strings.Split(v, ":")[1], 12)
			if err == nil {
				bookinfo.Price = float32(price)
			}
		}
		if strings.HasPrefix(v, "ISBN") {
			bookinfo.ISBN = strings.Split(v, ":")[1]
		}
	}

	interestsectl := doc.Find("#interest_sectl")
	score := interestsectl.Find("strong.rating_num").Text()
	tempscore, err := strconv.ParseFloat(strings.Trim(score, " "), 12)
	if err == nil {
		bookinfo.DoubanScore = float32(tempscore)
	}

	commDiv := doc.Find(".rating_people > span")
	commNum := commDiv.Text()
	bookinfo.CommentNum, _ = strconv.Atoi(strings.Trim(commNum, " "))

	/*
		linkReport := doc.Find("#link-report")
		bookinfo.BriefIntroduction = linkReport.Find("div.intro").Text()

		indentDiv := doc.Find("div.indent > div > div.intro")
		bookinfo.AuthorIntroduction = indentDiv.Text()
	*/

	return bookinfo, nil
}
