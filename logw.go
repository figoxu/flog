package flog

import (
	"sync"
)

var defaultlog *logBean = getdefaultLogger()
var skip int = 4

func GetLogger() (l *logger) {
	l = new(logger)
	l.lb = getdefaultLogger()
	return
}

func getdefaultLogger() (lb *logBean) {
	lb = &logBean{}
	lb.mu = new(sync.Mutex)
	lb.setConsole(true)
	return
}
