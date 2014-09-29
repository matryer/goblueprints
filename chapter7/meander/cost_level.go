package meander

import (
	"strings"
)

type Cost int8

const (
	Cost1 Cost = iota
	Cost2
	Cost3
	Cost4
	Cost5
)

func (l Cost) String() string {
	return strings.Repeat("$", int(l)+1)
}

func ParseCost(s string) Cost {
	return Cost(len(s)) - 1
}

type CostRange struct {
	From Cost
	To   Cost
}

func (r CostRange) String() string {
	return r.From.String() + "..." + r.To.String()
}

func ParseCostRange(s string) *CostRange {
	segs := strings.Split(s, "...")
	return &CostRange{
		From: ParseCost(segs[0]),
		To:   ParseCost(segs[1]),
	}
}
