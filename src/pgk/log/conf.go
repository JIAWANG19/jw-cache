package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"jw-cache/src/pgk/setting"
	"log"
	"os"
	"time"
)

var (
	logger      *logrus.Logger
	logLevel    string
	logFormat   string
	logLevelMap = map[string]logrus.Level{
		"debug": logrus.DebugLevel,
		"info":  logrus.InfoLevel,
		"warn":  logrus.WarnLevel,
		"error": logrus.ErrorLevel,
	}
)

// loadLogConf 加载配置文件
func loadLogConf() {
	logConf, err := setting.Cfg.GetSection("log")
	if err != nil {
		log.Fatalf("Fail to get section 'log': %v", err)
	}
	logLevel = logConf.Key("level").String()
	logFormat = logConf.Key("file_format").String()
}

// getFileName 获取日志文件名的格式
func getFileName() string {
	filename := time.Now().Format(logFormat) + ".log"
	return "log/" + filename
}

type CustomFormatter struct {
}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	level := entry.Level.String()
	message := entry.Message
	return []byte(fmt.Sprintf("[%s] [%s] %s\n", timestamp, level, message)), nil
}

func init() {
	logger = logrus.New()
	logger.SetFormatter(&CustomFormatter{})

	loadLogConf()
	fullFilename := getFileName()

	file, err := os.OpenFile(fullFilename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Fail to open log file %s: %v", fullFilename, err)
	}
	// 日志打印到文件和控制台
	logger.SetOutput(io.MultiWriter(file, os.Stdout))

	if level, ok := logLevelMap[logLevel]; !ok {
		log.Fatalf("bad log level in conf: %s", logLevel)
	} else {
		logger.SetLevel(level)
	}
}

// Debug Debug
func Debug(s string, v ...interface{}) {
	logger.Debugf(s, v...)
}

// Info Info
func Info(s string, v ...interface{}) {
	logger.Infof(s, v...)
}

func Warn(s string, v ...interface{}) {
	logger.Warnf(s, v...)
}

// Error Error
func Error(s string, v ...interface{}) {
	logger.Errorf(s, v...)
}
