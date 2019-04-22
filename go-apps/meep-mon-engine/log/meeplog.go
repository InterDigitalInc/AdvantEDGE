/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
*/
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
		"meep.component": "mon-engine",
	}).Info(args)
}

func Debug(args ...interface{}) {
	log.WithFields(log.Fields{
		"meep.component": "mon-engine",
	}).Debug(args)
}

func Warn(args ...interface{}) {
	log.WithFields(log.Fields{
		"meep.component": "mon-engine",
	}).Warn(args)
}
func Error(args ...interface{}) {
	log.WithFields(log.Fields{
		"meep.component": "mon-engine",
	}).Error(args)
}
func Panic(args ...interface{}) {
	log.WithFields(log.Fields{
		"meep.component": "mon-engine",
	}).Panic(args)
}

func Fatal(args ...interface{}) {
	log.WithFields(log.Fields{
		"meep.component": "mon-engine",
	}).Fatal(args)
}

func WithFields(fields log.Fields) *log.Entry {
	return log.WithFields(fields)
}
