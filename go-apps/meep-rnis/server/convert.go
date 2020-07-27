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

package server

import (
	"encoding/json"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
)

func convertJsonToEcgi(jsonInfo string) *Ecgi {

	var obj Ecgi
	err := json.Unmarshal([]byte(jsonInfo), &obj)
	if err != nil {
		log.Error(err.Error())
		return nil
	}
	return &obj
}

func convertEcgiToJson(obj *Ecgi) string {

	jsonInfo, err := json.Marshal(*obj)
	if err != nil {
		log.Error(err.Error())
		return ""
	}

	return string(jsonInfo)
}

func convertJsonToUeData(jsonData string) *UeData {

	var obj UeData
	err := json.Unmarshal([]byte(jsonData), &obj)
	if err != nil {
		log.Error(err.Error())
		return nil
	}
	return &obj
}

func convertUeDataToJson(obj *UeData) string {

	jsonData, err := json.Marshal(*obj)
	if err != nil {
		log.Error(err.Error())
		return ""
	}

	return string(jsonData)
}

func convertJsonToCellChangeSubscription(jsonInfo string) *CellChangeSubscription {

	var obj CellChangeSubscription
	err := json.Unmarshal([]byte(jsonInfo), &obj)
	if err != nil {
		log.Error(err.Error())
		return nil
	}
	return &obj
}

func convertCellChangeSubscriptionToJson(obj *CellChangeSubscription) string {

	jsonInfo, err := json.Marshal(*obj)
	if err != nil {
		log.Error(err.Error())
		return ""
	}

	return string(jsonInfo)
}

func convertJsonToRabEstSubscription(jsonInfo string) *RabEstSubscription {

	var obj RabEstSubscription
	err := json.Unmarshal([]byte(jsonInfo), &obj)
	if err != nil {
		log.Error(err.Error())
		return nil
	}
	return &obj
}

func convertRabEstSubscriptionToJson(obj *RabEstSubscription) string {

	jsonInfo, err := json.Marshal(*obj)
	if err != nil {
		log.Error(err.Error())
		return ""
	}

	return string(jsonInfo)
}

func convertJsonToRabRelSubscription(jsonInfo string) *RabRelSubscription {

	var obj RabRelSubscription
	err := json.Unmarshal([]byte(jsonInfo), &obj)
	if err != nil {
		log.Error(err.Error())
		return nil
	}
	return &obj
}

func convertRabRelSubscriptionToJson(obj *RabRelSubscription) string {

	jsonInfo, err := json.Marshal(*obj)
	if err != nil {
		log.Error(err.Error())
		return ""
	}

	return string(jsonInfo)
}
