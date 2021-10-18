/*
 * Copyright (c) 2021  InterDigital Communications, Inc
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

package subscriptions

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
)

// Generate a random string
func generateRand(n int) (string, error) {
	data := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, data); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(data), nil
}

func convertSubToJson(sub *Subscription) (string, error) {
	jsonSub, err := json.Marshal(sub)
	if err != nil {
		log.Error(err.Error())
		return "", err
	}
	return string(jsonSub), nil
}

func convertJsonToSub(jsonSub string) (*Subscription, error) {
	var sub Subscription
	err := json.Unmarshal([]byte(jsonSub), &sub)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return &sub, nil
}
