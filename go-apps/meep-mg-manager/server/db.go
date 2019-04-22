/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */
package server

import (
	"errors"
	"net"
	"reflect"
	"time"

	log "github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-mg-manager/log"

	"github.com/KromDaniel/rejonson"
	"github.com/go-redis/redis"
)

var dbClient *rejonson.Client
var dbClientStarted = false

var pubsub *redis.PubSub

// DBConnect - Establish connection to DB
func DBConnect() error {
	if dbClientStarted == false {
		err := openDB()
		if err != nil {
			return err
		}
	}
	return nil
}

func openDB() error {

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "meep-redis-master:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	rejonsonClient := rejonson.ExtendClient(redisClient)

	pong, err := rejonsonClient.Ping().Result()

	if pong == "" {
		log.Info("pong is null")
		return err
	}

	if err != nil {
		log.Info("Redis DB not accessible")
		return err
	}
	dbClientStarted = true
	dbClient = rejonsonClient
	log.Info("Redis DB opened and well!")
	return nil
}

// DBFlush - Empty DB
func DBFlush(module string) error {
	var cursor uint64
	var err error
	log.Debug("DBFlush module: ", module)

	// Find all module keys
	// Process in chunks of 50 matching entries to optimize processing speed & memory
	keyMatchStr := module + ":*"
	for {
		var keys []string
		keys, cursor, err = dbClient.Scan(cursor, keyMatchStr, 50).Result()
		if err != nil {
			log.Debug("ERROR: ", err)
			break
		}

		// Delete all matching entries
		if len(keys) > 0 {
			_, err = dbClient.Del(keys...).Result()
			if err != nil {
				log.Debug("Failed to retrieve entry fields")
				break
			}
		}

		// Stop searching if cursor is back at beginning
		if cursor == 0 {
			break
		}
	}

	return nil
}

// DBForEachEntry - Search for matching keys and run handler for each entry
func DBForEachEntry(keyMatchStr string, entryHandler func(string, map[string]string, interface{}) error, userData interface{}) error {
	var cursor uint64
	var err error

	// Process in chunks of 50 matching entries to optimize processing speed & memory
	for {
		var keys []string
		keys, cursor, err = dbClient.Scan(cursor, keyMatchStr, 50).Result()
		if err != nil {
			log.Debug("ERROR: ", err)
			break
		}

		if len(keys) > 0 {
			for i := 0; i < len(keys); i++ {
				fields, err := dbClient.HGetAll(keys[i]).Result()
				if err != nil || fields == nil {
					log.Debug("Failed to retrieve entry fields")
					break
				}

				// Invoke handler to process entry
				err = entryHandler(keys[i], fields, userData)
				if err != nil {
					return err
				}
			}
		}

		// Stop searching if cursor is back at beginning
		if cursor == 0 {
			break
		}
	}
	return nil
}

// DBSetEntry - Update existing entry or create new entry if it does not exist
func DBSetEntry(key string, fields map[string]interface{}) error {
	// Update existing entry or create new entry if it does not exist
	_, err := dbClient.HMSet(key, fields).Result()
	if err != nil {
		return err
	}
	return nil
}

// DBRemoveEntry - Remove entry from DB
func DBRemoveEntry(key string) error {
	_, err := dbClient.Del(key).Result()
	if err != nil {
		return err
	}
	return nil
}

// DBJsonGetEntry - Retrieve entry from DB
func DBJsonGetEntry(key string, path string) (string, error) {
	// Update existing entry or create new entry if it does not exist
	json, err := dbClient.JsonGet(key, path).Result()
	if err != nil {
		return "", err
	}
	return json, nil
}

// DBJsonSetEntry - Update existing entry or create new entry if it does not exist
func DBJsonSetEntry(key string, path string, json string) error {
	// Update existing entry or create new entry if it does not exist
	_, err := dbClient.JsonSet(key, path, json).Result()
	if err != nil {
		return err
	}
	return nil
}

// DBJsonDelEntry - Remove existing entry
func DBJsonDelEntry(key string, path string) error {
	_, err := dbClient.JsonDel(key, path).Result()
	if err != nil {
		return err
	}
	return nil
}

// Subscribe - Register as a listener for provided channels
func Subscribe(channels ...string) error {
	pubsub = dbClient.Subscribe(channels...)
	return nil
}

// Listen - Wait for subscribed events
func Listen(handler func(string, string)) error {

	// Make sure listener is subscribed to pubsub
	if pubsub == nil {
		return errors.New("Not subscribed to pubsub")
	}

	// Main listening loop
	for {
		// Wait for subscribed channel events, or timeout
		msg, err := pubsub.ReceiveTimeout(time.Second)
		if err != nil {
			if reflect.TypeOf(err) == reflect.TypeOf(&net.OpError{}) &&
				reflect.TypeOf(err.(*net.OpError).Err).String() == "*net.timeoutError" {
				// Timeout, ignore and wait for next event
				continue
			}
		}

		// Process published event
		switch m := msg.(type) {

		// Process Subscription
		case *redis.Subscription:
			log.Info("Subscription Message: ", m.Kind, " to channel ", m.Channel, ". Total subscriptions: ", m.Count)

		// Process received Message
		case *redis.Message:
			log.Info("MSG on ", m.Channel, ": ", m.Payload)
			handler(m.Channel, m.Payload)
		}
	}
}

// Publish - Publish message to channel
func Publish(channel string, message string) error {
	log.Info("Publish to channel: ", channel, " Message: ", message)
	_, err := dbClient.Publish(channel, message).Result()
	return err
}
