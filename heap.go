package main

import (
	"container/heap"
	"fmt"
	"strconv"

	"github.com/go-echarts/go-echarts/v2/opts"
	spatialmap "github.com/nrhvyc/go-visualization/spatial_map"
)

type minIntHeapNode struct {
	value int
	// name  string
}

type intHeap []minIntHeapNode

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
	*h = append(*h, x.(minIntHeapNode))
}
func (h *intHeap) Pop() any {
	old := *h
	n := len(old)
	popped := old[n-1]
	*h = old[:n-1]
	return popped
}
func (h *intHeap) Parent(childIndex int) *minIntHeapNode {
	return &(*h)[(childIndex-1)/2]
}
func (h *intHeap) LeftChild(parentIndex int) *minIntHeapNode {
	return &(*h)[parentIndex*2+1]
}
func (h *intHeap) RightChild(parentIndex int) *minIntHeapNode {
	return &(*h)[parentIndex*2+2]
}

func createHeap() (nodes []opts.GraphNode, links []opts.GraphLink) {
	h := intHeap{}
	// for _, v := range []int{9, 86, 1, 2, 10, 50, 37} {
	// for _, v := range []int{9, 86, 1, 2, 10, 50, 37, 900, 43, 513, 367} {
	for _, v := range []int{54, 23, 87, 42, 10, 68, 31, 91, 12, 79, 61, 36, 17, 99, 72, 45, 6, 29, 83, 57, 14, 67, 38, 94, 20, 76, 49, 8, 63, 27, 89, 52, 19, 73, 39, 97, 25, 80, 65, 16, 84, 58, 11, 70, 33, 96, 22, 78, 55, 2, 69} {
		// for _, v := range []int{9, 86, 1, 2} {
		h = append(h, minIntHeapNode{value: v})
	}
	heap.Init(&h)
	// a := heap.Pop(&h).(node)
	// fmt.Printf("a: %d\n", a)

	return buildGraphFromHeap(&h)
}

type nodeType interface {
	*minIntHeapNode
}

// type direction int

// const (
// 	_ direction = iota
// 	directionLeft
// 	directionRight
// )

// func (d direction) String() string {
// 	return [3]string{"", "left", "right"}[d]
// }

// type treeNode[T nodeType] struct {
// 	depth           int
// 	x, y            int
// 	node            T
// 	parent          *treeNode[T]
// 	parentDirection direction
// }

// type tree struct {
// 	nodes []*treeNode[*minIntHeapNode]
// 	// edges [][]int
// }

type stack[T nodeType] []stackNode[T]
type stackNode[T nodeType] struct {
	node        T
	index       int // index of the node: node[T][index]
	depth       int
	xOffset     int
	childOffset int // which node in the sequence of the parent is this node (root: 0, first: 1, second: 2)
}

func (s *stack[T]) IsEmpty() bool {
	return len(*s) == 0
}
func (s *stack[T]) Push(n stackNode[T]) {
	*s = append(*s, n)
}
func (s *stack[T]) Pop() (node stackNode[T]) {
	topIndex := len(*s) - 1
	top := (*s)[topIndex]
	*s = (*s)[:topIndex] // remove top
	return top
}

func buildGraphFromHeap(h *intHeap) (graphNodes []opts.GraphNode, links []opts.GraphLink) {
	if h == nil {
		return
	}

	const heightOffset int = 200
	const widthOffset int = 200

	nodeStack := stack[*minIntHeapNode]{
		// Insert heap root node
		stackNode[*minIntHeapNode]{
			node:  &(*h)[0],
			index: 0,
			depth: 0,
		},
	}

	graphNodes = make([]opts.GraphNode, h.Len())
	spatialMap := spatialmap.NewSpatialMap[opts.GraphNode](1)

	// depth first search
	for !nodeStack.IsEmpty() {
		node := nodeStack.Pop()

		graphNode := opts.GraphNode{
			Value: float32(node.node.value),
			Name:  strconv.Itoa(node.node.value),
			ItemStyle: &opts.ItemStyle{
				Color: fmt.Sprintf("%x", node.node.value),
			},
			SymbolSize: 20,
			X:          float32(node.xOffset*10 + widthOffset),
			Y:          float32(node.depth*10 + heightOffset),
		}
		// Keep offsetting nodes if they overlap
		for overlapNodes := spatialMap.Get(int(graphNode.X), int(graphNode.Y)); len(overlapNodes) > 0; {
			graphNode.X += 10
			overlapNodes = spatialMap.Get(int(graphNode.X), int(graphNode.Y))
		}
		if node.childOffset != 0 {
			parentIndex := (node.index - node.childOffset) / 2
			avgChildX := (graphNodes[(parentIndex-1)/2].X + graphNodes[(parentIndex-2)/2].X) / 2
			graphNodes[parentIndex].X = avgChildX
		}

		graphNodes[node.index] = graphNode
		spatialMap.Add(int(graphNode.X), int(graphNode.Y), &graphNode)

		rightChildIndex := node.index*2 + 2
		if rightChildIndex < h.Len() {
			nodeStack.Push(stackNode[*minIntHeapNode]{
				node:        &(*h)[rightChildIndex],
				index:       rightChildIndex,
				depth:       node.depth + 1,
				xOffset:     node.xOffset + 1,
				childOffset: 2,
			})
			links = append(links, opts.GraphLink{
				Source: node.index,
				Target: rightChildIndex,
			})
		}

		leftChildIndex := node.index*2 + 1
		if leftChildIndex < h.Len() {
			nodeStack.Push(stackNode[*minIntHeapNode]{
				node:        &(*h)[leftChildIndex],
				index:       leftChildIndex,
				depth:       node.depth + 1,
				xOffset:     node.xOffset - 1,
				childOffset: 1,
			})
			links = append(links, opts.GraphLink{
				Source: node.index,
				Target: leftChildIndex,
			})
		}
	}

	return
}
