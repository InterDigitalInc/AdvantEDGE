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
		"meep.component": "sidecar",
	}).Info(args)
}

func Debug(args ...interface{}) {
	log.WithFields(log.Fields{
		"meep.component": "sidecar",
	}).Debug(args)
}

func Warn(args ...interface{}) {
	log.WithFields(log.Fields{
		"meep.component": "sidecar",
	}).Warn(args)
}

func Error(args ...interface{}) {
	log.WithFields(log.Fields{
		"meep.component": "sidecar",
	}).Error(args)
}

func Panic(args ...interface{}) {
	log.WithFields(log.Fields{
		"meep.component": "sidecar",
	}).Panic(args)
}

func Fatal(args ...interface{}) {
	log.WithFields(log.Fields{
		"meep.component": "sidecar",
	}).Fatal(args)
}

func WithFields(fields log.Fields) *log.Entry {
	return log.WithFields(fields)
}
