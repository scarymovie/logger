package logger

import "sync"

func ResetGlobal() {
	once = sync.Once{}
	global = nil
}
