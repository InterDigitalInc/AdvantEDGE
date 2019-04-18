package main

import (
	"errors"
	log "github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-tc-sidecar/log"
	"net"
	"reflect"
	"time"

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

// DBEntryExists - true if entry exists; false otherwise
func DBEntryExists(key string) bool {
	value := dbClient.Exists(key).Val()
	if value == 0 {
		return false
	}
	return true
}

// DBAddEntry - Add entry to DB
func DBAddEntry(key string, fields map[string]string) error {
	m := convertMapStrStrToMapStrInt(fields)
	_, err := dbClient.HMSet(key, m).Result()
	if err != nil {
		return err
	}
	return nil
}

func convertMapStrStrToMapStrInt(src map[string]string) (dst map[string]interface{}) {
	dst = make(map[string]interface{})
	for key, value := range src {
		dst[key] = value
	}
	return dst
}

// DBRemoveEntry - Remove entry from DB
func DBRemoveEntry(key string) error {
	_, err := dbClient.Del(key).Result()
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
