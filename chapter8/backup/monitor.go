package backup

import (
	"fmt"
	"path/filepath"
	"time"
)

// Monitor checks paths and archives any that
// have changed.
type Monitor struct {
	Paths       map[string]string
	Archiver    Archiver
	Destination string
}

// Now checks all directories in Paths with the latest hash.
// Archive will be called for any paths whose hashes do not match.
func (m *Monitor) Now() (int, error) {
	var counter int
	for path, lastHash := range m.Paths {
		newHash, err := DirHash(path)
		if err != nil {
			return 0, err
		}
		if newHash != lastHash {
			err := m.act(path)
			if err != nil {
				return counter, err
			}
			m.Paths[path] = newHash // update the hash
			counter++
		}
	}
	return counter, nil
}

func (m *Monitor) act(path string) error {
	dirname := filepath.Base(path)
	filename := fmt.Sprintf(m.Archiver.DestFmt(), time.Now().UnixNano())
	return m.Archiver.Archive(path, filepath.Join(m.Destination, dirname, filename))
}
