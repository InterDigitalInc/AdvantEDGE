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
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"

	"github.com/KromDaniel/rejonson"
	"github.com/go-redis/redis"
)

var dbClient *rejonson.Client
var dbClientStarted = false

// RedisDBConnect - Establish connection to DB
func RedisDBConnect() error {
	if !dbClientStarted {
		err := openDB()
		if err != nil {
			return err
		}
	}
	return nil
}

func openDB() error {
	goRedisClient := redis.NewClient(&redis.Options{
		Addr:     "meep-redis-master:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	client := rejonson.ExtendClient(goRedisClient)

	pong, err := client.Ping().Result()
	if pong == "" {
		log.Info("pong is null")
		return err
	}

	if err != nil {
		log.Info("Redis DB not accessible")
		return err
	}
	dbClientStarted = true
	dbClient = client

	log.Info("Redis DB opened and well!")
	return nil
}

// RedisDBFlush - Empty DB
func RedisDBFlush(module string) error {
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

// RedisDBSetEntry - Update existing entry or create new entry if it does not exist
func RedisDBSetEntry(key string, fields map[string]interface{}) error {
	// Update existing entry or create new entry if it does not exist
	_, err := dbClient.HMSet(key, fields).Result()
	if err != nil {
		return err
	}
	return nil
}

// RedisDBRemoveEntry - Remove existing entries
func RedisDBRemoveEntry(keys ...string) error {
	_, err := dbClient.Del(keys...).Result()
	if err != nil {
		return err
	}
	return nil
}

// RedisDBJsonSetEntry - Update existing entry or create new entry if it does not exist
func RedisDBJsonSetEntry(key string, path string, json string) error {
	// Update existing entry or create new entry if it does not exist
	_, err := dbClient.JsonSet(key, path, json).Result()
	if err != nil {
		return err
	}
	return nil
}

// RedisDBJsonDelEntry - Remove existing entry
func RedisDBJsonDelEntry(key string, path string) error {
	_, err := dbClient.JsonDel(key, path).Result()
	if err != nil {
		return err
	}
	return nil
}

// RedisDBForEachEntry - Search for matching keys and run handler for each entry
func RedisDBForEachEntry(keyMatchStr string, entryHandler func(string, map[string]string, interface{}) error, userData interface{}) error {
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

// RedisDBPublish - Publish message to channel
func RedisDBPublish(channel string, message string) error {
	log.Info("Publish to channel: ", channel, " Message: ", message)
	_, err := dbClient.Publish(channel, message).Result()
	return err
}
