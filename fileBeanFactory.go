package flog

import (
	"fmt"
	"sync"
)

type fileBeanFactory struct {
	fbs map[string]*fileBean
	mu  *sync.RWMutex
}

var fbf = &fileBeanFactory{fbs: make(map[string]*fileBean, 0), mu: new(sync.RWMutex)}

func (p *fileBeanFactory) add(dir, filename string, _suffix int, maxsize int64, maxfileCount int32) {
	p.mu.Lock()
	defer p.mu.Unlock()
	id := md5str(fmt.Sprint(dir, filename))
	if _, ok := p.fbs[id]; !ok {
		p.fbs[id] = newFileBean(dir, filename, _suffix, maxsize, maxfileCount)
	}
}

func (p *fileBeanFactory) get(id string) *fileBean {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.fbs[id]
}
