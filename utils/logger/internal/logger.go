package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"tradeengine/cfg"

	"github.com/magiconair/properties"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

/*
Global Variables
*/
const Config_Update_Period = 5 * time.Second

/*
Define Structure
*/
type MyLogger struct {
	logger        *zap.SugaredLogger
	fileLoggerPtr *lumberjack.Logger
	atomicLevel   zap.AtomicLevel
	propertyPath  string
	property      MyLoggerProperty
}
type MyLoggerProperty struct {
	level               string
	isLogToStd          bool
	enableStdColor      bool
	isLogToFile         bool
	filename            string
	maxFileSize         int
	maxBackupFileNumber int
	isLogColorCustom    bool
	logColor            map[string]int
}

/*
constructors
*/
func NewMyLogger(propertyPath string) MyLogger {
	return MyLogger{
		propertyPath: propertyPath,
	}
}

/*
Define Export Functions
*/
func (l *MyLogger) Debug(template string, a ...interface{}) {
	if l.logger == nil {
		l.init()
	}
	l.logger.Debug(fmt.Sprintf(template, a...))
}
func (l *MyLogger) Info(template string, a ...interface{}) {
	if l.logger == nil {
		l.init()
	}
	l.logger.Info(fmt.Sprintf(template, a...))
}
func (l *MyLogger) Warn(template string, a ...interface{}) {
	if l.logger == nil {
		l.init()
	}
	l.logger.Warn(fmt.Sprintf(template, a...))
}
func (l *MyLogger) Error(template string, a ...interface{}) {
	if l.logger == nil {
		l.init()
	}
	l.logger.Error(fmt.Sprintf(template, a...))
}

// ========================
func (l *MyLogger) init() {
	l.atomicLevel = zap.NewAtomicLevel()
	l.createLogger()
	l.createUpdateThread()
}
func (l *MyLogger) reload() {
	l.createLogger()
}
func (l *MyLogger) clear() {
	l.logger.Sync()
	if l.fileLoggerPtr != nil {
		l.fileLoggerPtr.Close()
	}
	l.logger = nil
}
func (l *MyLogger) createLogger() {
	// Read Config
	var config = loadProperty(l.propertyPath)
	l.property = config

	// Set Logger-Level
	logLevel := parseLevel(config.level)
	l.atomicLevel.SetLevel(logLevel)

	// Create Write-Syncer
	var writeSyncerArr []zapcore.WriteSyncer
	if config.isLogToStd && !config.enableStdColor {
		writeToStd := zapcore.AddSync(os.Stdout)
		writeSyncerArr = append(writeSyncerArr, writeToStd)
	}
	if config.isLogToFile {
		fileLogger := &lumberjack.Logger{
			Filename:   config.filename,
			MaxSize:    config.maxFileSize,
			MaxBackups: config.maxBackupFileNumber,
		}
		writeToFile := zapcore.AddSync(fileLogger)
		writeSyncerArr = append(writeSyncerArr, writeToFile)
		l.fileLoggerPtr = fileLogger
	}

	// Create Encoder
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:       "T",
		LevelKey:      "L",
		NameKey:       "N",
		CallerKey:     "C",
		MessageKey:    "M",
		StacktraceKey: "S",

		EncodeTime:   encodeFmtTime,
		EncodeLevel:  encodeFmtLevel,
		EncodeCaller: encodeFmtCaller,
	}

	var core zapcore.Core
	if config.isLogToStd && config.enableStdColor {
		consoleCore := newColorConsoleCore(l)
		fileCore := zapcore.NewCore(
			NewTradeEngineLogEncoder(encoderCfg),
			zapcore.NewMultiWriteSyncer(writeSyncerArr...),
			l.atomicLevel,
		)
		core = zapcore.NewTee(consoleCore, fileCore)
	} else {
		core = zapcore.NewCore(
			NewTradeEngineLogEncoder(encoderCfg),
			zapcore.NewMultiWriteSyncer(writeSyncerArr...),
			l.atomicLevel,
		)
	}

	l.logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)).Sugar()
}

func (l *MyLogger) createUpdateThread() {
	go func() {
		for {
			time.Sleep(Config_Update_Period)
			l.update()
		}
	}()
}

func (l *MyLogger) update() {
	var newConfig = loadProperty(l.propertyPath)
	ret := l.property.checkDifference(&newConfig)

	if ret == onlyLevelDiff {
		newLevel := parseLevel(newConfig.level)
		l.atomicLevel.SetLevel(newLevel)
	} else if ret == allDiff {
		l.clear()
		l.reload()
	}
}

// ========================
const (
	sameProperty int = iota
	onlyLevelDiff
	allDiff
)

func (p1 *MyLoggerProperty) checkDifference(p2 *MyLoggerProperty) (result int) {
	switch {
	case (p1.isLogToFile != p2.isLogToFile),
		(p1.isLogToStd != p2.isLogToStd),
		(p1.filename != p2.filename),
		(p1.maxFileSize != p2.maxFileSize),
		(p1.maxBackupFileNumber != p2.maxBackupFileNumber):
		return allDiff
	case (p1.level != p2.level):
		return onlyLevelDiff
	}
	return sameProperty
}

// ========================
func loadProperty(propertyPath string) MyLoggerProperty {
	var config = MyLoggerProperty{
		level:               "FATAL",
		filename:            "log/default.log",
		maxFileSize:         50,
		maxBackupFileNumber: 1,
	}

	path := filepath.Join(cfg.LogDir, propertyPath)
	p, err := properties.LoadFile(path, properties.UTF8)
	if err != nil {
		return config
	}

	config.level = p.GetString("log.level", "FATAL")
	config.isLogToStd = p.GetBool("log.outToStd", false)
	config.enableStdColor = p.GetBool("log.outStd.color.enable", false)
	config.isLogToFile = p.GetBool("log.outToFile", false)
	config.filename = filepath.Join(cfg.RootPath, p.GetString("log.file.filepath", "log/default.log"))
	config.maxFileSize = p.GetInt("log.file.maxFileSize", 100)
	config.maxBackupFileNumber = p.GetInt("log.file.maxBackupNumber", 1)
	config.isLogColorCustom = p.GetBool("log.outStd.color.custom", false)
	config.logColor = map[string]int{}
	for _, l := range []string{"debug", "info", "warning", "error", "dpanic", "panic", "fatal"} {
		for _, f := range []string{"fg", "bg"} {
			s := p.GetString("log.outStd.color.custom."+l+"."+f, "0x000000FF")
			n, _ := strconv.ParseInt(s, 16, 32)
			config.logColor[l+"_"+f] = int(n)
		}
	}

	return config
}
func parseLevel(levelTxt string) zapcore.Level {
	switch levelTxt {
	case "DEBUG":
		return zap.DebugLevel
	case "INFO":
		return zap.InfoLevel
	case "WARN":
		return zap.WarnLevel
	case "ERROR":
		return zap.ErrorLevel
	case "FATAL":
		return zapcore.FatalLevel
	}
	return zapcore.FatalLevel
}
