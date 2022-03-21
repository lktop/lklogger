package lklogger


import (
	"fmt"
	colorful "github.com/gookit/color"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

var (
	colorDebug             = colorful.HEX("7EC0EE").Sprintf("[DEBUG]")
	colorInfo              = colorful.Blue.Sprintf("[INFO]")
	colorWarn              = colorful.HEX("B89715").Sprintf("[WARN]")
	colorPanic             = colorful.Red.Sprintf("[PANIC]")
	colorError             = colorful.Red.Sprintf("[ERROR]")
	colorFatal             = colorful.Red.Sprintf("[FATAL]")
	colorDPanic            = colorful.Red.Sprintf("[DPANIC]")
)

type LkLogger struct {
	zapLog *zap.Logger
}

func (l *LkLogger) Debug(s string)  {
	l.zapLog.Debug(s)
}

func (l *LkLogger) Info(s string)  {
	l.zapLog.Info(s)
}

func (l *LkLogger) Warn(s string)  {
	l.zapLog.Warn(s)
}

func (l *LkLogger) Error(s string)  {
	l.zapLog.Error(s)
}

func (l *LkLogger) DebugSf(format string,v ...interface{})  {
	l.zapLog.Debug(fmt.Sprintf(format,v...))
}

func (l *LkLogger) InfoSf(format string,v ...interface{})  {
	l.zapLog.Info(fmt.Sprintf(format,v...))
}

func (l *LkLogger) WarnSf(format string,v ...interface{})  {
	l.zapLog.Warn(fmt.Sprintf(format,v...))
}

func (l *LkLogger) ErrorSf(format string,v ...interface{})  {
	l.zapLog.Error(fmt.Sprintf(format,v...))
}

// CONSOLE OUT

func ConsoleTimeEncode (t time.Time, enc zapcore.PrimitiveArrayEncoder){
	enc.AppendString(colorful.HEX("3FAF56").Sprintf("[" + t.Format("2006-01-02 15:04:05.000") + "]"))
}

func ConsoleLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	switch level {
	case zapcore.DebugLevel:
		enc.AppendString(colorDebug)
	case zapcore.InfoLevel:
		enc.AppendString(colorInfo)
	case zapcore.WarnLevel:
		enc.AppendString(colorWarn)
	case zapcore.ErrorLevel:
		enc.AppendString(colorError)
	case zapcore.DPanicLevel:
		enc.AppendString(colorDPanic)
	case zapcore.PanicLevel:
		enc.AppendString(colorPanic)
	case zapcore.FatalLevel:
		enc.AppendString(colorFatal)
	default:
		enc.AppendString(colorful.HEX("FFFF00").Sprintf("[%d]", level))
	}
}

func ConsoleCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("- " + caller.TrimmedPath() + " -")
}

func GetConsoleEncoder() zapcore.EncoderConfig{
	ConsoleEncoderConfig := zap.NewProductionEncoderConfig()	//NewJSONEncoder()输出json格式，NewConsoleEncoder()输出普通文本格式
	ConsoleEncoderConfig.EncodeTime = ConsoleTimeEncode			//指定时间格式
	ConsoleEncoderConfig.EncodeLevel = ConsoleLevelEncoder		//按级别显示不同颜色，不需要的话取值zapcore.CapitalLevelEncoder就可以了
	ConsoleEncoderConfig.EncodeCaller = ConsoleCallerEncoder    //显示完整文件路径
	ConsoleEncoderConfig.ConsoleSeparator = " "
	return ConsoleEncoderConfig
}

// FILE OUT

func FileTimeEncode (t time.Time, enc zapcore.PrimitiveArrayEncoder){
	enc.AppendString("[" + t.Format("2006-01-02 15:04:05.000") + "]")
}

func FileLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	switch level {
	case zapcore.DebugLevel:
		enc.AppendString("[DEBUG]")
	case zapcore.InfoLevel:
		enc.AppendString("[INFO]")
	case zapcore.WarnLevel:
		enc.AppendString("[WARN]")
	case zapcore.ErrorLevel:
		enc.AppendString("[ERROR]")
	case zapcore.DPanicLevel:
		enc.AppendString("[DPANIC]")
	case zapcore.PanicLevel:
		enc.AppendString("[PANIC]")
	case zapcore.FatalLevel:
		enc.AppendString("[FATAL]")
	default:
		enc.AppendString(fmt.Sprintf("[%d]", level))
	}
}

func FileCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("- " + caller.TrimmedPath() + " -")
}

func GetFileEncoder() zapcore.EncoderConfig{
	ConsoleEncoderConfig := zap.NewProductionEncoderConfig()	//NewJSONEncoder()输出json格式，NewConsoleEncoder()输出普通文本格式
	ConsoleEncoderConfig.EncodeTime = FileTimeEncode			//指定时间格式
	ConsoleEncoderConfig.EncodeLevel = FileLevelEncoder		//按级别显示不同颜色，不需要的话取值zapcore.CapitalLevelEncoder就可以了
	ConsoleEncoderConfig.EncodeCaller = FileCallerEncoder    //显示完整文件路径
	ConsoleEncoderConfig.ConsoleSeparator = " "
	return ConsoleEncoderConfig
}

