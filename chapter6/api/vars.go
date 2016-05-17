package main

import (
	"net/http"
	"sync"
)

var l sync.RWMutex // protects vars
var vars map[*http.Request]map[string]interface{}

// GetVar gets the value of the key for the specified http.Request.
func GetVar(r *http.Request, key string) interface{} {
	l.RLock()
	value := vars[r][key]
	l.RUnlock()
	return value
}

// SetVar sets the key to the value for the specified http.Request.
func SetVar(r *http.Request, key string, value interface{}) {
	l.Lock()
	vars[r][key] = value
	l.Unlock()
}

// OpenVars opens the vars for the specified http.Request.
// Must be called before GetVar or SetVar is called for each
// request.
func OpenVars(r *http.Request) {
	l.Lock()
	if vars == nil {
		vars = make(map[*http.Request]map[string]interface{})
	}
	vars[r] = map[string]interface{}{}
	l.Unlock()
}

// CloseVars closes the vars for the specified
// http.Request.
// Must be called when all var activity is completed to
// clean up any used memory.
func CloseVars(r *http.Request) {
	l.Lock()
	delete(vars, r)
	l.Unlock()
}
