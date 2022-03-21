package lklogger

import (
	"testing"
)

var log *LkLogger

func TestA(t *testing.T) {
	log = NewLKLogger(false,false)
	log.Debug("hello debug")
	log.Info("hello info")
	log.Warn("hello warn")
	log.Error("hello error")

	log.DebugSf("hello debug %s","123")
	log.InfoSf("hello info %s","123")
	log.WarnSf("hello warn %s","123")
	log.ErrorSf("hello error %s","123")
}

func TestB(t *testing.T) {
	log = NewLKLogger(true,false)
	log.Debug("hello debug")
	log.Info("hello info")
	log.Warn("hello warn")
	log.Error("hello error")
}

func TestC(t *testing.T) {
	callErrorFunc := func(){
		log.Error("hello error")
	}
	log = NewLKLogger(true,true)
	log.Debug("hello debug")
	log.Info("hello info")
	log.Warn("hello warn")
	callErrorFunc()
}

func TestD(t *testing.T) {
	callErrorFunc := func(){
		log.Error("hello error")
	}
	log = NewLKLoggerAll(true,true)
	log.Debug("hello debug")
	log.Info("hello info")
	log.Warn("hello warn")
	callErrorFunc()
}
