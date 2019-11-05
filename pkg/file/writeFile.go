package file

import (
	"io"
	"os"
)

// WriteContentIntoFile function.
// Didn't test net.
func WriteContentIntoFile(filename string, data []byte) error {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	defer f.Close()

	if err != nil {
		return err
	}
	n, err := f.Write(data)

	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	if err1 := f.Close(); err == nil {
		err = err1
	}

	return err
}
