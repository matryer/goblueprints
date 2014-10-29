package backup_test

import (
	"strings"
	"testing"

	"github.com/matryer/goblueprints/chapter8/backup"
	"github.com/stretchr/testify/require"
)

func TestMonitor(t *testing.T) {

	a := &TestArchiver{}
	m := &backup.Monitor{
		Destination: "test/archive",
		Paths: map[string]string{
			"test/hash1": "abc",
			"test/hash2": "def",
		},
		Archiver: a,
	}

	n, err := m.Now()
	require.NoError(t, err)
	require.Equal(t, 2, n)

	require.Equal(t, 2, len(a.Archives))
	for _, call := range a.Archives {
		require.True(t, strings.HasPrefix(call.Dest, m.Destination))
		require.True(t, strings.HasSuffix(call.Dest, ".zip"))
	}

}
