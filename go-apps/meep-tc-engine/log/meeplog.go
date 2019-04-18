package logmain

import (
	log "github.com/sirupsen/logrus"
)

func MeepJSONLogInit() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)
}

func Info(args ...interface{}) {
	log.WithFields(log.Fields{
		"meep.component": "tc-engine",
	}).Info(args)
}

func Debug(args ...interface{}) {
	log.WithFields(log.Fields{
		"meep.component": "tc-engine",
	}).Debug(args)
}

func Warn(args ...interface{}) {
	log.WithFields(log.Fields{
		"meep.component": "tc-engine",
	}).Warn(args)
}
func Error(args ...interface{}) {
	log.WithFields(log.Fields{
		"meep.component": "tc-engine",
	}).Error(args)
}
func Panic(args ...interface{}) {
	log.WithFields(log.Fields{
		"meep.component": "tc-engine",
	}).Panic(args)
}

func Fatal(args ...interface{}) {
	log.WithFields(log.Fields{
		"meep.component": "tc-engine",
	}).Fatal(args)
}

func WithFields(fields log.Fields) *log.Entry {
	return log.WithFields(fields)
}