func NewLKLogger(callerPath bool,StackTrace bool)*LkLogger {
	var coreArr []zapcore.Core

	//获取编码器
	ConosleEncoder := zapcore.NewConsoleEncoder(GetConsoleEncoder())
	FileEncoder := zapcore.NewConsoleEncoder(GetFileEncoder())

	//日志级别
	highPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool{	//error级别
		return lev >= zap.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {	//info和debug级别,debug级别是最低的
		return lev < zap.ErrorLevel && lev >= zap.DebugLevel
	})

	//info文件writeSyncer
	infoFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./log/info.log", 	//日志文件存放目录，如果文件夹不存在会自动创建
		MaxSize:    2,            		//文件大小限制,单位MB
		MaxBackups: 100,            	//最大保留日志文件数量
		MaxAge:     30,           		//日志文件保留天数
		Compress:   false,        		//是否压缩处理
	})

	InfoConsoleCore := zapcore.NewCore(ConosleEncoder, zapcore.AddSync(os.Stdout), lowPriority) 	//第三个及之后的参数为写入文件的日志级别,ErrorLevel模式只记录error级别的日志
	InfoFileCore := zapcore.NewCore(FileEncoder, infoFileWriteSyncer, lowPriority) 					//第三个及之后的参数为写入文件的日志级别,ErrorLevel模式只记录error级别的日志
	coreArr = append(coreArr, InfoConsoleCore)
	coreArr = append(coreArr, InfoFileCore)

	errorFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./log/error.log", 		//日志文件存放目录
		MaxSize:    1,            			//文件大小限制,单位MB
		MaxBackups: 5,            			//最大保留日志文件数量
		MaxAge:     30,           			//日志文件保留天数
		Compress:   false,        			//是否压缩处理
	})
	errorConsoleCore := zapcore.NewCore(ConosleEncoder, zapcore.AddSync(os.Stdout), highPriority) //第三个及之后的参数为写入文件的日志级别,ErrorLevel模式只记录error级别的日志
	errorFileCore := zapcore.NewCore(FileEncoder, errorFileWriteSyncer, highPriority) //第三个及之后的参数为写入文件的日志级别,ErrorLevel模式只记录error级别的日志
	coreArr = append(coreArr, errorConsoleCore)
	coreArr = append(coreArr, errorFileCore)

	var zapLog *zap.Logger
	var lZapOption []zap.Option
	if callerPath{
		lZapOption = append(lZapOption, zap.AddCaller())  //zap.AddCaller()为显示文件名和行号，可省略
	}
	if StackTrace{
		lZapOption = append(lZapOption, zap.AddStacktrace(zapcore.ErrorLevel))  //zap.AddStacktrace()为显示调用堆栈
	}
	zapLog = zap.New(zapcore.NewTee(coreArr...),lZapOption...)
	return &LkLogger{zapLog: zapLog}
}

func NewLKLoggerAll(callerPath bool,StackTrace bool)*LkLogger {
	var coreArr []zapcore.Core

	//获取编码器
	ConosleEncoder := zapcore.NewConsoleEncoder(GetConsoleEncoder())
	FileEncoder := zapcore.NewConsoleEncoder(GetFileEncoder())

	//日志级别
	highPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool{	//error级别
		return lev >= zap.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {	//info和debug级别,debug级别是最低的
		return lev < zap.ErrorLevel && lev >= zap.DebugLevel
	})

	//文件writeSyncer
	FileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./log/all_log.log", 	//日志文件存放目录，如果文件夹不存在会自动创建
		MaxSize:    2,            		//文件大小限制,单位MB
		MaxBackups: 100,            	//最大保留日志文件数量
		MaxAge:     30,           		//日志文件保留天数
		Compress:   false,        		//是否压缩处理
	})

	InfoConsoleCore := zapcore.NewCore(ConosleEncoder, zapcore.AddSync(os.Stdout), lowPriority) 	//第三个及之后的参数为写入文件的日志级别,ErrorLevel模式只记录error级别的日志
	InfoFileCore := zapcore.NewCore(FileEncoder, FileWriteSyncer, lowPriority) 					//第三个及之后的参数为写入文件的日志级别,ErrorLevel模式只记录error级别的日志
	coreArr = append(coreArr, InfoConsoleCore)
	coreArr = append(coreArr, InfoFileCore)


	errorConsoleCore := zapcore.NewCore(ConosleEncoder, zapcore.AddSync(os.Stdout), highPriority) //第三个及之后的参数为写入文件的日志级别,ErrorLevel模式只记录error级别的日志
	errorFileCore := zapcore.NewCore(FileEncoder, FileWriteSyncer, highPriority) //第三个及之后的参数为写入文件的日志级别,ErrorLevel模式只记录error级别的日志
	coreArr = append(coreArr, errorConsoleCore)
	coreArr = append(coreArr, errorFileCore)

	var zapLog *zap.Logger
	var lZapOption []zap.Option
	if callerPath{
		lZapOption = append(lZapOption, zap.AddCaller())  //zap.AddCaller()为显示文件名和行号，可省略
	}
	if StackTrace{
		lZapOption = append(lZapOption, zap.AddStacktrace(zapcore.ErrorLevel))  //zap.AddStacktrace()为显示调用堆栈
	}
	zapLog = zap.New(zapcore.NewTee(coreArr...),lZapOption...)
	return &LkLogger{zapLog: zapLog}
}
