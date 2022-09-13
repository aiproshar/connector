package syncsys

import "sync"

var mtx = sync.Mutex{}

var concurrent = 0

func Add() {
	mtx.Lock()
	concurrent++
	mtx.Unlock()
}
func Delete() {
	mtx.Lock()
	concurrent--
	mtx.Unlock()
}

func ReadCount() int {
	return concurrent
}
