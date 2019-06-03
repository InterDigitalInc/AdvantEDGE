/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */

package redis

import (
	log "github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-loc-serv/log"

	"errors"
	"net"
	"reflect"
	"time"

	"github.com/KromDaniel/rejonson"
	"github.com/go-redis/redis"
)

// RedisDBConnect - Establish connection to DB
func RedisDBConnect(db int) (*rejonson.Client, error) {

	dbClient, err := openDB(db)
	if err != nil {
		return nil, err
	}
	return dbClient, nil
}

func openDB(db int) (*rejonson.Client, error) {
	goRedisClient := redis.NewClient(&redis.Options{
		Addr:     "meep-redis-master:6379",
		Password: "", // no password set
		DB:       db, // 0 for use default DB
	})
	client := rejonson.ExtendClient(goRedisClient)

	pong, err := client.Ping().Result()
	if pong == "" {
		log.Info("pong is null")
		return nil, err
	}

	if err != nil {
		log.Info("Redis DB not accessible")
		return nil, err
	}

	return client, nil
}

// RedisDBFlush - Empty DB
func RedisDBFlush(dbClient *rejonson.Client, keyMatchStr string) error {
	var cursor uint64
	var err error

	// Find all module keys
	// Process in chunks of 50 matching entries to optimize processing speed & memory
	keyMatchStr = keyMatchStr + "*"
	for {
		var keys []string
		keys, cursor, err = dbClient.Scan(cursor, keyMatchStr, 50).Result()
		if err != nil {
			//log.Debug("ERROR: ", err)
			break
		}

		// Delete all matching entries
		if len(keys) > 0 {
			_, err = dbClient.Del(keys...).Result()
			if err != nil {
				//log.Debug("Failed to retrieve entry fields")
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

// DBJsonGetEntry - Retrieve entry from DB
func DBJsonGetEntry(dbClient *rejonson.Client, key string, path string) (string, error) {
	// Update existing entry or create new entry if it does not exist
	json, err := dbClient.JsonGet(key, path).Result()
	if err != nil {
		return "", err
	}
	return json, nil
}

// RedisDBSetEntry - Update existing entry or create new entry if it does not exist
func RedisDBSetEntry(dbClient *rejonson.Client, key string, fields map[string]interface{}) error {
	// Update existing entry or create new entry if it does not exist
	_, err := dbClient.HMSet(key, fields).Result()
	if err != nil {
		return err
	}
	return nil
}

// RedisDBRemoveEntry - Remove existing entries
func RedisDBRemoveEntry(dbClient *rejonson.Client, keys ...string) error {
	_, err := dbClient.Del(keys...).Result()
	if err != nil {
		return err
	}
	return nil
}

// RedisDBJsonSetEntry - Update existing entry or create new entry if it does not exist
func DBJsonSetEntry(dbClient *rejonson.Client, key string, path string, json string) error {
	// Update existing entry or create new entry if it does not exist
	_, err := dbClient.JsonSet(key, path, json).Result()
	if err != nil {
		return err
	}
	return nil
}

// RedisDBJsonDelEntry - Remove existing entry
func DBJsonDelEntry(dbClient *rejonson.Client, key string, path string) error {
	_, err := dbClient.JsonDel(key, path).Result()
	if err != nil {
		return err
	}
	return nil
}

// RedisDBForEachEntry - Search for matching keys and run handler for each entry
func RedisDBForEachEntry(dbClient *rejonson.Client, keyMatchStr string, entryHandler func(string, map[string]string, interface{}) error, userData interface{}) error {
	var cursor uint64
	var err error

	// Process in chunks of 50 matching entries to optimize processing speed & memory
	for {
		var keys []string
		keys, cursor, err = dbClient.Scan(cursor, keyMatchStr, 50).Result()
		if err != nil {
			//log.Debug("ERROR: ", err)
			break
		}

		if len(keys) > 0 {
			for i := 0; i < len(keys); i++ {
				fields, err := dbClient.HGetAll(keys[i]).Result()
				if err != nil || fields == nil {
					//log.Debug("Failed to retrieve entry fields")
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

// RedisDBForEachEntry - Search for matching keys and run handler for each entry
func RedisDBForEachJsonEntry(dbClient *rejonson.Client, keyMatchStr string, param1 string, param2 string, entryHandler func(string, string, string, string, interface{}) error, userData interface{}) error {
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
				jsonInfo, err := dbClient.JsonGet(keys[i], ".").Result()
				if err != nil || jsonInfo == "" {
					log.Debug("Failed to retrieve entry fields")
					break
				}

				// Invoke handler to process entry
				err = entryHandler(keys[i], jsonInfo, param1, param2, userData)
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

// Subscribe - Register as a listener for provided channels
func Subscribe(dbClient *rejonson.Client, channels ...string) (*redis.PubSub, error) {
	pubsub := dbClient.Subscribe(channels...)
	return pubsub, nil
}

// Listen - Wait for subscribed events
func Listen(pubsub *redis.PubSub, handler func(string, string)) error {

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
			//log.Info("Subscription Message: ", m.Kind, " to channel ", m.Channel, ". Total subscriptions: ", m.Count)

		// Process received Message
		case *redis.Message:
			//log.Info("MSG on ", m.Channel, ": ", m.Payload)
			handler(m.Channel, m.Payload)
		}
	}
}

// RedisDBPublish - Publish message to channel
func RedisDBPublish(dbClient *rejonson.Client, channel string, message string) error {
	//log.Info("Publish to channel: ", channel, " Message: ", message)
	_, err := dbClient.Publish(channel, message).Result()
	return err
}

func DbJsonGet(dbClient *rejonson.Client, resourceName string, elementPath string) string {

	if dbClient == nil {
		err := errors.New("Database client is nil")
		log.Error(err)
		return ""
	}

	jsonInfo, err := DBJsonGetEntry(dbClient, elementPath+":"+resourceName, ".")
	if err != nil {
		log.Error(err.Error())
		return ""
	}

	return jsonInfo

}

func DbJsonGetList(dbClient *rejonson.Client, elem1 string, elem2 string, elementPath string, entryHandler func(string, string, string, string, interface{}) error, dataList interface{}) error {

	if dbClient == nil {
		err := errors.New("Database client is nil")
		log.Error(err)
		return err
	}

	keyName := elementPath + "*"
	err := RedisDBForEachJsonEntry(dbClient, keyName, elem1, elem2, entryHandler, dataList)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func DbJsonSet(dbClient *rejonson.Client, name string, jsonInfo string, elementPath string) error {

	if dbClient == nil {
		err := errors.New("Database client is nil")
		log.Error(err)
		return err
	}

	err := DBJsonSetEntry(dbClient, elementPath+":"+name, ".", jsonInfo)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func DbJsonDelete(dbClient *rejonson.Client, resourceName string, elementPath string) error {

	if dbClient == nil {
		err := errors.New("Database client is nil")
		log.Error(err)
		return err
	}

	err := DBJsonDelEntry(dbClient, elementPath+":"+resourceName, ".")
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}
