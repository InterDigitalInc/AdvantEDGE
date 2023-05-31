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

package meepdaimgr

import (
	"encoding/json"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
)

func convertPolygonToJson(area *Polygon) string {
	if area == nil { // This is not an error
		return ""
	}
	jsonInfo, err := json.Marshal(*area)
	if err != nil {
		log.Error(err.Error())
		return ""
	}
	return string(jsonInfo)
}

func convertJsonToPolygon(jsonInfo string) *Polygon {
	var obj Polygon
	if jsonInfo != "" { // Empty string is processed as an empty array
		err := json.Unmarshal([]byte(jsonInfo), &obj)
		if err != nil {
			log.Error(err.Error())
			return nil
		}
	}
	return &obj
}

func convertCivicAddressElementToJson(civicAddressElement *CivicAddressElement) string {
	if civicAddressElement == nil { // This is not an error
		return ""
	}
	jsonInfo, err := json.Marshal(*civicAddressElement)
	if err != nil {
		log.Error(err.Error())
		return ""
	}
	return string(jsonInfo)
}

func convertJsonToCivicAddressElement(jsonInfo string) *CivicAddressElement {
	var obj CivicAddressElement
	if jsonInfo != "" { // Empty string is processed as an empty array
		err := json.Unmarshal([]byte(jsonInfo), &obj)
		if err != nil {
			log.Error(err.Error())
			return nil
		}
	}
	return &obj
}

func convertApplicationListToJson(applicationList *ApplicationList) string {
	if applicationList == nil { // This is not an error
		return ""
	}
	jsonInfo, err := json.Marshal(*applicationList)
	if err != nil {
		log.Error(err.Error())
		return ""
	}
	return string(jsonInfo)
}

func convertJsonToApplicationList(jsonInfo string) *ApplicationList {
	var obj ApplicationList
	if jsonInfo != "" { // Empty string is processed as an empty array
		err := json.Unmarshal([]byte(jsonInfo), &obj)
		if err != nil {
			log.Error(err.Error())
			return nil
		}
	}
	return &obj
}

func NilToEmptyString(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}

func EmptyToNilString(s *string) *string {
	if s != nil && *s == "" {
		return nil
	}

	return s
}

func NilToEmptyUri(s *Uri) Uri {
	if s == nil {
		return ""
	}

	return *s
}

func EmptyToNilUri(s *Uri) *Uri {
	if s != nil && *s == "" {
		return nil
	}

	return s
}
