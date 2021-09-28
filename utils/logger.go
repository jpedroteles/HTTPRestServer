package utils

import (
	"log"
	"os"
)

var(
	//AppInfoLogger logs information related with app
	AppInfoLogger    *log.Logger
	//AppWarningLogger logs warnings related with app
	AppWarningLogger *log.Logger
	//AppErrorLogger logs errors related with app
	AppErrorLogger   *log.Logger

	//HTTPInfoLogger logs information related with HTTP request
	HTTPInfoLogger    *log.Logger
	//HTTPWarningLogger logs warnings related with HTTP request
	HTTPWarningLogger *log.Logger
	//HTTPErrorLogger logs errors related with HTTP request
	HTTPErrorLogger   *log.Logger
)

//SetUpLogger sets up two loggers, one for general app logging appLog and htaccessLog for http requests related logs
func SetUpLogger(){
	appLog, err:= os.OpenFile("app.log", os.O_APPEND| os.O_CREATE | os.O_WRONLY,0666)
	if err != nil{
		log.Fatal(err)
	}

	AppInfoLogger = log.New(appLog, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	AppWarningLogger = log.New(appLog, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	AppErrorLogger = log.New(appLog, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	htaccessLog, err:= os.OpenFile("htaccess.log", os.O_APPEND| os.O_CREATE | os.O_WRONLY,0666)
	if err != nil{
		log.Fatal(err)
	}

	HTTPInfoLogger = log.New(htaccessLog, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	HTTPWarningLogger = log.New(htaccessLog, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	HTTPErrorLogger = log.New(htaccessLog, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
