package meander_test

import (
	"testing"

	"github.com/matryer/goblueprints/chapter7/meander"
	"github.com/stretchr/testify/require"
)

type obj struct {
	value1 string
	value2 string
	value3 string
}

func (o *obj) Public() interface{} {
	return map[string]interface{}{"one": o.value1, "three": o.value3}
}

func TestPublic(t *testing.T) {

	o := &obj{
		value1: "value1",
		value2: "value2",
		value3: "value3",
	}

	v, ok := meander.Public(o).(map[string]interface{})
	require.True(t, ok, "Result should be msi")
	require.Equal(t, v["one"], "value1")
	require.Nil(t, v["two"], "value2")
	require.Equal(t, v["three"], "value3")

}
