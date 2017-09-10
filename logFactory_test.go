package flog

import "testing"

func TestGetLog(t *testing.T) {
	log_figo := GetLog("figo")
	log_figo.SetRollingDaily("./","figo.log")
	log_figo.Error("Hello")
	log_figo.Error("world")
	log_andy := GetLog("andy")
	log_andy.SetRollingDaily("./","andy.log")
	log_andy.Error("foo")
	log_andy.Error("bar")
}