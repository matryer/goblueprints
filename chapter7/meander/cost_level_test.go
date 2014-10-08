package meander_test

import (
	"github.com/matryer/goblueprints/chapter7/meander"
	"github.com/stretchr/testify/require"

	"testing"
)

func TestCostValues(t *testing.T) {

	require.Equal(t, int(meander.Cost1), 1)
	require.Equal(t, int(meander.Cost2), 2)
	require.Equal(t, int(meander.Cost3), 3)
	require.Equal(t, int(meander.Cost4), 4)
	require.Equal(t, int(meander.Cost5), 5)

}

func TestCostString(t *testing.T) {

	require.Equal(t, meander.Cost1.String(), "$")
	require.Equal(t, meander.Cost2.String(), "$$")
	require.Equal(t, meander.Cost3.String(), "$$$")
	require.Equal(t, meander.Cost4.String(), "$$$$")
	require.Equal(t, meander.Cost5.String(), "$$$$$")

}

func TestParseCost(t *testing.T) {

	require.Equal(t, meander.Cost1, meander.ParseCost("$"))
	require.Equal(t, meander.Cost2, meander.ParseCost("$$"))
	require.Equal(t, meander.Cost3, meander.ParseCost("$$$"))
	require.Equal(t, meander.Cost4, meander.ParseCost("$$$$"))
	require.Equal(t, meander.Cost5, meander.ParseCost("$$$$$"))

}

func TestParseCostRange(t *testing.T) {

	var l *meander.CostRange
	l = meander.ParseCostRange("$$...$$$")
	require.Equal(t, l.From, meander.Cost2)
	require.Equal(t, l.To, meander.Cost3)

	l = meander.ParseCostRange("$...$$$$$")
	require.Equal(t, l.From, meander.Cost1)
	require.Equal(t, l.To, meander.Cost5)

}

func TestCostRangeString(t *testing.T) {

	require.Equal(t, "$$...$$$$", (&meander.CostRange{From: meander.Cost2, To: meander.Cost4}).String())

}
