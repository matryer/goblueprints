package backup_test

import (
	"os"
	"testing"

	"github.com/matryer/goblueprints/chapter8/backup"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T) {
	os.MkdirAll("test/output", 0777)
}
func teardown(t *testing.T) {
	os.RemoveAll("test/output")
}

func TestZipArchive(t *testing.T) {
	setup(t)
	defer teardown(t)

	err := backup.Zip.Archive("test/hash1", "test/output/1.zip")
	require.NoError(t, err)

	// unzip
	err = backup.Zip.Restore("test/output/1.zip", "test/output/restored")
	require.NoError(t, err)

}
