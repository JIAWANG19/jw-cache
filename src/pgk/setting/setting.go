package setting

import (
	"github.com/go-ini/ini"
	"log"
)

var (
	Cfg *ini.File
)

func init() {
	var err error
	Cfg, err = ini.Load("conf/conf.ini")
	if err != nil {
		log.Fatalf("Fail to parse 'conf/conf.ini': %v", err)
	}
}
