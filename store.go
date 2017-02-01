package session

import (
	"net/http"
	"encoding/hex"
	"crypto/rand"
	"time"
	"log"
)

var expiration = time.Hour*48 // deffault expiration is 2 days
const n = 16
var store map[string]*Session

func InitStore(){
	store = make(map[string]*Session)
}

func Get(r *http.Request) *Session{
	cookie, err := r.Cookie("authid")
	if(err != nil){
		log.Println(err)
		return nil
	}
	key := cookie.Value
	log.Println("Got cookie: "+key)
	session := store[key]
	if(session == nil) {
		return nil
	}

	if(session.expires.Before(time.Now())){ // is expired
		log.Println("Session is expired.")
		delete(store, key)
		return nil
	}
	return session
}
func New(w http.ResponseWriter) *Session{
	key := ""
	exists := true;
	for exists{
		key = RandomString(n)
		if key == ""{
			return nil
		}
		_, exists = store[key]
	}
	log.Println("new cookie key: "+key)// DEBUG

	expTime := time.Now().Add(expiration)
    session := &Session{
    	values: make(map[string]interface{}),
    	expires: expTime}
    store[key] = session

    http.SetCookie(w ,&http.Cookie{
    	Name: "authid",
    	Value: key,
    	Expires: expTime})

    return session
}

func RandomString(len uint8) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return hex.EncodeToString(b)
}
