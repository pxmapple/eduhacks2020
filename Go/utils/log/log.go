package log

import (
	"eduhacks2020/Go/pkg/setting"
	rotateLogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
)

// Setup 初始化
func Setup() {
	basePath := getCurrentDirectory()

	writer, err := rotateLogs.New(
		basePath+"/log/info/"+"%Y-%m-%d"+".log",
		rotateLogs.WithLinkName("log.log"), // 生成软链，指向最新日志文件
		//rotateLogs.WithMaxAge(maxAge),      // 文件最大保存时间
	)
	if err != nil {
		logrus.Errorf("config local file system logger error. %+v", errors.WithStack(err))
	}

	errorWriter, err := rotateLogs.New(
		basePath+"/log/error/"+"%Y-%m-%d"+".log",
		rotateLogs.WithLinkName("error.log"), // 生成软链，指向最新日志文件
		//rotateLogs.WithMaxAge(maxAge),        // 文件最大保存时间
	)
	if err != nil {
		logrus.Errorf("config local file system logger error. %+v", errors.WithStack(err))
	}
	lfHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer, // 为不同级别设置不同的输出目的
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: errorWriter,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		PrettyPrint:     false, //是否格式化json格式
		FieldMap: logrus.FieldMap{
			"host": setting.GlobalSetting.LocalHost,
		},
	})
	//logrus.SetFormatter(&logrus.JSONFormatter{
	//	TimestampFormat: "2006-01-02 15:04:05",
	//	PrettyPrint:     false, //是否格式化json格式
	//	FieldMap: logrus.FieldMap{
	//		"host": setting.GlobalSetting.LocalHost,
	//	}})
	//logrus.SetReportCaller(true) //是否记录代码位置
	logrus.AddHook(lfHook)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
	})
	//logrus.SetOutput()
}

//获取当前程序运行的文件夹
func getCurrentDirectory() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return strings.Replace(dir, "\\", "/", -1)
}
