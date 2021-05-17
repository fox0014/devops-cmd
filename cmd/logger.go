package cmd

import (
	"log"
	"os"
)

var (
	Info  *log.Logger
	Error *log.Logger
)

func init() {
	//f, err := os.OpenFile("./Log.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	//if err != nil {
	//	log.Fatal("os.OpenFile err", err)
	//}
	//defer f.Close()
	//writers := []io.Writer{
	//	f,
	//	os.Stdout}
	//实例化Writer接口
	//fileAndStdoutWriter := io.MultiWriter(writers...)
	//MultiWriter创建一个Writer接口，会将提供给其的数据写入所有创建时提供的Writer接口
	//Info = log.New(fileAndStdoutWriter, "", log.Ldate|log.Ltime|log.Lshortfile)
	Info = log.New(os.Stdout, "[INFO]", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(os.Stderr, "[ERROR]", log.Ldate|log.Ltime|log.Lshortfile)
}
