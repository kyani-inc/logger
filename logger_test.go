package logger

import (
	"testing"
)

// TODO: Hook up a unit test for sumologic
var log = New(DefaultConfig())

func TestNewLogger(t *testing.T) {
	method := "this is the method"
	url := "this is the url"
	log.Infof("[BDT] Calling:%s %s", method, url)
}

func TestNewLoggerSumo(t *testing.T) {

	a := New(Config{Appname: "logger"})

	method := "this is the method"
	url := "this is the url"
	log.Infof("[BDT] Calling:%s %s", method, url)
}
