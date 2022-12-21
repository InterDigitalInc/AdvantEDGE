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

package server

import (
	"encoding/json"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
)

func convertJsonToApInfoComplete(jsonInfo string) *ApInfoComplete {
	var obj ApInfoComplete
	err := json.Unmarshal([]byte(jsonInfo), &obj)
	if err != nil {
		log.Error(err.Error())
		return nil
	}
	return &obj
}

func convertApInfoCompleteToJson(obj *ApInfoComplete) string {
	jsonInfo, err := json.Marshal(*obj)
	if err != nil {
		log.Error(err.Error())
		return ""
	}
	return string(jsonInfo)
}

func convertJsonToStaData(jsonData string) *StaData {
	var obj StaData
	err := json.Unmarshal([]byte(jsonData), &obj)
	if err != nil {
		log.Error(err.Error())
		return nil
	}
	return &obj
}

func convertStaDataToJson(obj *StaData) string {
	jsonData, err := json.Marshal(*obj)
	if err != nil {
		log.Error(err.Error())
		return ""
	}
	return string(jsonData)
}

func convertAssocStaSubscriptionToJson(obj *AssocStaSubscription) string {
	jsonInfo, err := json.Marshal(*obj)
	if err != nil {
		log.Error(err.Error())
		return ""
	}
	return string(jsonInfo)
}

func convertJsonToAssocStaSubscription(jsonData string) *AssocStaSubscription {
	var obj AssocStaSubscription
	err := json.Unmarshal([]byte(jsonData), &obj)
	if err != nil {
		log.Error(err.Error())
		return nil
	}
	return &obj
}

func convertJsonToStaDataRateSubscription(jsonData string) *StaDataRateSubscription {
	var obj StaDataRateSubscription
	err := json.Unmarshal([]byte(jsonData), &obj)
	if err != nil {
		log.Error(err.Error())
		return nil
	}
	return &obj
}

func convertJsonToMeasurementReportSubscription(jsonData string) *MeasurementReportSubscription {
	var obj MeasurementReportSubscription
	err := json.Unmarshal([]byte(jsonData), &obj)
	if err != nil {
		log.Error(err.Error())
		return nil
	}
	return &obj
}

func convertSubscriptionLinkListToJson(obj *SubscriptionLinkList) string {
	jsonInfo, err := json.Marshal(*obj)
	if err != nil {
		log.Error(err.Error())
		return ""
	}
	return string(jsonInfo)
}

func convertAssocStaNotificationToJson(obj *AssocStaNotification) string {
	jsonInfo, err := json.Marshal(*obj)
	if err != nil {
		log.Error(err.Error())
		return ""
	}
	return string(jsonInfo)
}

func convertStaDataRateNotificationToJson(obj *StaDataRateNotification) string {
	jsonInfo, err := json.Marshal(*obj)
	if err != nil {
		log.Error(err.Error())
		return ""
	}
	return string(jsonInfo)
}

func convertExpiryNotificationToJson(obj *ExpiryNotification) string {
	jsonInfo, err := json.Marshal(*obj)
	if err != nil {
		log.Error(err.Error())
		return ""
	}
	return string(jsonInfo)
}

func convertTestNotificationToJson(obj *TestNotification) string {
	jsonInfo, err := json.Marshal(*obj)
	if err != nil {
		log.Error(err.Error())
		return ""
	}
	return string(jsonInfo)
}

func convertApInfoListToJson(obj *[]ApInfo) string {
	jsonInfo, err := json.Marshal(*obj)
	if err != nil {
		log.Error(err.Error())
		return ""
	}
	return string(jsonInfo)
}

func convertStaInfoListToJson(obj *[]StaInfo) string {
	jsonInfo, err := json.Marshal(*obj)
	if err != nil {
		log.Error(err.Error())
		return ""
	}
	return string(jsonInfo)
}
