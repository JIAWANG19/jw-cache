package logger

import (
	"jw-cache/src/pgk/log"
	"os"
	"testing"
)

func init() {
	os.Chdir("D:\\APPDatas\\GoProject\\JWCache")
}

func TestLog(t *testing.T) {
	//currentDir, err := os.Getwd()
	//if err != nil {
	//	panic(err)
	//}
	//println(currentDir)
	//err = os.Chdir("D:\\APPDatas\\GoProject\\JWCache")
	//if err != nil {
	//	panic(err)
	//}
	//
	//currentDir, err = os.Getwd()
	//if err != nil {
	//	panic(err)
	//}
	//println(currentDir)

	log.Info("12312")
}
