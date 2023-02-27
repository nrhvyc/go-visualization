package main

import (
	"container/heap"
	"fmt"
	"strconv"

	"github.com/go-echarts/go-echarts/v2/opts"
)

type node struct {
	value int
	// name  string
}

type intHeap []node

func (h intHeap) Len() int {
	return len(h)
}
func (h intHeap) Less(i, j int) bool {
	return h[i].value < h[j].value
}
func (h *intHeap) Swap(i, j int) {
	x := *h
	temp := x[i]
	x[i] = x[j]
	x[j] = temp
}
func (h *intHeap) Push(x any) {
	*h = append(*h, x.(node))
}

func (h *intHeap) Pop() any {
	old := *h
	n := len(old)
	popped := old[n-1]
	*h = old[:n-1]
	return popped
}

func createHeap() *graph {
	h := intHeap{}
	for _, v := range []int{9, 86, 1, 2, 10, 50, 37} {
		h = append(h, node{value: v})
	}
	heap.Init(&h)
	// a := heap.Pop(&h).(node)
	// fmt.Printf("a: %d\n", a)

	return buildGraphFromHeap(&h)
}

type graph struct {
	nodes []*node
	edges [][]int
}

func buildGraphFromHeap(h *intHeap) *graph {
	g := &graph{
		nodes: make([]*node, h.Len()),
		edges: make([][]int, h.Len()),
	}

	// Add nodes to the graph
	for i := 0; i < h.Len(); i++ {
		n := (*h)[i]
		g.nodes[i] = &n
		g.edges[i] = []int{}
	}

	// Add edges to the graph
	for i := 0; i < h.Len(); i++ {
		leftChildIndex := 2*i + 1
		if leftChildIndex < h.Len() {
			g.edges[i] = append(g.edges[i], leftChildIndex)
		}

		rightChildIndex := 2*i + 2
		if rightChildIndex < h.Len() {
			g.edges[i] = append(g.edges[i], rightChildIndex)
		}
	}

	return g
}

func (g graph) chartGraph() (nodes []opts.GraphNode, links []opts.GraphLink) {
	nodes = make([]opts.GraphNode, len(g.nodes))
	for i, node := range g.nodes {
		nodes[i] = opts.GraphNode{
			Value: float32(node.value),
			Name:  strconv.Itoa(node.value),
			ItemStyle: &opts.ItemStyle{
				Color: fmt.Sprintf("%x", node.value),
			},
		}
	}

	for source, targets := range g.edges {
		for _, target := range targets {
			links = append(links, opts.GraphLink{
				Source: source,
				Target: target,
			})
		}
	}

	return nodes, links
}
