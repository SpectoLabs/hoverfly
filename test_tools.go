package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"testing"

	"github.com/garyburd/redigo/redis"
)

// Client structure to be injected into functions to perform HTTP calls
type Client struct {
	HTTPClient *http.Client
}

func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func refute(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		t.Errorf("Did not expect %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func testTools(code int, body string) (*httptest.Server, *DBClient) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, body)
	}))

	tr := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}

	// getting redis configuration
	redisAddress := os.Getenv("RedisAddress")
	if redisAddress == "" {
		redisAddress = ":6379"
	}
	AppConfig.redisAddress = redisAddress

	redisPool := getRedisPool()

	cache := Cache{pool: redisPool, prefix: "genproxy_test:"}

	// preparing client
	dbClient := &DBClient{
		http:  &http.Client{Transport: tr},
		cache: cache,
	}
	return server, dbClient
}

// teardown does some cleanup after tests
func teardown() {
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.cache.pool.Close()

	// deleting cache

	client := dbClient.cache.pool.Get()
	defer client.Close()

	values, _ := redis.Strings(client.Do("KEYS", "genproxy_test:*"))

	for _, v := range values {
		fmt.Println(v)
		client.Do("DEL", v)
	}
}
