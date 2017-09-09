package flog

type logger struct {
	lb *logBean
}

func (p *logger) SetConsole(isConsole bool) {
	p.lb.setConsole(isConsole)
}

func (p *logger) SetLevel(_level LEVEL) {
	p.lb.setLevel(_level)
}

func (p *logger) SetFormat(logFormat string) {
	p.lb.setFormat(logFormat)
}

func (p *logger) SetRollingFile(fileDir, fileName string, maxNumber int32, maxSize int64, _unit UNIT) {
	p.lb.setRollingFile(fileDir, fileName, maxNumber, maxSize, _unit)
}

func (p *logger) SetRollingDaily(fileDir, fileName string) {
	p.lb.setRollingDaily(fileDir, fileName)
}

func (p *logger) Debug(v ...interface{}) {
	p.lb.debug(v...)
}

func (p *logger) Info(v ...interface{}) {
	p.lb.info(v...)
}
func (p *logger) Warn(v ...interface{}) {
	p.lb.warn(v...)
}
func (p *logger) Error(v ...interface{}) {
	p.lb.error(v...)
}
func (p *logger) Fatal(v ...interface{}) {
	p.lb.fatal(v...)
}

func (p *logger) SetLevelFile(level LEVEL, dir, fileName string) {
	p.lb.setLevelFile(level, dir, fileName)
}
