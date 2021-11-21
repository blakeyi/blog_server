package log

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

type MyFormatter struct{}

func (m *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	var newLog string

	//HasCaller()为true才会有调用信息
	if entry.HasCaller() {
		fName := filepath.Base(entry.Caller.File)
		newLog = fmt.Sprintf("[%s] [%s] [%s:%d %s] %s\n",
			timestamp, entry.Level, fName, entry.Caller.Line, entry.Caller.Function, entry.Message)
	} else {
		newLog = fmt.Sprintf("[%s] [%s] %s\n", timestamp, entry.Level, entry.Message)
	}

	b.WriteString(newLog)
	return b.Bytes(), nil
}

func Init()  {
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&MyFormatter{})
	file, err := os.OpenFile("httpsvr.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatal("log init failed,", err)
		return
	}
	logrus.SetOutput(file)
}
