package logmain

import (
	log "github.com/sirupsen/logrus"
)

const meepComponentField = "meep.component"
const meepComponent = "mg-manager"

// MeepJSONLogInit -- Log init function
func MeepJSONLogInit() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)
}

// Info -- Info logs
func Info(args ...interface{}) {
	log.WithFields(log.Fields{meepComponentField: meepComponent}).Info(args)
}

// Debug -- Debug logs
func Debug(args ...interface{}) {
	log.WithFields(log.Fields{meepComponentField: meepComponent}).Debug(args)
}

// Warn -- Warn logs
func Warn(args ...interface{}) {
	log.WithFields(log.Fields{meepComponentField: meepComponent}).Warn(args)
}

// Error -- Error logs
func Error(args ...interface{}) {
	log.WithFields(log.Fields{meepComponentField: meepComponent}).Error(args)
}

// Panic -- Panic logs
func Panic(args ...interface{}) {
	log.WithFields(log.Fields{meepComponentField: meepComponent}).Panic(args)
}

// Fatal -- Fatal logs
func Fatal(args ...interface{}) {
	log.WithFields(log.Fields{meepComponentField: meepComponent}).Fatal(args)
}

// WithFields -- Log with fields
func WithFields(fields log.Fields) *log.Entry {
	return log.WithFields(fields)
}
