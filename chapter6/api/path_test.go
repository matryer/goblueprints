package main_test

import (
	"testing"

	"github.com/matryer/goblueprints/chapter6/api"
	"github.com/stretchr/testify/require"
)

func TestPath(t *testing.T) {

	var p *main.Path

	p = main.NewPath("/polls")
	require.False(t, p.HasID())
	require.Equal(t, "polls", p.Path)

	p = main.NewPath("/polls/")
	require.False(t, p.HasID())
	require.Equal(t, "polls", p.Path)

	p = main.NewPath("polls/")
	require.False(t, p.HasID())
	require.Equal(t, "polls", p.Path)

	p = main.NewPath("polls")
	require.False(t, p.HasID())
	require.Equal(t, "polls", p.Path)

	p = main.NewPath("/polls/1")
	require.True(t, p.HasID())
	require.Equal(t, "polls", p.Path)
	require.Equal(t, "1", p.ID)

	p = main.NewPath("/polls/1/")
	require.True(t, p.HasID())
	require.Equal(t, "polls", p.Path)
	require.Equal(t, "1", p.ID)

	p = main.NewPath("polls/1/")
	require.True(t, p.HasID())
	require.Equal(t, "polls", p.Path)
	require.Equal(t, "1", p.ID)

	p = main.NewPath("polls/1")
	require.True(t, p.HasID())
	require.Equal(t, "polls", p.Path)
	require.Equal(t, "1", p.ID)

}
