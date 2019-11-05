package log

import (
	"fmt"
	"github.com/kwuyoucloud/spider/pkg/file"
	"github.com/kwuyoucloud/spider/pkg/file/xml"
	"github.com/kwuyoucloud/spider/pkg/timer"
	"log"
	"os"
	"time"
)

var (
	// IntoFileFlag is a switch, whether log content will write into file.
	// If InfoFileFlag is true, so all of log content will write into file.
	// If InfoFileFlag is false, all of log content will not write into file. But, if other switch blew is true, the log
	// will show in command line.
	intoFileFlag bool
	fatalLogFlag bool
	infoLogFlag  bool
	debugLogFlag bool
	errorLogFlag bool
	// CommandLogFlag is a bool value which means whether type log in command line.
	CommandLogFlag bool
	logDir         = "log"
)

// Init switches from logconfigure xml file.
func init() {
	CommandLogFlag = false
	var xmlconfig *xml.LogConfigure = xml.GetLogConfigureFromXML()

	intoFileFlag = xmlconfig.FileSwitch
	fatalLogFlag = xmlconfig.FatalLog
	infoLogFlag = xmlconfig.InfoLog
	debugLogFlag = xmlconfig.DebugLog
	errorLogFlag = xmlconfig.ErrorLog

	// If the switch which write log into file is true,
	// Create log file folder first.
	if intoFileFlag {
		os.Mkdir("log", 0755)
	}
}

// FatalLog will check fatalswitch status, if is true, print the log text and write text into log file.
func FatalLog(v ...interface{}) {
	if fatalLogFlag {
		if CommandLogFlag {
			log.Println(" - [Fatal] ", v)
		}
		if intoFileFlag {
			logHeader := " - [Fatal] "
			logContent := fmt.Sprintln(v...)
			//Write the log content into file.
			WriteLog(logHeader + logContent)
		}
	}
}

// InfoLog will check infologswitch status, if is true, print the log text and write text into log file.
func InfoLog(v ...interface{}) {
	if infoLogFlag {
		if CommandLogFlag {
			log.Println(" - [Info] ", v)
		}
		if intoFileFlag {
			logHeader := " - [Info] "
			logContent := fmt.Sprintln(v...)
			//Write the log content into file.
			WriteLog(logHeader + logContent)
		}
	}
}

// DebugLog will check debuglogswitch status, if is true, print the log text and write text into log file.
func DebugLog(v ...interface{}) {
	if debugLogFlag {
		if CommandLogFlag {
			log.Println(" - [Debug] ", v)
		}
		if intoFileFlag {
			logHeader := " - [Debug] "
			logContent := fmt.Sprintln(v...)
			//Write the log content into file.
			WriteLog(logHeader + logContent)
		}
	}
}

// ErrLog will check errorlogswitch status, if is true, print the log text and write text into log file.
func ErrLog(v ...interface{}) {
	if errorLogFlag {
		if CommandLogFlag {
			log.Println(" - [Error] ", v)
		}
		if intoFileFlag {
			logHeader := " - [Error] "
			logContent := fmt.Sprintln(v...)
			//Write the log content into file.
			WriteLog(logHeader + logContent)
		}
	}
}

// Fatalln will check fatalinswitch status, if is true, print the log text and write text into log file.
func Fatalln(v ...interface{}) {
	if CommandLogFlag {
		log.Fatalln(" - [Fatal] ", v)
	}
	if intoFileFlag {
		logHeader := " - [Fatal] "
		logContent := fmt.Sprintln(v...)
		//Write the log content into file.
		WriteLog(logHeader + logContent)
	}
}

// Get logpath/year-month-date.log file name, which is today.log file.
func getTodayLogFileName() string {
	today := timer.GetTodayString()

	return logDir + "/" + today + ".log"
}

// WriteLog is a function to write log context into log file.
func WriteLog(logContent string) error {
	// Get log file name path string.
	logName := getTodayLogFileName()
	err := file.WriteContentIntoFile(logName, []byte(time.Now().Local().String()[0:19]+logContent))
	if err != nil {
		fmt.Println("Failed to open the file", err.Error())
		return err
	}
	return nil
}

// CheckErr will check error is nil or not, if error is not nil, write errlog,
// if there is no error, write infolog.
func CheckErr(err error, noerrstr string) {
	if err != nil {
		ErrLog(err)
	} else {
		InfoLog(noerrstr)
	}

}
