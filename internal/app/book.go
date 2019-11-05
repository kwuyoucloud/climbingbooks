package app

import (
	"database/sql"
	"fmt"
	mysql "github.com/kwuyoucloud/spider/pkg/database/mysql"
)

// BookInfo store book information
type BookInfo struct {
	ID                 string  `json:"id"`
	Name               string  `json:"name"`
	Author             string  `json:"author"`
	AuthorIntroduction string  `json:"authorintroduction"`
	Press              string  `json:"press"`
	BriefIntroduction  string  `json:"briefintroduction"`
	Translator         string  `json:"translator"`
	PageNum            int     `json:"pagenum"`
	Price              float32 `json:"price"`
	ISBN               string  `json:"isbn"`
	DoubanScore        float32 `json:"score"`
	CommentNum         int     `json:"commentnum"`
}

const (
	// USERNAME is database username
	USERNAME = "root"
	// PASSWORD is database password
	PASSWORD = "sunzwa"
	// NETWORK is the protocal
	NETWORK = "tcp"
	// SERVER is the server ip address
	SERVER = "localhost"
	// PORT is the database service port
	PORT = 3306
	// DATABASE is the database name
	DATABASE = "bookdb"
)

var createBookTableStr = `create table bookinfo( 
	id varchar(20) PRIMARY KEY, 
	name varchar(100), 
	author varchar(100), 
	authorintro varchar(1000), 
	press varchar(100), 
	briefintro varchar(1000), 
	translator varchar(100), 
	pagenum int, 
	price float(5,2) DEFAULT 0.0, 
	isbn varchar(50), 
	doubanscore float(2,2) DEFAULT 0.0, 
	commentnum int 
);`
var dbConnNum = 100

var dbchan = make(chan *sql.DB, dbConnNum)

func init() {
	for i := 0; i < dbConnNum; i++ {
		db, err := mysql.NewDBConnect(USERNAME, PASSWORD, NETWORK, SERVER, PORT, DATABASE)
		if err != nil {
			fmt.Printf("Create database connection[%d] with an error: %s \n", i, err)
		}
		dbchan <- db
	}
}

// InsertBookInfo insert book information into database
func InsertBookInfo(bookInfo *BookInfo) error {
	insertStr := "insert into bookinfo(name, author, authorintro, press," +
		"briefintro, translator, pagenum, price, isbn, doubanscore, commentnum) " +
		" values(?,?,?,?,?,?,?,?,?,?,?)"
	db := <-dbchan
	defer func() {
		dbchan <- db
	}()
	err := mysql.InsertData(db, insertStr, bookInfo.Name, bookInfo.Author, bookInfo.AuthorIntroduction, bookInfo.Press,
		bookInfo.BriefIntroduction, bookInfo.Translator, bookInfo.PageNum, bookInfo.Price, bookInfo.ISBN, bookInfo.DoubanScore,
		bookInfo.CommentNum)
	if err != nil {
		fmt.Println("Insert data into database with an error: ", err)
		return err
	}

	return nil
}

// ToString return bookinfo string.
func (bookinfo *BookInfo) ToString() string {
	return fmt.Sprintf("Book Information{Name: %s, Author: %s, AuthorIntroduction: %s, Press: %s, BriefIntroduction: %s, Translator: %s, PageNum: %d, Price: %f, ISBN: %s, DoubanScore: %f, CommentNum: %d}",
		bookinfo.Name, bookinfo.Author, bookinfo.AuthorIntroduction, bookinfo.Press, bookinfo.BriefIntroduction, bookinfo.Translator, bookinfo.PageNum, bookinfo.Price, bookinfo.ISBN, bookinfo.DoubanScore, bookinfo.CommentNum)
}
