/*
 * Copyright (c) 2019  InterDigital Communications, Inc
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package logger

import (
	"fmt"
	"path"
	"runtime"

	logrus "github.com/sirupsen/logrus"
)

var componentName string

type Fields map[string]interface{}

func MeepTextLogInit(name string) {
	Formatter := new(logrus.TextFormatter)
	Formatter.TimestampFormat = "2006-01-02T15:04:05.999Z07:00"
	Formatter.FullTimestamp = true
	logrus.SetFormatter(Formatter)
	//logrus.SetLevel(logrus.TraceLevel)
	logrus.SetLevel(logrus.DebugLevel)
	componentName = name
}

func MeepJSONLogInit(name string) {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	//logrus.SetLevel(logrus.TraceLevel)
	logrus.SetLevel(logrus.DebugLevel)
	componentName = name
}

// getLogCaller
// setReportCaller cannot ignore caller levels, so cannot work in a wrapper. Fix not merged in github, so doing it manually
func getLogCaller() string {
	_, file, line, _ := runtime.Caller(2)

	location := fmt.Sprintf("%v:%v", path.Base(file), line)
	return location
}

func Info(args ...interface{}) {
	logrus.WithFields(logrus.Fields{
		"meep.component": componentName,
	}).Info(args...)
}

func Debug(args ...interface{}) {
	logrus.WithFields(logrus.Fields{
		"meep.component": componentName,
		"meep.from":      getLogCaller(),
	}).Debug(args...)
}

func Trace(args ...interface{}) {
	logrus.WithFields(logrus.Fields{
		"meep.component": componentName,
		"meep.from":      getLogCaller(),
	}).Trace(args...)
}

func Warn(args ...interface{}) {
	logrus.WithFields(logrus.Fields{
		"meep.component": componentName,
		"meep.from":      getLogCaller(),
	}).Warn(args...)
}

func Error(args ...interface{}) {
	logrus.WithFields(logrus.Fields{
		"meep.component": componentName,
		"meep.from":      getLogCaller(),
	}).Error(args...)
}

func Panic(args ...interface{}) {
	logrus.WithFields(logrus.Fields{
		"meep.component": componentName,
		"meep.from":      getLogCaller(),
	}).Panic(args...)
}

func Fatal(args ...interface{}) {
	logrus.WithFields(logrus.Fields{
		"meep.component": componentName,
		"meep.from":      getLogCaller(),
	}).Fatal(args...)
}

func WithFields(fields Fields) *logrus.Entry {
	return logrus.WithFields(logrus.Fields(fields))
}
