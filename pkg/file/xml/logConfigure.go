package xml

import (
	"encoding/xml"
	"fmt"
	"github.com/kwuyoucloud/spider/pkg/file"
)

// Define the location of configuration file.
var logXMLName = "config/log.config"

// LogConfigure struct,
// Fileswitch is a switch, if value is true, the application will write log into file.
// else, if value is false, the application will not write log into file.
//
// InfoLog, DebugLog, FatalLog, ErrorLog is switches, if the value is true,
// then application will write infolog into file,
// else will not write infolog into file.
type LogConfigure struct {
	XMLName    xml.Name `xml:"logconfig"`
	FileSwitch bool     `xml:"fileswitch"`
	InfoLog    bool     `xml:"infolog"`
	DebugLog   bool     `xml:"debuglog"`
	FatalLog   bool     `xml:"fatallog"`
	ErrorLog   bool     `xml:"errorlog"`
}

func init() {
	logXMLName = "config/log.config"
}

// GetLogConfigureFromXML can get LogConfigure struct from cofiguration file.
func GetLogConfigureFromXML() *LogConfigure {
	configXMLType, err := file.ReadFile(logXMLName)

	if err != nil {
		fmt.Println(err)
	}

	//read log.config file to get status of every log rank
	v := LogConfigure{}
	err = xml.Unmarshal(configXMLType, &v)

	return &v

}
