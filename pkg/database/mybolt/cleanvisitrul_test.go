package mybolt

import (
	"fmt"
	"strings"
	"testing"
)

var visitDBName = "visit.db"
var visitedBucketName = "visitedlink"

// Test test
func TestCleanVisitURL(t *testing.T) {
	fmt.Println("Start clean visited url")

	keys := GetAllKeys(visitDBName, visitedBucketName)
	headStr := "https://book.douban.com/subject/"

	for i := range keys {
		key := keys[i]
		if !strings.HasPrefix(key, headStr) {
			err := DelKeyonBolt(visitDBName, visitedBucketName, key)
			fmt.Println("Delete key: ", key)
			if err != nil {
				fmt.Printf("Delete key[%s] with an error: %s \n", key, err.Error())
			}

			continue
		}

		tailStr := strings.TrimLeft(key, headStr)
		tailStr = strings.Trim(tailStr, "/")
		if strings.Index(tailStr, "/") != -1 {
			err := DelKeyonBolt(visitDBName, visitedBucketName, key)
			fmt.Println("Delete key: ", key)
			if err != nil {
				fmt.Printf("Delete key[%s] with an error: %s \n", key, err.Error())
			}
		}
	}
}
