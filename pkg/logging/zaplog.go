package logging

import (
	"fmt"
	"oamp/global"
	"oamp/pkg/setting"
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//初始化日志,传递日志级别参数
func InitLog(level string) {
	var (
		logLevel = zap.InfoLevel //定义出初始的logLevel为info
	)
	//通过switch,根据传入的level不同,给logLevel赋予不同的值
	switch level {
	case "debug":
		logLevel = zap.DebugLevel //debug
	case "info":
		logLevel = zap.InfoLevel //info
	case "warn":
		logLevel = zap.WarnLevel //warn
	case "error":
		logLevel = zap.ErrorLevel //error
	case "panic":
		logLevel = zap.PanicLevel //panic
	case "fatal":
		logLevel = zap.FatalLevel //fatal
	default:
		logLevel = zap.InfoLevel //默认就是info
	}
	//日志级别,zap.LevelEnablerFunc(func(lev zapcore.Level) bool 用来划分不同级别的输出,zapcore.Level类型为int8
	debugPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { //debug级别
		return lev == zap.DebugLevel
	})
	infoPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { //info级别
		return lev == zap.InfoLevel
	})
	errorPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { //error级别
		return lev == zap.ErrorLevel
	})
	//定义接口类型切片,zapcore.Core为一个接口
	var cores []zapcore.Core
	//如果日志级别为debug,那么向切片中添加数据
	if logLevel == zap.DebugLevel {
		cores = append(cores, getEncoderCore(fmt.Sprintf("./%s/server_debug.log", setting.LogSetting.Directory), debugPriority))
	} else if logLevel == zap.InfoLevel {
		cores = append(cores, getEncoderCore(fmt.Sprintf("./%s/server_info.log", setting.LogSetting.Directory), infoPriority))
		cores = append(cores, getEncoderCore(fmt.Sprintf("./%s/server_error.log", setting.LogSetting.Directory), errorPriority))
	}
	//通过zap.New创建一个logger,NewTee创建一个Core，将日志条目复制到两个或更多的底层Core中,...表示解压缩切片
	logger := zap.New(zapcore.NewTee(cores[:]...))
	//判断是否显示行号
	if setting.LogSetting.ShowLine {
		logger = logger.WithOptions(zap.AddCaller())
	}
	//将创建后的logger赋值给全局变量global.Log
	global.Log = logger
}

//函数参数为上面传递过来的文件名和级别,返回为接口类型
func getEncoderCore(fileName string, level zapcore.LevelEnabler) (core zapcore.Core) {
	writer := GetWriteSyncer(fileName) //调用函数对日志进行切割,并返回写入位置
	/*
		     zapcore.NewCore方法中三个参数:
			 getEncoder(): 以什么格式写入日志
			 writer: 日志写到哪里
			 level:   什么级别的日志可以被写入
	*/
	return zapcore.NewCore(getEncoder(), writer, level)
}

//定义日志输出格式为json还是普通格式
func getEncoder() zapcore.Encoder {
	//获取配置文件的输出格式,json or console
	if setting.LogSetting.Format == "json" {
		return zapcore.NewJSONEncoder(getEncoderConfig())
	}
	return zapcore.NewConsoleEncoder(getEncoderConfig())
}

//自定义日志格式,zapcore.EncoderConfig为结构体
func getEncoderConfig() (config zapcore.EncoderConfig) {
	config = zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  setting.LogSetting.StacktraceKey, //栈名
		LineEnding:     zapcore.DefaultLineEnding,        //默认的结尾\n
		EncodeLevel:    zapcore.LowercaseLevelEncoder,    //小写字母输出
		EncodeTime:     CustomTimeEncoder,                //自定义时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,   //编码间隔
		EncodeCaller:   zapcore.FullCallerEncoder,        //控制打印的文件位置是绝对路径,ShortCallerEncoder 是相对路径
	}
	//根据配置文件重新配置编码颜色和字体
	switch {
	case setting.LogSetting.EncodeLevel == "LowercaseLevelEncoder": // 小写编码器(默认)
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	case setting.LogSetting.EncodeLevel == "LowercaseColorLevelEncoder": // 小写编码器带颜色
		config.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	case setting.LogSetting.EncodeLevel == "CapitalLevelEncoder": // 大写编码器
		config.EncodeLevel = zapcore.CapitalLevelEncoder
	case setting.LogSetting.EncodeLevel == "CapitalColorLevelEncoder": // 大写编码器带颜色
		config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	default:
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	}
	return config
}

// 自定义日志输出时间格式
func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

//日志切割,并返回写入位置
func GetWriteSyncer(file string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   file,                          // 日志文件的位置
		MaxSize:    setting.LogSetting.MaxSize,    // 在进行切割之前，日志文件的最大大小（以MB为单位）
		MaxBackups: setting.LogSetting.MaxBackups, // 保留旧文件的最大个数
		MaxAge:     setting.LogSetting.MaxAge,     // 保留旧文件的最大天数
		Compress:   setting.LogSetting.Compress,   // 是否压缩/归档旧文件
	}
	//如果setting.LogSetting.LogToConsole为true,那么日志既输入到控制台又输出到文件
	if setting.LogSetting.LogToConsole {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberJackLogger))
	}
	return zapcore.AddSync(lumberJackLogger)
}
