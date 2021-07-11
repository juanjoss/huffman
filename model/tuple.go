package model

/*
	Used to build the list with the frequencies of
	each symbol in the data to be compressed.
*/

type Pair struct {
	Symbol    rune
	Frequency int
}

type Pairs []Pair

func (p Pairs) Len() int { return len(p) }

func (p Pairs) Less(i, j int) bool { return p[i].Frequency < p[j].Frequency }

func (p Pairs) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
