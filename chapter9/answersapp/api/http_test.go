package api

import (
	"net/http"
	"testing"
)

func TestPathParams(t *testing.T) {
	r, err := http.NewRequest("GET", "1/2/3/4/5", nil)
	if err != nil {
		t.Errorf("NewRequest: %s", err)
	}
	params := pathParams(r, "one/two/three/four")
	if len(params) != 4 {
		t.Errorf("expected 4 params but got %d: %v", len(params), params)
	}
	for k, v := range map[string]string{
		"one":   "1",
		"two":   "2",
		"three": "3",
		"four":  "4",
	} {
		if params[k] != v {
			t.Errorf("%s: %s != %s", k, params[k], v)
		}
	}
	params = pathParams(r, "one/two/three/four/five/six")
	if len(params) != 5 {
		t.Errorf("expected 5 params but got %d: %v", len(params), params)
	}
	for k, v := range map[string]string{
		"one":   "1",
		"two":   "2",
		"three": "3",
		"four":  "4",
		"five":  "5",
	} {
		if params[k] != v {
			t.Errorf("%s: %s != %s", k, params[k], v)
		}
	}
}
