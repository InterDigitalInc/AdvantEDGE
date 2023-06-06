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

type message_broker_interface interface {
	Init(tm *TrafficMgr) (err error)

	Run(tm *TrafficMgr) (err error)

	Stop(tm *TrafficMgr) (err error)

	Send(tm *TrafficMgr, msgContent string, msgEncodeFormat string, stdOrganization string, msgType *int32) (err error)
}
