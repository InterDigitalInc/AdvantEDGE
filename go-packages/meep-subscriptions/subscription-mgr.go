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
	"errors"
	"sync"
	"time"

	dkm "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-key-mgr"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	redis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis"
)

type ExpiredSubscriptionCb func(*Subscription)
type PeriodicSubscriptionCb func(*Subscription)
type TestNotificationCb func(*Subscription) error
type NewWebsocketCb func(*Subscription) (string, error)

type SubscriptionMgrCfg struct {
	Module         string
	Sandbox        string
	Mep            string
	Service        string
	Basekey        string
	MetricsEnabled bool
	ExpiredSubCb   ExpiredSubscriptionCb
	PeriodicSubCb  PeriodicSubscriptionCb
	TestNotifCb    TestNotificationCb
	NewWsCb        NewWebsocketCb
}

type SubscriptionMgr struct {
	cfg           *SubscriptionMgrCfg
	rc            *redis.Connector
	baseKey       string
	subscriptions map[string]*Subscription
	mutex         sync.Mutex
	ticker        *time.Ticker
}

const subRedisTable = 0
const periodicCounterPending = -1

// NewSubscriptionMgr - Create and initialize a Subscription Manager instance
func NewSubscriptionMgr(cfg *SubscriptionMgrCfg, addr string) (sm *SubscriptionMgr, err error) {

	// Create new Subscription Manager instance
	log.Info("Creating new Subscription Manager")
	sm = new(SubscriptionMgr)
	sm.cfg = cfg

	// Connect to Redis DB
	sm.rc, err = redis.NewConnector(addr, subRedisTable)
	if err != nil {
		log.Error("Failed connection to Subscription Manager redis DB. Error: ", err)
		return nil, err
	}
	log.Info("Connected to Subscription Manager Redis DB")

	// Get base store key
	if cfg.Basekey != "" {
		sm.baseKey = cfg.Basekey
	} else {
		sm.baseKey = dkm.GetKeyRoot(cfg.Sandbox) + cfg.Module + ":mep:" + cfg.Mep + ":"
	}

	// Initialize subscription cache from store
	var subList []*Subscription
	var subListPtr = &subList
	key := sm.baseKey + "sub:*:*"
	err = sm.rc.ForEachJSONEntry(key, populateSubList, subListPtr)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	sm.subscriptions = make(map[string]*Subscription)
	for _, sub := range *subListPtr {
		sm.subscriptions[sub.Cfg.Id] = sub
		log.Debug("id: ", sub.Cfg.Id, " sub: ", sub)
	}

	// Start ticker
	sm.ticker = time.NewTicker(time.Second)
	go func() {
		for range sm.ticker.C {
			sm.runTicker()
		}
	}()

	log.Info("Created Subscription Manager")
	return sm, nil
}

func (sm *SubscriptionMgr) CreateSubscription(cfg *SubscriptionCfg, jsonSubOrig string) (*Subscription, error) {
	// Validate params
	if cfg == nil {
		return nil, errors.New("Missing subscription config")
	}

	// Generate subscription ID if none provided
	if cfg.Id == "" {
		cfg.Id = sm.GenerateSubscriptionId()
	}

	// Create new subscription
	sub, err := newSubscription(cfg, jsonSubOrig)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Process new subscription
	err = sm.processSubscription(sub)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	// Store new subscription
	err = sm.storeSubscription(sub)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return sub, nil
}

