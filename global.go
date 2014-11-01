package main

import (
	"crypto/md5"
	"os"
)

var global struct {
	hostname       string
	machineId      []byte
	pid            int
	requestCounter uint32
}

func init() {
	global.pid = os.Getpid()
	var err error
	global.hostname, err = os.Hostname()
	if err != nil {
		panic("failed to get hostname: " + err.Error())
	}

	hw := md5.New()
	if _, err := hw.Write([]byte(global.hostname)); err != nil {
		panic("unable to md5 hostname: " + err.Error())
	}
	global.machineId = make([]byte, 3)
	copy(global.machineId, hw.Sum(nil))
}
