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

func convertJsonToMobilityProcedureSubscription(jsonInfo string) *MobilityProcedureSubscription {
	var obj MobilityProcedureSubscription
	err := json.Unmarshal([]byte(jsonInfo), &obj)
	if err != nil {
		log.Error(err.Error())
		return nil
	}
	return &obj
}

func convertMobilityProcedureSubscriptionToJson(obj *MobilityProcedureSubscription) string {
	jsonInfo, err := json.Marshal(*obj)
	if err != nil {
		log.Error(err.Error())
		return ""
	}
	return string(jsonInfo)
}

func convertMobilityProcedureNotificationToJson(obj *MobilityProcedureNotification) string {
	jsonInfo, err := json.Marshal(*obj)
	if err != nil {
		log.Error(err.Error())
		return ""
	}
	return string(jsonInfo)
}

func convertJsonToAdjacentAppInfoSubscription(jsonInfo string) *AdjacentAppInfoSubscription {
	var obj AdjacentAppInfoSubscription
	err := json.Unmarshal([]byte(jsonInfo), &obj)
	if err != nil {
		log.Error(err.Error())
		return nil
	}
	return &obj
}

func convertAdjacentAppInfoSubscriptionToJson(obj *AdjacentAppInfoSubscription) string {
	jsonInfo, err := json.Marshal(*obj)
	if err != nil {
		log.Error(err.Error())
		return ""
	}
	return string(jsonInfo)
}

func convertAdjacentAppInfoNotificationToJson(obj *AdjacentAppInfoNotification) string {
	jsonInfo, err := json.Marshal(*obj)
	if err != nil {
		log.Error(err.Error())
		return ""
	}
	return string(jsonInfo)
}

func convertRegistrationInfoToJson(obj *RegistrationInfo) string {
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

func convertSubscriptionLinkListToJson(obj *SubscriptionLinkList) string {
	jsonInfo, err := json.Marshal(*obj)
	if err != nil {
		log.Error(err.Error())
		return ""
	}
	return string(jsonInfo)
}

func convertDevInfoToJson(obj *DevInfo) string {
	jsonInfo, err := json.Marshal(*obj)
	if err != nil {
		log.Error(err.Error())
		return ""
	}
	return string(jsonInfo)
}

func convertProblemDetailstoJson(probdetails *ProblemDetails) string {
	jsonInfo, err := json.Marshal(*probdetails)
	if err != nil {
		log.Error(err.Error())
		return ""
	}
	return string(jsonInfo)
}

// func convertJsonToDevInfo(jsonInfo string) *DevInfo {
// 	var obj DevInfo
// 	err := json.Unmarshal([]byte(jsonInfo), &obj)
// 	if err != nil {
// 		log.Error(err.Error())
// 		return nil
// 	}
// 	return &obj
// }
