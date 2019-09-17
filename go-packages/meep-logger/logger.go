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
	log "github.com/sirupsen/logrus"
)

var componentName string

type Fields map[string]interface{}

func MeepTextLogInit(name string) {
	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.DebugLevel)
	componentName = name
}

func MeepJSONLogInit(name string) {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)
	componentName = name
}

func Info(args ...interface{}) {
	log.WithFields(log.Fields{
		"meep.component": componentName,
	}).Info(args...)
}

func Debug(args ...interface{}) {
	log.WithFields(log.Fields{
		"meep.component": componentName,
	}).Debug(args...)
}

func Warn(args ...interface{}) {
	log.WithFields(log.Fields{
		"meep.component": componentName,
	}).Warn(args...)
}
func Error(args ...interface{}) {
	log.WithFields(log.Fields{
		"meep.component": componentName,
	}).Error(args...)
}
func Panic(args ...interface{}) {
	log.WithFields(log.Fields{
		"meep.component": componentName,
	}).Panic(args...)
}

func Fatal(args ...interface{}) {
	log.WithFields(log.Fields{
		"meep.component": componentName,
	}).Fatal(args...)
}

func WithFields(fields Fields) *log.Entry {
	return log.WithFields(log.Fields(fields))
}
