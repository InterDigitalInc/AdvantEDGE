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
	"os"

	log "github.com/sirupsen/logrus"
)

func MeepJSONLogInit() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)
	//force output to stdout rather than default stderr
	log.SetOutput(os.Stdout)
}

func Info(args ...interface{}) {
	log.WithFields(log.Fields{
		"meep.component": "virt-engine",
	}).Info(args)
}

func Debug(args ...interface{}) {
	log.WithFields(log.Fields{
		"meep.component": "virt-engine",
	}).Debug(args)
}

func Warn(args ...interface{}) {
	log.WithFields(log.Fields{
		"meep.component": "virt-engine",
	}).Warn(args)
}
func Error(args ...interface{}) {
	log.WithFields(log.Fields{
		"meep.component": "virt-engine",
	}).Error(args)
}
func Panic(args ...interface{}) {
	log.WithFields(log.Fields{
		"meep.component": "virt-engine",
	}).Panic(args)
}

func Fatal(args ...interface{}) {
	log.WithFields(log.Fields{
		"meep.component": "virt-engine",
	}).Fatal(args)
}

func WithFields(fields log.Fields) *log.Entry {
	return log.WithFields(fields)
}
