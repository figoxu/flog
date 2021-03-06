package flog

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"
	"strconv"
)

type LEVEL int32
type UNIT int64
type _ROLLTYPE int

const _DATEFORMAT = "2006-01-02"

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

var skip int = 4

type LogBean struct {
	mu              *sync.Mutex
	logLevel        LEVEL
	maxFileSize     int64
	maxFileCount    int32
	consoleAppender bool
	rolltype        _ROLLTYPE
	format          string
	id              string
	d, i, w, e, f   string //id
}

func (p *LogBean) SetConsole(isConsole bool) {
	p.consoleAppender = isConsole
}

func (p *LogBean) SetLevelFile(level LEVEL, dir, fileName string) {
	key := md5str(fmt.Sprint(dir, fileName))
	switch level {
	case DEBUG:
		p.d = key
	case INFO:
		p.i = key
	case WARN:
		p.w = key
	case ERROR:
		p.e = key
	case FATAL:
		p.f = key
	default:
		return
	}
	var _suffix = 0
	if p.maxFileCount < 1<<31-1 {
		for i := 1; i < int(p.maxFileCount); i++ {
			if isExist(dir + "/" + fileName + "." + strconv.Itoa(i)) {
				_suffix = i
			} else {
				break
			}
		}
	}
	fbf.add(dir, fileName, _suffix, p.maxFileSize, p.maxFileCount)
}

func (p *LogBean) SetLevel(_level LEVEL) {
	p.logLevel = _level
}

func (p *LogBean) SetFormat(logFormat string) {
	p.format = logFormat
}

func (p *LogBean) SetRollingFile(fileDir, fileName string, maxNumber int32, maxSize int64, _unit UNIT) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if maxNumber > 0 {
		p.maxFileCount = maxNumber
	} else {
		p.maxFileCount = 1<<31 - 1
	}
	p.maxFileSize = maxSize * int64(_unit)
	p.rolltype = _ROLLFILE
	mkdirlog(fileDir)
	var _suffix = 0
	for i := 1; i < int(maxNumber); i++ {
		if isExist(fileDir + "/" + fileName + "." + strconv.Itoa(i)) {
			_suffix = i
		} else {
			break
		}
	}
	p.id = md5str(fmt.Sprint(fileDir, fileName))
	fbf.add(fileDir, fileName, _suffix, p.maxFileSize, p.maxFileCount)
}

func (p *LogBean) SetRollingDaily(fileDir, fileName string) {
	p.rolltype = _DAILY
	mkdirlog(fileDir)
	p.id = md5str(fmt.Sprint(fileDir, fileName))
	fbf.add(fileDir, fileName, 0, 0, 0)
}

func (p *LogBean) console(v ...interface{}) {
	s := fmt.Sprint(v...)
	if p.consoleAppender {
		_, file, line, _ := runtime.Caller(skip)
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		if p.format == "" {
			log.Println(file, strconv.Itoa(line), s)
		} else {
			vs := make([]interface{}, 0)
			vs = append(vs, file)
			vs = append(vs, strconv.Itoa(line))
			for _, vv := range v {
				vs = append(vs, vv)
			}
			format := fmt.Sprint(p.format," \n ", )
			log.Printf(format, vs...)
		}
	}
}

func (p *LogBean) log(level string, v ...interface{}) {
	defer catchError()
	s := fmt.Sprint(v...)
	length := len([]byte(s))
	var lg *fileBean = fbf.get(p.id)
	var _level = ALL
	switch level {
	case "debug":
		if p.d != "" {
			lg = fbf.get(p.d)
		}
		_level = DEBUG
	case "info":
		if p.i != "" {
			lg = fbf.get(p.i)
		}
		_level = INFO
	case "warn":
		if p.w != "" {
			lg = fbf.get(p.w)
		}
		_level = WARN
	case "error":
		if p.e != "" {
			lg = fbf.get(p.e)
		}
		_level = ERROR
	case "fatal":
		if p.f != "" {
			lg = fbf.get(p.f)
		}
		_level = FATAL
	}
	if lg != nil {
		p.fileCheck(lg)
		lg.addsize(int64(length))
		if p.logLevel <= _level {
			if lg != nil {
				if p.format == "" {
					lg.write(level, s)
				} else {
					lg.writef(p.format, v...)
				}
			}
			p.console(v...)
		}
	} else {
		p.console(v...)
	}
}

func (p *LogBean) Debug(v ...interface{}) {
	p.log("debug", v...)
}
func (p *LogBean) Info(v ...interface{}) {
	p.log("info", v...)
}
func (p *LogBean) Warn(v ...interface{}) {
	p.log("warn", v...)
}
func (p *LogBean) Error(v ...interface{}) {
	p.log("error", v...)
}
func (p *LogBean) Fatal(v ...interface{}) {
	p.log("fatal", v...)
}

func (p *LogBean) fileCheck(fb *fileBean) {
	defer catchError()
	if p.isMustRename(fb) {
		p.mu.Lock()
		defer p.mu.Unlock()
		if p.isMustRename(fb) {
			fb.rename(p.rolltype)
		}
	}
}

//--------------------------------------------------------------------------------

func (p *LogBean) isMustRename(fb *fileBean) bool {
	switch p.rolltype {
	case _DAILY:
		t, _ := time.Parse(_DATEFORMAT, time.Now().Format(_DATEFORMAT))
		if t.After(*fb._date) {
			return true
		}
	case _ROLLFILE:
		return fb.isOverSize()
	}
	return false
}
