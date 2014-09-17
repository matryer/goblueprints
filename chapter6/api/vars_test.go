package main_test

import (
	"net/http"
	"testing"

	"github.com/matryer/goblueprints/chapter6/api"
	"github.com/stretchr/testify/require"
)

func TestVars(t *testing.T) {

	// make two requests
	r1 := &http.Request{}
	r2 := &http.Request{}

	var1 := "var1"
	var2 := "var2"

	// open the vars for these requests
	main.OpenVars(r1)
	main.OpenVars(r2)

	// set a variable for both requests with same key
	main.SetVar(r1, "key", var1)
	main.SetVar(r2, "key", var2)

	require.Equal(t, var1, main.GetVar(r1, "key"))
	require.Equal(t, var2, main.GetVar(r2, "key"))

	// close the vars
	main.CloseVars(r1)
	main.CloseVars(r2)

}
