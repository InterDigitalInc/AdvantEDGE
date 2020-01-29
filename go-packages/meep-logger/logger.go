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
	"net/http"
	"time"
	"runtime"
	"path"
	"fmt"
	logrus "github.com/sirupsen/logrus"
)

var componentName string

type Fields map[string]interface{}

func MeepTextLogInit(name string) {
	logrus.SetLevel(logrus.DebugLevel)
	//SetSetReportCaller(true)
	componentName = name
}

func MeepJSONLogInit(name string) {
	logrus.SetFormatter(&logrus.JSONFormatter{})
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
		"meep.time": time.Now().String(),
	}).Info(args...)
}

func Debug(args ...interface{}) {
	logrus.WithFields(logrus.Fields{
		"meep.component": componentName,
		"meep.from": getLogCaller(),
		"meep.time": time.Now().String(),
	}).Debug(args...)
}

func Trace(args ...interface{}) {
        logrus.WithFields(logrus.Fields{
                "meep.component": componentName,
                "meep.from": getLogCaller(),
                "meep.time": time.Now().String(),
        }).Trace(args...)
}

func Warn(args ...interface{}) {
	logrus.WithFields(logrus.Fields{
		"meep.component": componentName,
		"meep.from": getLogCaller(),
		"meep.time": time.Now().String(),
	}).Warn(args...)
}
func Error(args ...interface{}) {
	logrus.WithFields(logrus.Fields{
		"meep.component": componentName,
		"meep.from": getLogCaller(),
		"meep.time": time.Now().String(),
	}).Error(args...)
}
func Panic(args ...interface{}) {
	logrus.WithFields(logrus.Fields{
		"meep.component": componentName,
		"meep.from": getLogCaller(),
		"meep.time": time.Now().String(),
	}).Panic(args...)
}

func Fatal(args ...interface{}) {
	logrus.WithFields(logrus.Fields{
		"meep.component": componentName,
		"meep.from": getLogCaller(),
		"meep.time": time.Now().String(),
	}).Fatal(args...)
}

func WithFields(fields Fields) *logrus.Entry {
	return logrus.WithFields(logrus.Fields(fields))
}

func httpLog(args ...interface{}) {
        logrus.WithFields(logrus.Fields{
                "meep.component": componentName,
                "meep.from": getLogCaller(),
                "meep.time": time.Now().String(),
        }).Debug(args...)
}

func HTTP(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		httpLog(
			r.Method,
			" ", r.RequestURI,
			" apiFct: ", name,
			" execTime: ", time.Since(start))
	})
}
