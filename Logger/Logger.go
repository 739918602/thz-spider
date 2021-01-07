package Logger

import (
	"io"
	"log"
	"os"
)

var (
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func init() {
	os.MkdirAll("/data/service_logs/", 0655)
	infoFile, ierr := os.OpenFile("/data/service_logs/info.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0655)
	warnFile, wrr := os.OpenFile("/data/service_logs/warn.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0655)
	errFile, eerr := os.OpenFile("/data/service_logs/errors.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0655)

	if ierr != nil || wrr != nil || eerr != nil {
		log.Fatalln("打开日志文件失败：", ierr, wrr, eerr)
	}

	Info = log.New(io.MultiWriter(os.Stdout, infoFile), "Info:", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(io.MultiWriter(os.Stdout, warnFile), "Warning:", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(io.MultiWriter(os.Stderr, errFile), "Error:", log.Ldate|log.Ltime|log.Lshortfile)

}
