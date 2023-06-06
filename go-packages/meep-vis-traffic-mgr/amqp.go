/*
 * Copyright (c) 2022  The AdvantEDGE Authors
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

package vistrafficmgr

import (
	//"encoding/hex"
	"net/url"
	//"errors"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	//amqp "github.com/hadihammurabi/go-rabbitmq"
)

type message_broker_amqp struct {
	running bool
}

func (amqp *message_broker_amqp) Init(tm *TrafficMgr) (err error) {
	log.Info("message_broker_amqp: Init")

	amqp.running = false

	u, err := url.ParseRequestURI(tm.broker)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	log.Info("url:%v\nscheme:%v host:%v Path:%v Port:%s", u, u.Scheme, u.Hostname(), u.Path, u.Port())

	// TODO

	return nil
}

func (amqp *message_broker_amqp) Run(tm *TrafficMgr) (err error) {
	log.Info("message_broker_amqp: Run")

	// TODO

	return nil
}

func (amqp *message_broker_amqp) Stop(tm *TrafficMgr) (err error) {
	log.Info("message_broker_amqp: Stop")

	// TODO

	return nil
}

func (amqp *message_broker_amqp) Send(tm *TrafficMgr, msgContent string, msgEncodeFormat string, stdOrganization string, msgType *int32) (err error) {
	log.Info("message_broker_amqp: Send")

	// TODO

	return nil
}
