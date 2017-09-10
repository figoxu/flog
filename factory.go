package flog

type LEVEL int32
type UNIT int64
type _ROLLTYPE int

const _DATEFORMAT = "2006-01-02"

var logLevel LEVEL = 1
var skip int = 4

const (
	_       = iota
	KB UNIT = 1 << (iota * 10)
	MB
	GB
	TB
)

const (
	ALL LEVEL = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
	OFF
)

const (
	_DAILY _ROLLTYPE = iota
	_ROLLFILE
)