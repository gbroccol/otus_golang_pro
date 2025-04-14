package hw03frequencyanalysis

import (
	"container/heap"
	"sort"
	"strings"
)

// Определяем структуру для max heap
type MaxHeap []Pair

func (h MaxHeap) Len() int { return len(h) }
func (h MaxHeap) Less(i, j int) bool {
	return h[i].Value > h[j].Value || (h[i].Value == h[j].Value && h[i].Key < h[j].Key)
}
func (h MaxHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *MaxHeap) Push(x interface{}) {
	*h = append(*h, x.(Pair))
}

func (h *MaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type Pair struct {
	Key   string
	Value int
}

func Top10(text string) []string {

	words := strings.Fields(text)
	frequencyMap := map[string]int{}

	for _, word := range words {
		frequencyMap[word] += 1
	}

	//return getTopUsingSlice(frequencyMap)
	return getTopUsingHeap(frequencyMap)
}

func getTopUsingSlice(frequencyMap map[string]int) []string {

	var pairs []Pair

	for k, v := range frequencyMap {
		pairs = append(pairs, Pair{Key: k, Value: v})
	}

	// Sort the pairs based on value, and by key if values are equal
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].Value == pairs[j].Value {
			return pairs[i].Key < pairs[j].Key
		}
		return pairs[i].Value > pairs[j].Value
	})

	// Extract the top 10 keys
	var top10Keys []string
	for i := 0; i < 10 && i < len(pairs); i++ {
		top10Keys = append(top10Keys, pairs[i].Key)
	}

	return top10Keys
}

func getTopUsingHeap(frequencyMap map[string]int) []string {

	h := &MaxHeap{}
	heap.Init(h)

	for k, v := range frequencyMap {
		heap.Push(h, Pair{Key: k, Value: v})
	}

	var top10Keys []string
	for i := 0; i < 10 && h.Len() > 0; i++ {
		top10Keys = append(top10Keys, heap.Pop(h).(Pair).Key)
	}

	return top10Keys
}
