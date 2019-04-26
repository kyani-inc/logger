package logger

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime"
	"strconv"
	"time"
)

func NewLogger(config Config) *Logger {
	return &Logger{
		Out:     os.Stderr,
		AppName: config.Appname,
		Host:    config.Host,
		Port:    config.Port,
	}
}

type Logger struct {
	Out     io.Writer
	AppName string
	Host    string
	Port    int
	Hook    *SumoLogicHook
}

type Log struct {
	AppName      string    `json:"app"`
	Date         time.Time `json:"date"`
	LogLevel     LogLevel  `json:"level"`
	Message      string    `json:"message"`
	HttpBodyDump string    `json:"http_body"`
	FuncCall     string    `json:"func_call"`
}

type LogLevel string

const (
	InfoLevel  LogLevel = "INFO"
	ErrorLevel LogLevel = "ERROR"
	PanicLevel LogLevel = "PANIC"
	FatalLevel LogLevel = "FATAL"
	WarnLevel  LogLevel = "WARN"
	DebugLevel LogLevel = "DEBUG"
)

func (logger *Logger) Info(message ...interface{}) {
	msg := fmt.Sprint(message)
	logger.logit(logger.Out, InfoLevel, msg, nil, nil)
}

func (logger *Logger) Infof(message string, args ...interface{}) {
	logger.logit(logger.Out, InfoLevel, fmt.Sprintf(message, args...), nil, nil)
}

func (logger *Logger) InfofWithRequest(req *http.Request, message string, args ...interface{}) {
	logger.logit(logger.Out, InfoLevel, fmt.Sprintf(message, args...), nil, req)
}

func (logger *Logger) InfofWithResponse(res *http.Response, message string, args ...interface{}) {
	logger.logit(logger.Out, InfoLevel, fmt.Sprintf(message, args...), res, nil)
}

func (logger *Logger) Error(message ...interface{}) {
	msg := fmt.Sprint(message)
	logger.logit(logger.Out, ErrorLevel, msg, nil, nil)
}

func (logger *Logger) Errorf(message string, args ...interface{}) {
	logger.logit(logger.Out, ErrorLevel, fmt.Sprintf(message, args...), nil, nil)
}

func (logger *Logger) ErrorfWithRequest(req *http.Request, message string, args ...interface{}) {
	logger.logit(logger.Out, ErrorLevel, fmt.Sprintf(message, args...), nil, req)
}

func (logger *Logger) ErrorfWithResponse(res *http.Response, message string, args ...interface{}) {
	logger.logit(logger.Out, ErrorLevel, fmt.Sprintf(message, args...), res, nil)
}

func (logger *Logger) logit(Out io.Writer, level LogLevel, message string, res *http.Response, req *http.Request) {
	pc, file, line, _ := runtime.Caller(2)
	funcCall := runtime.FuncForPC(pc).Name()

	var body []byte

	if res != nil {
		body, _ = httputil.DumpResponse(res, true)
	}
	if req != nil {
		body, _ = httputil.DumpRequest(req, true)
	}

	log := Log{HttpBodyDump: string(body), AppName: logger.AppName,
		Date:     time.Now().UTC(),
		LogLevel: level,
		Message:  message, FuncCall: fmt.Sprintf("FILE:%s, FUNC:%s, LINE:%s", file, funcCall, strconv.Itoa(line))}

	if logger.Hook != nil {
		logger.Hook.Fire(log)
	} else {
		bodyToWrite := ""
		if log.HttpBodyDump != "" {
			bodyToWrite = fmt.Sprintf("[%s] [%s] (%s) ***%s*** {%s} : %s", log.LogLevel, log.AppName, log.Date.String(), log.Message, log.HttpBodyDump, log.FuncCall)
		} else {
			bodyToWrite = fmt.Sprintf("[%s] [%s] (%s) ***%s*** %s", log.LogLevel, log.AppName, log.Date.String(), log.Message, log.FuncCall)
		}
		logger.Out.Write([]byte(bodyToWrite))
	}
}