func (sm *SubscriptionMgr) UpdateSubscription(sub *Subscription) error {
	// Validate params
	if sub == nil {
		return errors.New("Missing subscription")
	}

	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Update subscription
	err := sub.updateSubscription()
	if err != nil {
		log.Error(err.Error())
		return err
	}

	// Process updated subscription
	err = sm.processSubscription(sub)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	// Store updated subscription
	err = sm.storeSubscription(sub)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func (sm *SubscriptionMgr) SetSubscriptionJson(sub *Subscription, jsonSub string) error {
	// Validate params
	if sub == nil {
		return errors.New("Missing subscription")
	}

	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Set original JSON
	sub.JsonSubOrig = jsonSub

	// Store updated subscription
	err := sm.storeSubscription(sub)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func (sm *SubscriptionMgr) DeleteSubscription(sub *Subscription) error {
	// Validate params
	if sub == nil {
		return errors.New("Missing subscription")
	}

	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Delete subscription
	err := sm.delSubscription(sub)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func (sm *SubscriptionMgr) DeleteAllSubscriptions() error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Get subscriptions from cache
	subList := make([]*Subscription, 0, len(sm.subscriptions))
	for _, sub := range sm.subscriptions {
		subList = append(subList, sub)
	}

	// Delete subscriptions
	for _, sub := range subList {
		err := sm.delSubscription(sub)
		if err != nil {
			log.Error(err.Error())
		}
	}
	return nil
}

func (sm *SubscriptionMgr) DeleteFilteredSubscriptions(AppId string, Type string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Get filtered subscriptions from cache
	var subList []*Subscription
	for _, sub := range sm.subscriptions {
		if (AppId == "" || sub.Cfg.AppId == AppId) && (Type == "" || sub.Cfg.Type == Type) {
			subList = append(subList, sub)
		}
	}

	// Delete subscriptions
	for _, sub := range subList {
		err := sm.delSubscription(sub)
		if err != nil {
			log.Error(err.Error())
		}
	}
	return nil
}

func (sm *SubscriptionMgr) GetSubscription(Id string) (*Subscription, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Get subscription from cache
	sub, found := sm.subscriptions[Id]
	if !found {
		return nil, errors.New("Subscription ID not found")
	}
	return sub, nil
}

func (sm *SubscriptionMgr) GetAllSubscriptions() ([]*Subscription, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Get subscriptions from cache
	var subList []*Subscription
	for _, sub := range sm.subscriptions {
		subList = append(subList, sub)
	}
	return subList, nil
}

func (sm *SubscriptionMgr) GetFilteredSubscriptions(AppId string, Type string) ([]*Subscription, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Get filtered subscriptions from cache
	var subList []*Subscription
	for _, sub := range sm.subscriptions {
		if (AppId == "" || sub.Cfg.AppId == AppId) && (Type == "" || sub.Cfg.Type == Type) {
			subList = append(subList, sub)
		}
	}
	return subList, nil
}

func (sm *SubscriptionMgr) GenerateSubscriptionId() string {
	randomStr, _ := generateRand(12)
	// return uuid.New().String()
	return "sub-" + randomStr
}

func (sm *SubscriptionMgr) ReadyToSend(sub *Subscription) bool {
	if sub == nil {
		return false
	}

	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Subscription state
	if !sub.isReady() {
		return false
	}
	// Periodic interval
	if sub.Cfg.PeriodicInterval > 0 && sub.PeriodicCounter != periodicCounterPending {
		return false
	}
	return true
}

func (sm *SubscriptionMgr) SendNotification(sub *Subscription, notif []byte) error {
	// Validate params
	if sub == nil {
		return errors.New("Missing subscription")
	}

	// Send notification
	err := sub.sendNotification(notif, sm.cfg.Sandbox, sm.cfg.Service, sm.cfg.MetricsEnabled)
	if err != nil {
		log.Error(err.Error())
	}

	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Reset periodic counter if present
	if sub.PeriodicCounter == periodicCounterPending {
		sub.PeriodicCounter = sub.Cfg.PeriodicInterval
	}

	return err
}

func (sm *SubscriptionMgr) processSubscription(sub *Subscription) error {

	if sub.Mode == ModeWebsocket {
		// Create Websocket path handler
		if !sub.WsCreated {
			// Validate callback
			if sm.cfg.NewWsCb == nil {
				err := errors.New("Websockets not supported")
				log.Error(err.Error())
				return err
			}

			// Invoke callback to create websocket path
			wsUri, err := sm.cfg.NewWsCb(sub)
			if err != nil {
				log.Error(err.Error())
				return err
			}

			// Update Websocket URI & creation flag
			sub.Ws.Uri = wsUri
			sub.WsCreated = true
		}

	} else if sub.Mode == ModeDirect {
		// Send Test notification if necessary
		if sub.State == StateTestNotif && !sub.TestNotifSent && sub.Cfg.RequestTestNotif {
			// Validate callback
			if sm.cfg.TestNotifCb == nil {
				err := errors.New("Test notification not supported")
				log.Error(err.Error())
				return err
			}
			// Start goroutine to trigger test notification after subscription creation
			go func() {
				// Allow some time to complete subscription creation
				time.Sleep(100 * time.Millisecond)

				// Invoke callback to send test notification & wait for result
				err := sm.cfg.TestNotifCb(sub)

				sm.mutex.Lock()
				defer sm.mutex.Unlock()

				// Update subscription state according to test notification result
				if err != nil {
					sub.TestNotifSent = false
				} else {
					sub.State = StateReady
				}

				// Store updated subscription
				err = sm.storeSubscription(sub)
				if err != nil {
					log.Error(err.Error())
				}
			}()

			// Set flag indicating test notification was sent
			sub.TestNotifSent = true
		}
	}

	return nil
}

func (sm *SubscriptionMgr) delSubscription(sub *Subscription) error {
	// Delete subscription
	err := sub.deleteSubscription()
	if err != nil {
		log.Error(err.Error())
		return err
	}

	// Remove from store
	err = sm.rc.JSONDelEntry(sm.baseKey+"sub:"+sub.Cfg.Type+":"+sub.Cfg.Id, ".")
	if err != nil {
		log.Error(err.Error())
		return err
	}

	// Remove from cache
	delete(sm.subscriptions, sub.Cfg.Id)

	return nil
}

func (sm *SubscriptionMgr) storeSubscription(sub *Subscription) error {

	// Store subscription
	jsonSub, err := convertSubToJson(sub)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	key := sm.baseKey + "sub:" + sub.Cfg.Type + ":" + sub.Cfg.Id
	err = sm.rc.JSONSetEntry(key, ".", jsonSub)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	// Cache updated subscription
	sm.subscriptions[sub.Cfg.Id] = sub

	return nil
}

func (sm *SubscriptionMgr) runTicker() {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Check for expired subscriptions
	if sm.cfg.ExpiredSubCb != nil {
		var expiredSubList []*Subscription
		currentTime := time.Now()

		for _, sub := range sm.subscriptions {
			if sub.State == StateExpired {
				// Add to list of expired subscriptions
				expiredSubList = append(expiredSubList, sub)
			} else if sub.Cfg.ExpiryTime != nil && currentTime.After(*sub.Cfg.ExpiryTime) {
				// Set state to expired & invoke expiry callback
				sub.State = StateExpired

				log.Debug("Invoking expiry callback for sub: ", sub.Cfg.Id)
				go sm.cfg.ExpiredSubCb(sub)
			}
		}

		// Remove expired subscriptions from previous iteration
		for _, sub := range expiredSubList {
			_ = sm.delSubscription(sub)
		}
	}

	// Trigger periodic notifications
	if sm.cfg.PeriodicSubCb != nil {
		for _, sub := range sm.subscriptions {
			if sub.Cfg.PeriodicInterval > 0 {
				if sub.PeriodicCounter > 0 {
					sub.PeriodicCounter--
				}
				// If periodic interval is up, trigger notification if subscription is ready
				if sub.PeriodicCounter == 0 && sub.isReady() {
					// Set counter to -1; it will be reset when notification is sent
					sub.PeriodicCounter = periodicCounterPending

					// Invoke periodic callback
					log.Debug("Invoking periodic callback for sub: ", sub.Cfg.Id)
					go sm.cfg.PeriodicSubCb(sub)
				}
			}
		}
	}
}

func populateSubList(key string, jsonSub string, userData interface{}) error {
	// Get query params & userlist from user data
	subListPtr := userData.(*([]*Subscription))
	if subListPtr == nil {
		return errors.New("subList not found in userData")
	}

	// Retrieve user info from DB
	sub, err := convertJsonToSub(jsonSub)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	// Add subscription to list
	*subListPtr = append(*subListPtr, sub)
	return nil
}
