package db

import (
	"fmt"
	"io"
	"my_test/cn"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var L *logrus.Entry = nil

func InitLog() {
	tempL := logrus.New()
	logger := &lumberjack.Logger{
		// 日志输出文件路径 /data/log/kasa-lcs/kasa_lcs.log
		Filename: fmt.Sprintf("%s_%s.log", cn.Schema, cn.Name),
		// 日志文件最大 size, 单位是 MB
		MaxSize: 500, // megabytes
		// 最大过期日志保留的个数
		MaxBackups: 10,
		// 保留过期文件的最大时间间隔,单位是天
		MaxAge: 28, //days
		// 是否需要压缩滚动日志, 使用的 gzip 压缩
		Compress: true, // disabled by default
	}
	tempL.SetOutput(io.MultiWriter(logger, os.Stdout))
	tempL.SetFormatter(&logrus.JSONFormatter{})
	L = tempL.WithFields(logrus.Fields{
		"prdline": cn.PrdLine,
		"app":     cn.Schema + cn.Name,
	})
	L.Infof("%s-%s log start!", cn.Schema, cn.Name)
}
