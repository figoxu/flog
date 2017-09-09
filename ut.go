package flog

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"runtime/debug"
)

func mkdirlog(dir string) (e error) {
	_, er := os.Stat(dir)
	b := er == nil || os.IsExist(er)
	if !b {
		if err := os.MkdirAll(dir, 0666); err != nil {
			if os.IsPermission(err) {
				e = err
			}
		}
	}
	return
}

func fileSize(file string) int64 {
	f, e := os.Stat(file)
	if e != nil {
		fmt.Println(e.Error())
		return 0
	}
	return f.Size()
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func md5str(s string) string {
	m := md5.New()
	m.Write([]byte(s))
	return hex.EncodeToString(m.Sum(nil))
}

func catchError() {
	if err := recover(); err != nil {
		fmt.Println(string(debug.Stack()))
	}
}
