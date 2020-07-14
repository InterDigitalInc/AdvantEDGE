/*
 * Copyright (c) 2020  InterDigital Communications, Inc
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

package sessionsManager

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/rs/xid"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	redis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis"
)

type User struct {
	Username  string
	Password  string
	SessionId string
	Active    bool
}

var user1 = User{"u1", "1234", "NA", false}
var user2 = User{"u2", "2345", "NA", false}
var user3 = User{"u3", "3456", "NA", false}

// Map of configured users - Key=Username
var ConfiguredUsers map[string]*User

// Map of connected users - key=SessionId
var ConnectedUsers map[string]*User

// store will hold all session data
var CookieStore *sessions.CookieStore

var redisDBAddr string = "meep-redis-master.default.svc.cluster.local:6379"
var redisClient *redis.Connector

const sessionsDb = 0
const sessionsStoreKey = "session-store:"
const sessionsKey = "keys:"

//const encryptionKey = "encryptionKey:"
const sessionIdKey = "sid:"
const authKeyStr = "authKey"
const encryptionKeyStr = "encryptionKey"

func Init(redisAddr string) error {
	log.Info("Initializing sessions manager")

	//CookieStore = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
	//CookieStore = sessions.NewCookieStore([]byte("veryprivatekey"))

	if redisAddr == "" {
		redisAddr = redisDBAddr
	}
	// Connect to Redis DB
	var err error
	redisClient, err = redis.NewConnector(redisAddr, sessionsDb)
	if err != nil {
		log.Error("Failed connection to Sessions redis DB. Error: ", err)
		return err
	}
	log.Info("Connected to Metrics Redis DB")

	authKey, encryptionKey, err := getRedisSessionKeys()
	if err != nil {
		log.Error("Failed to get session manager keys: ", err)
		return err
	}

	sessionKeys := map[string]interface{}{authKeyStr: authKey, encryptionKeyStr: encryptionKey}

	if authKey == "" || encryptionKey == "" {
		if authKey == "" {
			authKey = string(securecookie.GenerateRandomKey(64))
			sessionKeys[authKeyStr] = authKey
		}
		if encryptionKey == "" {
			encryptionKey = string(securecookie.GenerateRandomKey(32))
			sessionKeys[encryptionKeyStr] = encryptionKey
		}
		setRedisSessionKeys(sessionKeys)
	}

	CookieStore = sessions.NewCookieStore(
		[]byte(authKey),
		[]byte(encryptionKey),
	)

	CookieStore.Options = &sessions.Options{
		MaxAge:   60 * 15,
		HttpOnly: true,
		Domain:   "10.3.16.150",
	}
	// uncomment when ready to test with HTTPS
	//Secure: true
	// Not sure SameSite is supported by Gorilla Sessions
	//SameSite --

	// add preconfigured users
	ConfiguredUsers = make(map[string]*User)
	ConfiguredUsers[user1.Username] = &user1
	ConfiguredUsers[user2.Username] = &user2
	ConfiguredUsers[user3.Username] = &user3

	// no active sessions
	ConnectedUsers = make(map[string]*User)

	return nil
	//PrintUsers()

}

// getRedisSessionId - Generic session getter
func getRedisSessionId(keySuffix string) (map[string]string, error) {

	key := sessionsStoreKey + sessionIdKey + keySuffix
	user, err := redisClient.GetEntry(key)
	if err != nil {
		log.Error("Failed to set entry: ", err)
		return nil, err
	}
	return user, nil
}

// setRedisSessionId - Generic session id setter
func setRedisSessionId(sessionId string, user map[string]interface{}) error {

	key := sessionsStoreKey + sessionIdKey + sessionId
	err := redisClient.SetEntry(key, user)
	if err != nil {
		log.Error("Failed to set entry: ", err)
		return err
	}
	return nil
}

// delRedisSessionId - Generic session delete
func delRedisSessionId(sessionId string) error {

	key := sessionsStoreKey + sessionIdKey + sessionId
	log.Info("SIMON delete ", key)
	err := redisClient.DelEntry(key)
	if err != nil {
		log.Error("Failed to delete entry: ", err)
		return err
	}
	return nil
}

// getRedisSessionKeys - Generic session getter
func getRedisSessionKeys() (string, string, error) {

	key := sessionsStoreKey + sessionsKey
	sessionKeys, err := redisClient.GetEntry(key)
	if err != nil {
		log.Error("Failed to get entry: ", err)
		return "", "", err
	}
	return sessionKeys[authKeyStr], sessionKeys[encryptionKeyStr], nil
}

// setRedisSessionKeyValue - Generic session setter
func setRedisSessionKeys(keys map[string]interface{}) error {

	key := sessionsStoreKey + sessionsKey
	err := redisClient.SetEntry(key, keys)
	if err != nil {
		log.Error("Failed to set entry: ", err)
		return err
	}
	return nil
}

// delRedisSessionKeys - Generic session keys delete
func delRedisSessionKeys(keySuffix string) error {

	key := sessionsStoreKey + sessionsKey
	err := redisClient.DelEntry(key)
	if err != nil {
		log.Error("Failed to delete entry: ", err)
		return err
	}
	return nil
}

func IsActiveSession(r *http.Request) bool {

	log.Info("----- IS ACTIVE -----")
	// Get session cookie
	session, err := CookieStore.Get(r, "authCookie")
	if err != nil {
		log.Info("Invalid session")
		return false
	}
	PrintSession(session)

	sid, _ := session.Values["SessionId"].(string)
	log.Info("SIMON sid ", sid)
	PrintConnectedUsers()
	_, ok := ConnectedUsers[sid]
	if !ok {
		log.Info("Invalid sid")
		//        return false
	}

	_, err = getRedisSessionId(sid)
	if err != nil {
		log.Info("Session not active")
		return false
	}
	return true
}

func deactivateUser(user *User) {
	log.Info("deactivate")
	// DB management
	_, ok := ConnectedUsers[user.SessionId]
	if ok {
		delete(ConnectedUsers, user.SessionId)
	}
	user.Active = false
	user.SessionId = "NA"
}

func AuthenticateNewUser(username string, password string, w http.ResponseWriter, r *http.Request) error {

	//PrintConnectedUsers()
	// Get session cookie
	session, err := CookieStore.Get(r, "authCookie")
	if err != nil {
		log.Info(err.Error())
		// Patch until encryption keys are persisted
		// Try to renew the cookie
		if session != nil {
			session, err = CookieStore.New(r, "authCookie")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return err
			}
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return err
		}
	}
	PrintSession(session)

	// Authenticate user
	// Does user exist?
	user, ok := ConfiguredUsers[username]
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return err
	}
	// Does user have an existing session
	if user.Active {
		deactivateUser(user)
	}
	// Does password match?
	if user.Password != password {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return err
	}

	// user exists & password match - Let's get it done
	//        user.Active = true
	user.SessionId = xid.New().String()
	ConnectedUsers[user.SessionId] = user
	session.Values["SessionId"] = user.SessionId
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	log.Info("SIMON session ", session.Values["SessionId"], "---", session)
	PrintSession(session)

	//valid session Id
	userMap := map[string]interface{}{"username": user.Username}
	setRedisSessionId(user.SessionId, userMap)
	PrintConnectedUsers()

	return nil
}

func AuthenticateUserDeletion(w http.ResponseWriter, r *http.Request) error {

	session, err := CookieStore.Get(r, "authCookie")
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return err
	}
	PrintSession(session)

	sid, _ := session.Values["SessionId"].(string)
	log.Info("SIMON delete sid", sid)
	PrintConnectedUsers()
	user, ok := ConnectedUsers[sid]
	if !ok {
		http.Error(w, "Invalid session id", http.StatusUnauthorized)
		return err

	}

	delRedisSessionId(sid)

	//log.Info("SIMON 2")
	deactivateUser(user)
	// Invalidate session
	session.Values["SessionId"] = ""
	session.Options.MaxAge = -1
	err = session.Save(r, w)
	if err != nil {
		log.Info("SIMON 3")
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return err
	}
	log.Info("SIMON 4")
	PrintConnectedUsers()

	return nil
}

func PrintSession(session *sessions.Session) {
	//str,_ := json.Marshal(session.Values)
	str, ok := session.Values["SessionId"].(string)
	if !ok {
		str = "nil"
	}
	log.Info("Session: " + str)
}

func PrintUser(msg string, user *User) {
	str, _ := json.Marshal(*user)
	log.Info(msg + " " + string(str))
}

func PrintUsers() {
	str, _ := json.Marshal(ConfiguredUsers)
	log.Info("ConfiguredUsers:" + string(str))
	str, _ = json.Marshal(ConnectedUsers)
	log.Info("ConnectedUsers:" + string(str))
}

func PrintConnectedUsers() {
	str, _ := json.Marshal(ConnectedUsers)
	log.Info("ConnectedUsers:" + string(str))
}
