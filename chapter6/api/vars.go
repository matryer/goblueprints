package main

import (
	"net/http"
	"sync"
)

var vars map[*http.Request]map[string]interface{}
var varsLock sync.RWMutex

// GetVar gets the value of the key for the specified request.
func GetVar(r *http.Request, key string) interface{} {
	var value interface{}
	varsLock.RLock()
	value = vars[r][key]
	varsLock.RUnlock()
	return value
}

// SetVar sets the key to the value for the specified request.
func SetVar(r *http.Request, key string, value interface{}) {
	varsLock.Lock()
	if vars == nil {
		vars = map[*http.Request]map[string]interface{}{}
	}
	if vars[r] == nil {
		vars[r] = map[string]interface{}{}
	}
	vars[r][key] = value
	varsLock.Unlock()
}
