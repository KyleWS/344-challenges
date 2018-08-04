package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"github.com/info344-a17/challenges-KyleIWS/servers/gateway/handlers"
	"github.com/info344-a17/challenges-KyleIWS/servers/gateway/indexes"
	"github.com/info344-a17/challenges-KyleIWS/servers/gateway/models/users"
	"github.com/info344-a17/challenges-KyleIWS/servers/gateway/sessions"
	mgo "gopkg.in/mgo.v2"
)

func ServiceProxy(addrs []string, ctx *handlers.Ctx) *httputil.ReverseProxy {
	nextIndex := 0
	mx := sync.Mutex{}
	return &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			sess, errSess := sessions.GetSessionID(r, ctx.Key)
			r.Header.Del("X-User")
			if errSess == nil {
				sessionState := &handlers.SessionState{}
				ctx.SessionsStore.Get(sess, sessionState)
				if sessionState.AuthenticatedUser != nil {
					jsonVal, err := json.Marshal(sessionState.AuthenticatedUser)
					if err == nil {
						r.Header.Add("X-User", string(jsonVal))
					} else {
						fmt.Printf("error marshalling json: %v\n", err)
					}
				}
			} else {
				fmt.Printf("error when trying to set xuser: %v\n", errSess)
			}
			mx.Lock()
			r.URL.Host = addrs[nextIndex%len(addrs)]
			nextIndex++
			mx.Unlock()
			r.URL.Scheme = "http"
		},
	}
}

//main is the main entry point for the server
func main() {

	// ADDR is which port our server runs on
	// ex: 127.0.0.1(:ADDR)
	addr := os.Getenv("ADDR")
	if len(addr) == 0 {
		addr = ":443"
	}
	// TLS
	tlskey := os.Getenv("TLSKEY")
	tlscert := os.Getenv("TLSCERT")
	if len(tlskey) == 0 || len(tlscert) == 0 {
		log.Fatal("Please set TLSKEY and TLSCERT\n")
	}
	// SESSIONKEY signs and validates our sessions
	sessionkey := os.Getenv("SESSIONKEY")
	if len(sessionkey) == 0 {
		// default key is md5 hash of "kyle is cool"
		sessionkey = "8b3f95a3bb29d578eb4544607856e4de"
	}
	// REDISADDR is the address at which exists our redis server?
	redisaddr := os.Getenv("REDISADDR")
	if len(redisaddr) == 0 {
		// default key is hash of "kyle is cool"
		redisaddr = "127.0.0.1:6379"
	}
	// DBADDR is the address at which exists our redis server?
	dbaddr := os.Getenv("DBADDR")
	if len(dbaddr) == 0 {
		// default key is hash of "kyle is cool"
		dbaddr = "127.0.0.1:27017"
	}
	redisClientInstance := redis.NewClient(&redis.Options{
		Addr: redisaddr,
	})
	sessionStoreInstance := sessions.NewRedisStore(redisClientInstance, time.Hour)
	sess, err := mgo.Dial(dbaddr)
	if err != nil {
		fmt.Printf("error connecting to db : %v\n", err)
	}
	// address(es) for summary microservice
	summarySvcAddr := os.Getenv("SUMMARYSVCADDR")
	if len(summarySvcAddr) == 0 {
		summarySvcAddr = "localhost:4001"
	}
	splitSummarySvcAddr := strings.Split(summarySvcAddr, ",")
	// address(es) for messaging microservice.
	messageSvcAddr := os.Getenv("MESSAGESVCADDR")
	if len(messageSvcAddr) == 0 {
		messageSvcAddr = "localhost:5000"
	}
	splitMessageSvcAddr := strings.Split(messageSvcAddr, ",")

	usersStoreInstance := users.NewMongoStore(sess, "website", "user")
	rootTrieNode := indexes.NewTrieNode(0, nil)
	usersStoreInstance.LoadExistingUsers(rootTrieNode)
	handlerMux := &handlers.Ctx{
		Key:           sessionkey,
		SessionsStore: sessionStoreInstance,
		UsersStore:    usersStoreInstance,
		RootTrieNode:  rootTrieNode,
	}

	masterMux := http.NewServeMux()
	masterMux.HandleFunc("/v1/users", handlerMux.UsersHandler)
	masterMux.HandleFunc("/v1/users/me", handlerMux.UsersMeHandler)
	masterMux.HandleFunc("/v1/sessions", handlerMux.SessionsHandler)
	masterMux.HandleFunc("/v1/sessions/mine", handlerMux.SessionsMineHandler)
	masterMux.Handle("/v1/summary", ServiceProxy(splitSummarySvcAddr, handlerMux))
	masterMux.Handle("/v1/messages/", ServiceProxy(splitMessageSvcAddr, handlerMux))
	masterMux.Handle("/v1/channels/", ServiceProxy(splitMessageSvcAddr, handlerMux))
	masterMuxCORS := &handlers.CORS{
		Handler: masterMux,
	}

	log.Printf("Server is started and listening for port %s!", addr)
	log.Fatal(http.ListenAndServeTLS(addr, tlscert, tlskey, masterMuxCORS))
}
