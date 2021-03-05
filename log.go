package nuwa

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

type logImp struct {
}

var instanceLog *logImp

func Loger() *logImp {
	if instanceLog == nil {
		instanceLog = &logImp{} // not thread safe
	}
	return instanceLog
}

func (l *logImp) Out(args ...interface{}) {
	log.Println(args...)
}

func StartLogToFile(dir string) {
	log.SetOutput(newLogWriter(dir))
}

type logWriter struct {
	path string
}

func newLogWriter(path string) *logWriter {
	_, notok := os.Stat(path)
	if notok != nil {
		err := os.MkdirAll(path, os.ModeDir)
		if err != nil {
			println(err.Error())
		}
	}

	return &logWriter{
		path: path,
	}
}

func (that *logWriter) Write(p []byte) (n int, err error) {
	fmt.Print(string(p))
	that.WriteFile(fmt.Sprintf(that.path+"/%s.log", time.Now().Format("2006_01_02_15")), p, 0644)
	return len(p), nil
}

func (that *logWriter) WriteFile(filename string, bytes []byte, perm os.FileMode) {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, perm)
	if err != nil {
		return
	}
	n, err := f.Write(bytes)
	if err == nil && n < len(bytes) {
		err = io.ErrShortWrite
	}
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return
}
