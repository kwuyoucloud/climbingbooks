package file

import (
	"fmt"
	"io/ioutil"
	"os"
)

var fileName string

// ReadFile can read file and return strings of whole file characters.
// The "filePaht" string include path and filename.
func ReadFile(filePath string) ([]byte, error) {
	f, err := os.Open(filePath)
	defer f.Close()

	if err != nil {
		fmt.Println("read file failuer:", err)
		return nil, err

	}

	return ioutil.ReadAll(f)
}
