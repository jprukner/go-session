package session
// TODO benchmark mutex locking vs channel and comunication sync
import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"sync"
	"time"
)

var expiration time.Duration
var n uint64
var mutex *sync.Mutex
var store map[string]*Session

func InitStore(defaultExpiration time.Duration, sessionIdLength uint64) {
	store = make(map[string]*Session)
	expiration = defaultExpiration
	n = sessionIdLength / 2 // because one byte is represented by two hex digits. e.g 0xFF is 0b11111111
	mutex = &sync.Mutex{}
}

func Get(r *http.Request) *Session {
	cookie, err := r.Cookie("authid")
	if err != nil {
		log.Println(err)
		return nil
	}
	key := cookie.Value
	log.Println("Got cookie: " + key)

	mutex.Lock()
	session := store[key]
	mutex.Unlock()

	if session == nil {
		return nil
	}

	if session.expires.Before(time.Now()) { // is expired
		log.Println("Session is expired.")
		delete(store, key)
		return nil
	}
	return session
}
func New(w http.ResponseWriter) *Session {
	key := ""
	exists := true
	for exists {
		key = RandomString(n)
		if key == "" {
			return nil
		}
		mutex.Lock()
		_, exists = store[key]
		mutex.Unlock()
	}
	log.Println("new cookie key: " + key) // DEBUG

	expTime := time.Now().Add(expiration)
	session := &Session{
		values:  make(map[string]interface{}),
		expires: expTime}

	mutex.Lock()
	store[key] = session
	mutex.Unlock()

	http.SetCookie(w, &http.Cookie{
		Name:    "authid",
		Value:   key,
		Expires: expTime})

	return session
}

func RandomString(len uint64) string {
	b := make([]byte, len)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return hex.EncodeToString(b)
}
