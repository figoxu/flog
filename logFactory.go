package flog

import (
	"github.com/quexer/utee"
	"sync"
)

var (
	logMap = utee.SyncMap{}
)

func GetLog(key string)*LogBean {
	m,ok:=logMap.Get(key)
	if !ok {
		lb:=&LogBean{}
		lb.mu = new(sync.Mutex)
		lb.SetConsole(true)
		logMap.Put(key,lb)
		m=lb
	}
	return m.(*LogBean)
}