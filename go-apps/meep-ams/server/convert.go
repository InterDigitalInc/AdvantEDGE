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

/*
func convertJsonToAppInfo(jsonInfo string) *AppInfo {

	var obj AppInfo
	err := json.Unmarshal([]byte(jsonInfo), &obj)
	if err != nil {
		log.Error(err.Error())
		return nil
	}
	return &obj
}

func convertAppInfoToJson(obj *AppInfo) string {

	jsonInfo, err := json.Marshal(*obj)
	if err != nil {
		log.Error(err.Error())
		return ""
	}

	return string(jsonInfo)
}

func convertJsonToPoaInfo(jsonInfo string) *PoaInfo {

	var obj PoaInfo
	err := json.Unmarshal([]byte(jsonInfo), &obj)
	if err != nil {
		log.Error(err.Error())
		return nil
	}
	return &obj
}

func convertPoaInfoToJson(obj *PoaInfo) string {

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

func convertJsonToDomainData(jsonData string) *DomainData {

	var obj DomainData
	err := json.Unmarshal([]byte(jsonData), &obj)
	if err != nil {
		log.Error(err.Error())
		return nil
	}
	return &obj
}

func convertDomainDataToJson(obj *DomainData) string {

	jsonData, err := json.Marshal(*obj)
	if err != nil {
		log.Error(err.Error())
		return ""
	}

	return string(jsonData)
}
*/
/*
func convertJsonToOneOfNotificationSubscription(jsonInfo string) *OneOfNotificationSubscription {

        var obj OneOfNotificationSubscription
        err := json.Unmarshal([]byte(jsonInfo), &obj)
        if err != nil {
                log.Error(err.Error())
                return nil
        }
        return &obj
}

func convertOneOfNotificationSubscriptionToJson(obj *OneOfNotificationSubscription) string {

        jsonInfo, err := json.Marshal(*obj)
        if err != nil {
                log.Error(err.Error())
                return ""
        }

        return string(jsonInfo)
}
*/
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

/*
func convertJsonToAdjacentAppInfoSubscription(jsonInfo string) *AdjacentAppInfoSubscription {

        var obj AdjacentAppInfoSubscription
        err := json.Unmarshal([]byte(jsonInfo), &obj)
        if err != nil {
                log.Error(err.Error())
                return nil
        }
        return &obj
}
*/
func convertAdjacentAppInfoSubscriptionToJson(obj *AdjacentAppInfoSubscription) string {

	jsonInfo, err := json.Marshal(*obj)
	if err != nil {
		log.Error(err.Error())
		return ""
	}

	return string(jsonInfo)
}

func convertJsonToRegistrationInfo(jsonInfo string) *RegistrationInfo {

	var obj RegistrationInfo
	err := json.Unmarshal([]byte(jsonInfo), &obj)
	if err != nil {
		log.Error(err.Error())
		return nil
	}
	return &obj
}

func convertRegistrationInfoToJson(obj *RegistrationInfo) string {

	jsonInfo, err := json.Marshal(*obj)
	if err != nil {
		log.Error(err.Error())
		return ""
	}

	return string(jsonInfo)
}
