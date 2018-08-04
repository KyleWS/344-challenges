package sessions

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
)

//RedisStore represents a session.Store backed by redis.
type RedisStore struct {
	//Redis client used to talk to redis server.
	Client *redis.Client
	//Used for key expiry time on redis.
	SessionDuration time.Duration
}

//NewRedisStore constructs a new RedisStore
func NewRedisStore(client *redis.Client, sessionDuration time.Duration) *RedisStore {
	store := &RedisStore{
		Client:          client,
		SessionDuration: sessionDuration,
	}

	return store
}

//Store implementation

//Save saves the provided `sessionState` and associated SessionID to the store.
//The `sessionState` parameter is typically a pointer to a struct containing
//all the data you want to associated with the given SessionID.
func (rs *RedisStore) Save(sid SessionID, sessionState interface{}) error {
	jsonVal, err := json.Marshal(sessionState)
	if err != nil {
		return err
	}
	key := sid.getRedisKey()
	_, err1 := rs.Client.Set(key, jsonVal, rs.SessionDuration).Result()
	if err1 != nil {
		return err1
	}
	return nil
}

//Get populates `sessionState` with the data previously saved
//for the given SessionID
func (rs *RedisStore) Get(sid SessionID, sessionState interface{}) error {
	result, err := rs.Client.Get(sid.getRedisKey()).Result()
	if err != nil {
		return ErrStateNotFound
	}
	json.Unmarshal([]byte(result), &sessionState)
	rs.Client.Expire(sid.getRedisKey(), rs.SessionDuration)
	return nil
}

//Delete deletes all state data associated with the SessionID from the store.
func (rs *RedisStore) Delete(sid SessionID) error {
	_, err := rs.Client.Del(sid.getRedisKey()).Result()
	if err != nil {
		return err
	}
	return nil
}

//getRedisKey() returns the redis key to use for the SessionID
func (sid SessionID) getRedisKey() string {
	//convert the SessionID to a string and add the prefix "sid:" to keep
	//SessionID keys separate from other keys that might end up in this
	//redis instance
	return "sid:" + sid.String()
}
