package model

type PriorityQueue []Tree

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Freq() < pq[j].Freq()
}

func (pq *PriorityQueue) Push(e interface{}) {
	// prepend needed to mantain every new node in front
	*pq = append(PriorityQueue{e.(Tree)}, *pq...)
}

func (pq *PriorityQueue) Pop() (e interface{}) {
	e = (*pq)[len(*pq)-1]
	*pq = (*pq)[:len(*pq)-1]
	return
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}
