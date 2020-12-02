/*
 * Copyright (c) 2020  InterDigital Communications, Inc
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

package helm

import (
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
)

type Task string

const (
	Install Task = "INSTALL"
	Delete  Task = "DELETE"
)

type Job struct {
	task        Task
	charts      []Chart
	sandboxName string
}

var queue *chan Job = nil

func startWorker() {
	if queue != nil {
		return
	}
	queueChan := make(chan Job, 5)
	queue = &queueChan

	go func() {
		for job := range queueChan {
			switch job.task {
			case Install:
				log.Debug("Installing ", len(job.charts), " Charts...")
				_ = installCharts(job.charts, job.sandboxName)
				log.Debug("Charts installed (", len(job.charts), ")")

			case Delete:
				log.Debug("Deleting ", len(job.charts), " Releases...")
				_ = deleteReleases(job.charts)
				log.Debug("Releases deleted (", len(job.charts), ")")
			}
		}
		queue = nil
	}()
}

func runTask(task Task, charts []Chart, sandboxName string) error {
	startWorker()
	var job Job = Job{task: task, charts: charts, sandboxName: sandboxName}
	*queue <- job
	return nil
}
