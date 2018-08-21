package adtracker

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/cocoonlife/timber"
	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"
)

// Storer is an interface which satisfies the datastore requirements of
// AdTracker
type Storer interface {
	Get(string) (int, error)
	Increment(string) error
}

// BasicStore is a basic in memory threadsafe store without persistence
// BasicStore allows us to run the AdTracker with minimal dependencies
// TODO: Replace BasicStore with usage of github.com/rafaeljusto/redigomock
type BasicStore struct {
	m map[string]int
	sync.Mutex
}

// Get retieves a value given a key
func (bc *BasicStore) Get(key string) (int, error) {
	bc.Lock()
	defer bc.Unlock()
	v, ok := bc.m[key]
	var err error
	if !ok {
		err = fmt.Errorf("ID %s does not exist", key)
	}
	return v, err
}

// Increment increments the value stored under the given key
func (bc *BasicStore) Increment(key string) error {
	bc.Lock()
	defer bc.Unlock()
	bc.m[key]++
	return nil
}

// NewBasicStore returns a new instance of BasicStore
func NewBasicStore() *BasicStore {
	return &BasicStore{m: make(map[string]int)}
}

// RedisStorage saves and loads files to Redis, implements Storer interface
type RedisStorage struct {
	prefix    string
	redisAddr string
	pool      *redis.Pool
}

// NewRedisStorage constructs a RedisStorage.
// We maintain a long lived connection.
// Every operation first checks to see if we have a client (and does nothing if we do)
// If we see any error, we close and release the client, which will cause the next operation
// to attempt a new connection.
// No retry logic is currently implemented.
func NewRedisStorage(redisAddr string, prefix string) *RedisStorage {
	rs := RedisStorage{
		prefix:    prefix,
		redisAddr: redisAddr,
		pool: &redis.Pool{
			MaxIdle:     3,
			IdleTimeout: 240 * time.Second,
			Dial: func() (redis.Conn, error) {
				return redis.Dial("tcp", redisAddr)
			},
		},
	}
	return &rs
}

func (rs *RedisStorage) makeKey(name string) string {
	key := rs.prefix + "/" + name
	return key
}

// Get fetches a value from redis
func (rs *RedisStorage) Get(key string) (int, error) {
	key = rs.makeKey(key)
	conn := rs.pool.Get()
	defer conn.Close()

	val, err := redis.String(conn.Do("GET", key))
	if err != nil {
		return 0, fmt.Errorf("Can't get [%s] from redis: %s", key, err)
	}
	i, err := strconv.Atoi(val)
	return i, err
}

// Increment stores a value to redis
func (rs *RedisStorage) Increment(key string) error {
	key = rs.makeKey(key)
	conn := rs.pool.Get()
	defer conn.Close()

	_, err := conn.Do("INCR", key)
	if err != nil {
		err = fmt.Errorf("Can't incr [%s] via redis: %s", key, err.Error())
	}
	return err
}

// This struct is easily marshallable as json and the definitions serve as
// documentation of the request and response bodies used by the server
// Using json allows us to easily extend the api in the future without breaking
// clients. We use json tags so api is not tightly coupled to variable names.

// Resp represents the data format of the response to HTTP GET requests
type Resp struct {
	Value int `json:"value"`
}

// App is an implementation of an adtracking server
// It would be nicer to use the same endpoint for both get and increment ops
// but different HTTP verbs.
type App struct {
	store Storer
}

// incHandler given an id, retrieves, increments and persists the
// value stored under that ID - endpoint /ad
func (a *App) incHandler(w http.ResponseWriter, r *http.Request) {
	a.incHandlerVars(w, r, mux.Vars(r))
}

func (a *App) incHandlerVars(w http.ResponseWriter, r *http.Request, vars map[string]string) {
	id := vars["id"]
	timber.Debugf("increment %s", id)

	err := a.store.Increment(id)
	if err != nil {
		timber.Errorf(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *App) getHandlerVars(w http.ResponseWriter, r *http.Request, vars map[string]string) {
	id := vars["id"]

	// Retrieve stored value for given id
	val, err := a.store.Get(id)

	timber.Debugf("get %s = %d", id, val)

	// Handle unknown ID
	if err != nil {
		timber.Errorf(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Marshal response json
	resp := Resp{val}
	respBytes, err := json.Marshal(resp)
	if err != nil {
		timber.Errorf(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write retrieved value to response
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(respBytes)
	if err != nil {
		timber.Errorf(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// getHandler given an id, retrieves and returns the associated
// stored value - endpoint /ad_count
func (a *App) getHandler(w http.ResponseWriter, r *http.Request) {
	a.getHandlerVars(w, r, mux.Vars(r))
}

// Run runs the server on the specified port
func Run(port int) {

	h := App{store: NewRedisStorage(":6379", "adtracker")}

	r := mux.NewRouter()
	r.HandleFunc("/ad/{id}", h.incHandler).Methods(http.MethodPut)
	r.HandleFunc("/ad/{id}", h.getHandler).Methods(http.MethodGet)

	s := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: r,

		// timeouts to prevent slow clients from holiding up a connection
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Fatal(s.ListenAndServe())
}
