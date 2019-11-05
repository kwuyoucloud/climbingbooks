package file

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/satori/go.uuid"
)

var foldDir = "download"

func init() {
	os.Mkdir(foldDir, 0755)
}

// Download function can download file from URL.
// Copied from internet.
// Download file from url
func Download(fileURL string, filename string, fileType string, p chan string) {
	fileName := ""

	if filename == "" {
		uid := uuid.NewV4() //If file name is null, create a random file name.
		fileName = uid.String()
	} else {
		fileName = filename
	}

	fileName = fileName + "." + fileType

	fmt.Println("Download file name is:" + fileName)
	f, err := os.Create(foldDir + "/" + fileName)
	if err != nil {
		fmt.Println("Create file failed.")
	}
	defer f.Close() //Close file before stop function.

	resp, err := http.Get(fileURL)
	if err != nil {
		fmt.Println("http.get error", err)
	}

	body, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		fmt.Println("Read data failed.")
	}
	defer resp.Body.Close() //Close before exit
	f.Write(body)

	p <- fileName //Send file name into channel.
}
