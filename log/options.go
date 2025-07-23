package log

import (
	"gatesvr/etc"
	"strings"
	"time"
)

const (
	defaultFile              = "./log/due.log"
	defaultLevel             = InfoLevel
	defaultFormat            = TextFormat
	defaultStdout            = true
	defaultFileMaxAge        = 7 * 24 * time.Hour
	defaultFileMaxSize       = 100
	defaultFileCutRule       = CutByDay
	defaultTimeFormat        = "2006/01/02 15:04:05.000000"
	defaultCallerFullPath    = false
	defaultClassifiedStorage = false
)

const (
	defaultFileKey              = "etc.log.file"
	defaultLevelKey             = "etc.log.level"
	defaultFormatKey            = "etc.log.format"
	defaultTimeFormatKey        = "etc.log.timeFormat"
	defaultStackLevelKey        = "etc.log.stackLevel"
	defaultFileMaxAgeKey        = "etc.log.fileMaxAge"
	defaultFileMaxSizeKey       = "etc.log.fileMaxSize"
	defaultFileCutRuleKey       = "etc.log.fileCutRule"
	defaultStdoutKey            = "etc.log.stdout"
	defaultCallerFullPathKey    = "etc.log.callerFullPath"
	defaultClassifiedStorageKey = "etc.log.classifiedStorage"
)

type options struct {
	file              string        // 输出的文件路径，有文件路径才会输出到文件，否则只会输出到终端
	level             Level         // 输出的最低日志级别，默认Info
	format            Format        // 输出的日志格式，Text或者Json，默认Text
	stdout            bool          // 是否输出到终端，debug模式下默认输出到终端
	timeFormat        string        // 时间格式，标准库时间格式，默认2006/01/02 15:04:05.000000
	stackLevel        Level         // 堆栈的最低输出级别，默认不输出堆栈
	fileMaxAge        time.Duration // 文件最大留存时间，默认7天
	fileMaxSize       int64         // 文件最大尺寸限制，单位（MB），默认100MB
	fileCutRule       CutRule       // 文件切割规则，默认按照天
	callerSkip        int           // 调用者跳过的层级深度
	callerFullPath    bool          // 是否启用调用文件全路径，默认短路径
	classifiedStorage bool          // 是否启用分级存储，默认不分级
}

type Option func(o *options)

func defaultOptions() *options {
	opts := &options{
		file:              defaultFile,
		level:             defaultLevel,
		format:            defaultFormat,
		stdout:            defaultStdout,
		timeFormat:        defaultTimeFormat,
		fileMaxAge:        defaultFileMaxAge,
		fileMaxSize:       defaultFileMaxSize,
		fileCutRule:       defaultFileCutRule,
		callerFullPath:    defaultCallerFullPath,
		classifiedStorage: defaultClassifiedStorage,
	}

	file := etc.Get(defaultFileKey).String()
	if file != "" {
		opts.file = file
	}

	level := etc.Get(defaultLevelKey).String()
	if lvl := ParseLevel(level); lvl != NoneLevel {
		opts.level = lvl
	}

	format := etc.Get(defaultFormatKey).String()
	switch strings.ToLower(format) {
	case JsonFormat.String():
		opts.format = JsonFormat
	case TextFormat.String():
		opts.format = TextFormat
	}

	timeFormat := etc.Get(defaultTimeFormatKey).String()
	if timeFormat != "" {
		opts.timeFormat = timeFormat
	}

	stackLevel := etc.Get(defaultStackLevelKey).String()
	if lvl := ParseLevel(stackLevel); lvl != NoneLevel {
		opts.stackLevel = lvl
	}

	fileMaxAge := etc.Get(defaultFileMaxAgeKey).Duration()
	if fileMaxAge > 0 {
		opts.fileMaxAge = fileMaxAge
	}

	fileMaxSize := etc.Get(defaultFileMaxSizeKey).Int64()
	if fileMaxSize > 0 {
		opts.fileMaxSize = fileMaxSize
	}

	fileCutRule := etc.Get(defaultFileCutRuleKey).String()
	switch strings.ToLower(fileCutRule) {
	case CutByYear.String():
		opts.fileCutRule = CutByYear
	case CutByMonth.String():
		opts.fileCutRule = CutByMonth
	case CutByDay.String():
		opts.fileCutRule = CutByDay
	case CutByHour.String():
		opts.fileCutRule = CutByHour
	case CutByMinute.String():
		opts.fileCutRule = CutByMinute
	case CutBySecond.String():
		opts.fileCutRule = CutBySecond
	}

	opts.stdout = etc.Get(defaultStdoutKey, defaultStdout).Bool()
	opts.callerFullPath = etc.Get(defaultCallerFullPathKey, defaultCallerFullPath).Bool()
	opts.classifiedStorage = etc.Get(defaultClassifiedStorageKey, defaultClassifiedStorage).Bool()

	return opts
}

// WithFile 设置输出的文件路径
func WithFile(file string) Option {
	return func(o *options) { o.file = file }
}

// WithLevel 设置输出的最低日志级别
func WithLevel(level Level) Option {
	return func(o *options) { o.level = level }
}

// WithFormat 设置输出的日志格式
func WithFormat(format Format) Option {
	return func(o *options) { o.format = format }
}

// WithStdout 设置是否输出到终端
func WithStdout(enable bool) Option {
	return func(o *options) { o.stdout = enable }
}

// WithTimeFormat 设置时间格式
func WithTimeFormat(format string) Option {
	return func(o *options) { o.timeFormat = format }
}

// WithStackLevel 设置堆栈的最小输出级别
func WithStackLevel(level Level) Option {
	return func(o *options) { o.stackLevel = level }
}

// WithFileMaxAge 设置文件最大留存时间
func WithFileMaxAge(maxAge time.Duration) Option {
	return func(o *options) { o.fileMaxAge = maxAge }
}

// WithFileMaxSize 设置输出的单个文件尺寸限制
func WithFileMaxSize(size int64) Option {
	return func(o *options) { o.fileMaxSize = size }
}

// WithFileCutRule 设置文件切割规则
func WithFileCutRule(cutRule CutRule) Option {
	return func(o *options) { o.fileCutRule = cutRule }
}

// WithCallerSkip 设置调用者跳过的层级深度
func WithCallerSkip(skip int) Option {
	return func(o *options) { o.callerSkip = skip }
}

// WithCallerFullPath 设置是否启用调用文件全路径
func WithCallerFullPath(enable bool) Option {
	return func(o *options) { o.callerFullPath = enable }
}

// WithClassifiedStorage 设置启用文件分级存储
// 启用后，日志将进行分级存储，大一级的日志将存储于小于等于自身的日志级别文件中
// 例如：InfoLevel级的日志将存储于due.debug.20220910.log、due.info.20220910.log两个日志文件中
func WithClassifiedStorage(enable bool) Option {
	return func(o *options) { o.classifiedStorage = enable }
}
