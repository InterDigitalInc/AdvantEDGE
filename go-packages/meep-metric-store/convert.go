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

package metricstore

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

func JsonNumToInt64(num json.Number) (val int64) {
	if intVal, err := num.Int64(); err == nil {
		val = intVal
	}
	return val
}
