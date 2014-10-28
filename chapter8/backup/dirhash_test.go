package backup_test

import (
	"testing"

	"github.com/matryer/goblueprints/chapter8/backup"
	"github.com/stretchr/testify/require"
)

func TestDirHash(t *testing.T) {

	hash1a, err := backup.DirHash("test/hash1")
	require.NoError(t, err)
	hash1b, err := backup.DirHash("test/hash1")
	require.NoError(t, err)

	require.Equal(t, hash1a, hash1b, "hash1 and hash1b should be identical")

	hash2, err := backup.DirHash("test/hash2")
	require.NoError(t, err)

	require.NotEqual(t, hash1a, hash2, "hash1 and hash2 should not be the same")

}
