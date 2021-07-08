package model

/*
	custom structure to build a sortered slice
	of (symbol, frequency).
*/

type Pair struct {
	Symbol    rune
	Frequency int
}

type Pairs []Pair

func (p Pairs) Len() int { return len(p) }

func (p Pairs) Less(i, j int) bool { return p[i].Frequency < p[j].Frequency }

func (p Pairs) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
