package main

import (
	log "github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-mon-engine/log"

	"github.com/go-redis/redis"
)

var dbClient *redis.Client
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
	db := redis.NewClient(&redis.Options{
		Addr:     "meep-redis-master:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := db.Ping().Result()

	if pong == "" {
		log.Info("pong is null")
		return err
	}

	if err != nil {
		log.Info("Redis DB not accessible")
		return err
	}
	dbClientStarted = true
	dbClient = db

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

// DBSetEntry - Update existing entry or create new entry if it does not exist
func DBSetEntry(key string, fields map[string]interface{}) error {
	// Update existing entry or create new entry if it does not exist
	_, err := dbClient.HMSet(key, fields).Result()
	if err != nil {
		return err
	}
	return nil
}

// DBRemoveEntry - Remove existing entries
func DBRemoveEntry(keys ...string) error {
	_, err := dbClient.Del(keys...).Result()
	if err != nil {
		return err
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
