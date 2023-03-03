package main

import (
	"container/heap"
	"fmt"
	"strconv"

	"github.com/go-echarts/go-echarts/v2/opts"
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

func createHeap() (nodes []opts.GraphNode, links []opts.GraphLink) {
	h := intHeap{}
	// for _, v := range []int{9, 86, 1, 2, 10, 50, 37} {
	for _, v := range []int{9, 86, 1, 2, 10, 50, 37, 900, 43, 513, 367} {
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

type direction int

const (
	_ direction = iota
	directionLeft
	directionRight
)

func (d direction) String() string {
	return [3]string{"", "left", "right"}[d]
}

type treeNode[T nodeType] struct {
	depth           int
	x, y            int
	node            T
	parent          *treeNode[T]
	parentDirection direction
}

// type tree struct {
// 	nodes []*treeNode[*minIntHeapNode]
// 	// edges [][]int
// }

type stack[T nodeType] []stackNode[T]
type stackNode[T nodeType] struct {
	node    T
	index   int // index of the node: node[T][index]
	depth   int
	xOffset int
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
	*s = (*s)[:topIndex]
	return top
}

func buildGraphFromHeap(h *intHeap) (nodes []opts.GraphNode, links []opts.GraphLink) {
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

	nodes = make([]opts.GraphNode, h.Len())

	// depth first search
	for !nodeStack.IsEmpty() {
		node := nodeStack.Pop()

		nodes[node.index] = opts.GraphNode{
			Value: float32(node.node.value),
			Name:  strconv.Itoa(node.node.value),
			ItemStyle: &opts.ItemStyle{
				Color: fmt.Sprintf("%x", node.node.value),
			},
			SymbolSize: 20,
			X:          float32(node.xOffset*10 + widthOffset),
			Y:          float32(node.depth*10 + heightOffset),
		}

		rightChildIndex := node.index*2 + 2
		if rightChildIndex < h.Len() {
			nodeStack.Push(stackNode[*minIntHeapNode]{
				node:    &(*h)[rightChildIndex],
				index:   rightChildIndex,
				depth:   node.depth + 1,
				xOffset: node.xOffset + 1,
			})
			links = append(links, opts.GraphLink{
				Source: node.index,
				Target: rightChildIndex,
			})
		}

		leftChildIndex := node.index*2 + 1
		if leftChildIndex < h.Len() {
			nodeStack.Push(stackNode[*minIntHeapNode]{
				node:    &(*h)[leftChildIndex],
				index:   leftChildIndex,
				depth:   node.depth + 1,
				xOffset: node.xOffset - 1,
			})
			links = append(links, opts.GraphLink{
				Source: node.index,
				Target: leftChildIndex,
			})
		}
	}

	return
}

// func buildTreeFromHeap(h *intHeap) *tree {
// 	g := &tree{
// 		nodes: make([]*treeNode[*minIntHeapNode], h.Len()),
// 		edges: make([][]int, h.Len()-1), // max amount of edges
// 		// edges: [][]int{}, // max amount of edges
// 	}
// 	// initialize the nodes
// 	for k := range g.nodes {
// 		g.nodes[k] = &treeNode[*minIntHeapNode]{}
// 	}

// 	depth := 0
// 	const heightCoefficient int = 20
// 	const widthCoefficient int = 20

// 	const branchingFactor = 2                            // for a binary tree this is 2; likely will be a variable later
// 	nextLevel := branchingFactor*depth + branchingFactor // index at which the next depth has been reached

// 	for i := 0; i < h.Len(); i++ {
// 		if i > nextLevel {
// 			depth++
// 			nextLevel = branchingFactor*depth + branchingFactor
// 			fmt.Printf("setting depth to %d", depth)
// 		} else if i == 1 {
// 			depth = 1
// 		}

// 		// Add node to the tree
// 		g.nodes[i].depth = depth
// 		g.nodes[i].y = depth * heightCoefficient
// 		g.nodes[i].node = &(*h)[i]

// 		// Add edges for the node to the tree
// 		leftChildIndex := 2*i + 1
// 		if leftChildIndex < h.Len() {
// 			if g.edges[i] == nil {
// 				g.edges[i] = []int{}
// 			}
// 			g.nodes[leftChildIndex].parent = g.nodes[i]
// 			g.nodes[leftChildIndex].parentDirection = directionRight
// 			g.edges[i] = append(g.edges[i], leftChildIndex)
// 			// g.nodes[leftChildIndex].x = -1*widthCoefficient - g.nodes[i].x //*g.nodes[i].depth
// 			if g.nodes[i].parent != nil {
// 				g.nodes[leftChildIndex].x = -1*widthCoefficient - g.nodes[i].parent.x //*g.nodes[i].depth
// 			} else {
// 				g.nodes[leftChildIndex].x = -1 * widthCoefficient
// 			}
// 		}

// 		rightChildIndex := 2*i + 2
// 		if rightChildIndex < h.Len() {
// 			if g.edges[i] == nil {
// 				g.edges[i] = []int{}
// 			}
// 			g.nodes[leftChildIndex].parentDirection = directionLeft
// 			g.nodes[rightChildIndex].parent = g.nodes[i]
// 			g.edges[i] = append(g.edges[i], rightChildIndex)
// 			if g.nodes[i].parent != nil {
// 				g.nodes[rightChildIndex].x = widthCoefficient + g.nodes[i].parent.x //*g.nodes[i].depth
// 			} else {
// 				g.nodes[rightChildIndex].x = widthCoefficient
// 			}
// 		}
// 	}

// 	// for i := range g.nodes {
// 	// 	g.nodes[i].
// 	// }

// 	return g
// }

// func chartTree() (nodes []opts.GraphNode, links []opts.GraphLink) {
// 	nodes = make([]opts.GraphNode, len(g.nodes))

// 	const heightOffset int = 200
// 	const widthOffset int = 200

// 	// for i, node := range g.nodes {
// 	for i, treeNode := range g.nodes {
// 		nodes[i] = opts.GraphNode{
// 			Value: float32(treeNode.node.value),
// 			Name:  strconv.Itoa(treeNode.node.value),
// 			ItemStyle: &opts.ItemStyle{
// 				Color: fmt.Sprintf("%x", treeNode.node.value),
// 			},
// 			SymbolSize: 20,
// 			X:          float32(treeNode.x + widthOffset),
// 			Y:          float32(treeNode.y + heightOffset),
// 		}
// 	}

// 	for source, targets := range g.edges {
// 		for _, target := range targets {
// 			links = append(links, opts.GraphLink{
// 				Source: source,
// 				Target: target,
// 				Label: &opts.EdgeLabel{
// 					Show:     true,
// 					Position: "middle",
// 					Formatter: fmt.Sprintf(
// 						"from: %d -> to: %d\nfrom id: %d-> to id: %d\nparent: %s x: %d",
// 						// "from: %+v -> to: %+v",
// 						g.nodes[source].x, g.nodes[target].x,
// 						source, target,
// 						g.nodes[target].parent.parentDirection.String(), g.nodes[target].parent.x),
// 					// Formatter: fmt.Sprintf(
// 					// 	"from: %.f -> to: %.f",
// 					// 	nodes[source].X, nodes[target].Y),
// 					// Formatter: fmt.Sprintf(
// 					// 	"from: %d -> to: %d\ntarget.x: %.f, target.y:%.f",
// 					// 	source, target, nodes[target].X, nodes[target].Y),
// 				},
// 			})
// 		}
// 	}

// 	return nodes, links
// }
