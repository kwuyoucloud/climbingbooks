package mysql

import (
	"database/sql"
	"fmt"
	// just use this to init database driver.
	_ "github.com/go-sql-driver/mysql"
	"github.com/kwuyoucloud/spider/pkg/log"
	"time"
)

// NewDBConnect get a new connection
func NewDBConnect(username, password, network, server string, port int, database string) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", username, password, network, server, port, database)
	DB, err := sql.Open("mysql", dsn)
	if err != nil {
		log.ErrLog("Can not conncet mysql database with an error: ", err)
		return nil, err
	}

	DB.SetConnMaxLifetime(100 * time.Second)
	DB.SetMaxOpenConns(100)
	DB.SetMaxIdleConns(16)

	return DB, nil
}

// InsertData comment
func InsertData(db *sql.DB, statment string, args ...interface{}) error {
	result, err := db.Exec(statment, args...)
	if err != nil {
		log.ErrLog(fmt.Sprintf("InsertData statment: %s,with an error: ", err))
		return err
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		log.ErrLog(fmt.Sprintf("InsertData statment: %s, get LastInsertID with an error: ", err))
		// return err
	}
	fmt.Println("LastInsertID: ", lastInsertID)
	rowsaffected, err := result.RowsAffected()
	if err != nil {
		log.ErrLog(fmt.Sprintf("InsertData statment: %s, get RowsAffected with an error: ", err))
		return err
	}
	fmt.Println("RowsAffected: ", rowsaffected)

	return nil
}
