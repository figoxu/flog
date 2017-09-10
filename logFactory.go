package flog

import "github.com/quexer/utee"

var (
	logMap = utee.SyncMap{}
)

func GetLog(key string)*logBean{
	m,ok:=logMap.Get(key)
	if !ok {
		lb:=&logBean{}
		logMap.Put(key,lb)
		m=lb
	}
	return m.(*logBean)
}