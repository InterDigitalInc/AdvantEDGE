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

package metrics

import (
	"encoding/json"
	"strconv"
)

func JsonNumToInt32(num json.Number) (val int32) {
	if intVal, err := strconv.Atoi(num.String()); err == nil {
		val = int32(intVal)
	}
	return val
}

func JsonNumToFloat64(num json.Number) (val float64) {
	if floatVal, err := num.Float64(); err == nil {
		val = floatVal
	}
	return val
}

func StrToInt32(str string) (val int32) {
	if intVal, err := strconv.Atoi(str); err == nil {
		val = int32(intVal)
	}
	return val
}

func StrToFloat64(str string) (val float64) {
	if floatVal, err := strconv.ParseFloat(str, 64); err == nil {
		val = floatVal
	}
	return val
}
