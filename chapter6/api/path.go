package main

import (
	"strings"
)

// PathSeparator is the character used to separate
// HTTP paths.
const PathSeparator = "/"

// Path represents the path of a request.
type Path struct {
	Path string
	ID   string
}

// NewPath makes a new Path from the specified
// path string.
func NewPath(p string) *Path {
	var id string
	p = strings.Trim(p, PathSeparator)
	s := strings.Split(p, PathSeparator)
	if len(s) > 1 {
		id = s[len(s)-1]
		p = strings.Join(s[:len(s)-1], PathSeparator)
	}
	return &Path{Path: p, ID: id}
}

// HasID gets whether this path has an ID
// or not.
func (p *Path) HasID() bool {
	return len(p.ID) > 0
}
