package log

import (
	"fmt"
	"testing"
)

func TestLog(t *testing.T) {
	fmt.Println("test log")
	t.Error("test error")
}
