package flog

import (
	"fmt"
	"log"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

type fileBean struct {
	id           string
	dir          string
	filename     string
	_suffix      int
	_date        *time.Time
	mu           *sync.RWMutex
	logfile      *os.File
	lg           *log.Logger
	filesize     int64
	maxFileSize  int64
	maxFileCount int32
}

func newFileBean(fileDir, fileName string, _suffix int, maxSize int64, maxfileCount int32) (fb *fileBean) {
	t, _ := time.Parse(_DATEFORMAT, time.Now().Format(_DATEFORMAT))
	fb = &fileBean{dir: fileDir, filename: fileName, _date: &t, mu: new(sync.RWMutex)}
	fb.logfile, _ = os.OpenFile(fileDir+"/"+fileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	fb.lg = log.New(fb.logfile, "", log.Ldate|log.Ltime|log.Lshortfile)
	fb._suffix = _suffix
	fb.maxFileSize = maxSize
	fb.maxFileCount = maxfileCount
	fb.filesize = fileSize(fileDir + "/" + fileName)
	fb._date = &t
	return
}

func (p *fileBean) nextSuffix() int {
	return int(p._suffix%int(p.maxFileCount) + 1)
}

func (p *fileBean) rename(rolltype _ROLLTYPE) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.close()
	nextfilename := ""
	switch rolltype {
	case _DAILY:
		nextfilename = fmt.Sprint(p.dir, "/", p.filename, ".", p._date.Format(_DATEFORMAT))
	case _ROLLFILE:
		nextfilename = fmt.Sprint(p.dir, "/", p.filename, ".", p.nextSuffix())
		p._suffix = p.nextSuffix()
	}
	if isExist(nextfilename) {
		os.Remove(nextfilename)
	}
	os.Rename(p.dir+"/"+p.filename, nextfilename)
	t, _ := time.Parse(_DATEFORMAT, time.Now().Format(_DATEFORMAT))
	p._date = &t
	p.logfile, _ = os.OpenFile(p.dir+"/"+p.filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	p.lg = log.New(p.logfile, "", log.Ldate|log.Ltime|log.Lshortfile)
	p.filesize = fileSize(p.dir + "/" + p.filename)
}

func (p *fileBean) addsize(size int64) {
	atomic.AddInt64(&p.filesize, size)
}

func (p *fileBean) write(level string, v ...interface{}) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	s := fmt.Sprint(v...)
	p.lg.Output(skip+1, fmt.Sprintln(level, s))
}

func (p *fileBean) writef(format string, v ...interface{}) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	p.lg.Output(skip+1, fmt.Sprintf(format, v...))
}

func (p *fileBean) isOverSize() bool {
	return p.filesize >= p.maxFileSize
}

func (p *fileBean) close() {
	p.logfile.Close()
}
