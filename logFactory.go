package flog

import (
	"github.com/quexer/utee"
	"sync"
)

var (
	logMap = utee.SyncMap{}
)

func GetLog(key string)*logBean{
	m,ok:=logMap.Get(key)
	if !ok {
		lb:=&logBean{}
		lb.mu = new(sync.Mutex)
		lb.SetConsole(true)
		logMap.Put(key,lb)
		m=lb
	}
	return m.(*logBean)
}