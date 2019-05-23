/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */

package redisdb

import (
	"errors"
	"strings"
	"time"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"

	"github.com/KromDaniel/rejonson"
	"github.com/go-redis/redis"
)

var dbClient *rejonson.Client
var dbClientStarted = false

var pubsub *redis.PubSub

// Connector - Implements a Redis connector
type Connector struct {
	addr          string
	connected     bool
	client        *rejonson.Client
	pubsub        *redis.PubSub
	isListening   bool
	doneListening chan bool
}

// NewConnector - Creates and initialize a Redis connector
func NewConnector(addr string) (rc *Connector, err error) {
	rc = new(Connector)
	err = rc.connectDB(addr)
	if err != nil {
		return nil, err
	}

	return rc, nil
}

func (rc *Connector) connectDB(addr string) error {
	if addr == "" {
		rc.addr = "meep-redis-master:6379"
	} else {
		rc.addr = addr
	}
	log.Debug("Redis Connector connecting to ", rc.addr)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     rc.addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	rc.client = rejonson.ExtendClient(redisClient)

	pong, err := rc.client.Ping().Result()

	if pong == "" || err != nil {
		log.Error("Redis Connector unable to connect ", rc.addr)
		return err
	}

	rc.connected = true
	log.Info("Redis Connector connected to ", rc.addr)
	return nil
}

// // DBFlush - Empty DB
// func DBFlush(module string) error {
// 	var cursor uint64
// 	var err error
// 	log.Debug("DBFlush module: ", module)

// 	// Find all module keys
// 	// Process in chunks of 50 matching entries to optimize processing speed & memory
// 	keyMatchStr := module + ":*"
// 	for {
// 		var keys []string
// 		keys, cursor, err = dbClient.Scan(cursor, keyMatchStr, 50).Result()
// 		if err != nil {
// 			log.Debug("ERROR: ", err)
// 			break
// 		}

// 		// Delete all matching entries
// 		if len(keys) > 0 {
// 			_, err = dbClient.Del(keys...).Result()
// 			if err != nil {
// 				log.Debug("Failed to retrieve entry fields")
// 				break
// 			}
// 		}

// 		// Stop searching if cursor is back at beginning
// 		if cursor == 0 {
// 			break
// 		}
// 	}

// 	return nil
// }

// // DBForEachEntry - Search for matching keys and run handler for each entry
// func DBForEachEntry(keyMatchStr string, entryHandler func(string, map[string]string, interface{}) error, userData interface{}) error {
// 	var cursor uint64
// 	var err error

// 	// Process in chunks of 50 matching entries to optimize processing speed & memory
// 	for {
// 		var keys []string
// 		keys, cursor, err = dbClient.Scan(cursor, keyMatchStr, 50).Result()
// 		if err != nil {
// 			log.Debug("ERROR: ", err)
// 			break
// 		}

// 		if len(keys) > 0 {
// 			for i := 0; i < len(keys); i++ {
// 				fields, err := dbClient.HGetAll(keys[i]).Result()
// 				if err != nil || fields == nil {
// 					log.Debug("Failed to retrieve entry fields")
// 					break
// 				}

// 				// Invoke handler to process entry
// 				err = entryHandler(keys[i], fields, userData)
// 				if err != nil {
// 					return err
// 				}
// 			}
// 		}

// 		// Stop searching if cursor is back at beginning
// 		if cursor == 0 {
// 			break
// 		}
// 	}
// 	return nil
// }

// // DBSetEntry - Update existing entry or create new entry if it does not exist
// func DBSetEntry(key string, fields map[string]interface{}) error {
// 	// Update existing entry or create new entry if it does not exist
// 	_, err := dbClient.HMSet(key, fields).Result()
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// // DBRemoveEntry - Remove entry from DB
// func DBRemoveEntry(key string) error {
// 	_, err := dbClient.Del(key).Result()
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// JSONGetEntry - Retrieve entry from DB
func (rc *Connector) JSONGetEntry(key string, path string) (string, error) {
	if !rc.connected {
		return "", errors.New("Redis Connector is disconnected (JSONGetEntry)")
	}
	// Update existing entry or create new entry if it does not exist
	json, err := rc.client.JsonGet(key, path).Result()
	if err != nil {
		return "", err
	}
	return json, nil
}

// // DBJsonSetEntry - Update existing entry or create new entry if it does not exist
// func DBJsonSetEntry(key string, path string, json string) error {
// 	// Update existing entry or create new entry if it does not exist
// 	_, err := dbClient.JsonSet(key, path, json).Result()
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// // DBJsonDelEntry - Remove existing entry
// func DBJsonDelEntry(key string, path string) error {
// 	_, err := dbClient.JsonDel(key, path).Result()
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// Subscribe - Register as a listener for provided channels
func (rc *Connector) Subscribe(channels ...string) error {
	if !rc.connected {
		return errors.New("Redis Connector is disconnected (Subscribe)")
	}

	rc.pubsub = rc.client.Subscribe(channels...)
	return nil
}

// Unsubscribe - Unregister as a listener for provided channels
func (rc *Connector) Unsubscribe(channels ...string) error {
	if !rc.connected {
		return errors.New("Redis Connector is disconnected (Unsubscribe)")
	}
	if rc.pubsub != nil {
		rc.pubsub.Unsubscribe(channels...)
	}
	return nil
}

// Listen - Wait for subscribed events
func (rc *Connector) Listen(handler func(string, string)) error {
	if !rc.connected {
		return errors.New("Redis Connector is disconnected (Listen)")
	}
	if rc.pubsub == nil {
		return errors.New("Not subscribed to pubsub (Listen)")
	}

	rc.isListening = true
	rc.doneListening = make(chan bool, 1)
	// Main listening loop
	for {
		// Wait for subscribed channel events, or timeout
		msg, err := rc.pubsub.ReceiveTimeout(time.Second)
		if err != nil {
			if !strings.Contains(err.Error(), "timeout") {
				log.Debug("Listen Error: ", err)
			}
		} else {
			channel := ""
			payload := ""

			// Process published event
			switch m := msg.(type) {
			// Process Subscription
			case *redis.Subscription:
				log.Info("Subscription Message: ", m.Kind, " to channel ", m.Channel, ". Total subscriptions: ", m.Count)
				continue
			// Process received Message
			case *redis.Message:
				channel = m.Channel
				payload = m.Payload
				log.Info("RX-MSG [", channel, "] ", payload)
				handler(channel, payload)
			}
		}

		if !rc.isListening {
			log.Debug("Redis Connector exiting listen routine")
			rc.doneListening <- true
			return nil
		}
	}
}

// StopListen - Stop the listening goroutine
func (rc *Connector) StopListen() {
	if rc.isListening {
		// stop the listen goroutine
		rc.isListening = false
		// synchronize on completion
		<-rc.doneListening
	}
}

// Publish - Publish message to channel
func (rc *Connector) Publish(channel string, message string) error {
	if !rc.connected {
		return errors.New("Redis Connector is disconnected (Publish)")
	}

	log.Info("TX-MSG [", channel, "] ", message)
	_, err := rc.client.Publish(channel, message).Result()
	return err
}
