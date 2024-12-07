package v1

import (
	"net/http"
	"sync"
)

var once sync.Once
var instance *http.Client

func GetClient() *http.Client {
	once.Do(func() {
		instance = &http.Client{}
	})
	return instance
}
