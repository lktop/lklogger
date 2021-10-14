package lklogger

import (
	"go.uber.org/zap"
	"testing"
)

var log *zap.Logger


func TestA(t *testing.T) {
	log = NewLKLogger(false,false)
	log.Debug("hello debug")
	log.Info("hello info")
	log.Warn("hello warn")
	log.Error("hello error")
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

