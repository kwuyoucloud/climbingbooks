package timer

import (
	"fmt"
	"strings"
	"time"
)

// PrintSthEveryOneSecond is not usful function in real , just a thought to do something need to do every second.
func PrintSthEveryOneSecond(printstr string) {
	t := time.NewTimer(time.Second * 1)
	<-t.C
	fmt.Println(printstr)
	PrintSthEveryOneSecond(printstr)
}

// GetTime will return the string of now.
func GetTime() string {
	return time.Now().String()
}

// GetSimpleTime will return string of now, just like "2019-10-15 10:35:23"
func GetSimpleTime() string {
	longNow := GetTime()
	return strings.Split(longNow, ".")[0]
}

// GetFormateSimpleTime will return string of now, just like "2019-10-15T10:35:23Z"
func GetFormateSimpleTime() string {
	tStr := GetSimpleTime()
	tStr = strings.ReplaceAll(tStr, " ", "T")
	return tStr + "Z"
}

// GetTodayString will return the string of today, in formate like '2019-05-13'.
// Get today's string, like "2019-05-06".
func GetTodayString() string {
	return time.Now().String()[0:10]
}

// GetTodayShortString get today string like "20190506"
func GetTodayShortString() string {
	return strings.ReplaceAll(GetTodayString(), "-", "")
}

// GetTheDayAfterTodayShortString return a short day string after today.
func GetTheDayAfterTodayShortString(daynum int) string {
	today := time.Now()
	wantday := today.AddDate(0, 0, daynum)
	wantdayString := wantday.String()[0:10]
	return strings.ReplaceAll(wantdayString, "-", "")
}
