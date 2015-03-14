package graph

import (
	"fmt"
	"sort"
	"time"
)

const debug = false

// NodeFetcher is a function that can lazily load Node data.
type NodeFetcher func(*Node) error

var defaultNodeFetcher NodeFetcher = func(n *Node) error {
	n.Data = true
	return nil
}

// Node represents a graph node.
type Node struct {
	ID         string
	Data       interface{}
	neighbours []*Node
	load       NodeFetcher
	group      *NodeGroup
	//	paths      map[string][]Path
}

// NewNode constructs a new node with an ID and a lazy loader, and returns a pointer to the newly constructed Node.
// Parameter 1: id - ID for the new Node.
// Parameter 2: loader - NodeFetcher to lazy load the Node. Defaults to an empty NodeFetcher if not specified.
// Parameter 3: group - NodeGroup that this Node should belong to. Defaults to a default NodeGroup if not specified.
func NewNode(id string, otherArgs ...interface{}) *Node {
	var loader NodeFetcher
	var group *NodeGroup

	if len(otherArgs) > 0 && otherArgs[0] != nil {
		loader = otherArgs[0].(NodeFetcher)
	}
	if len(otherArgs) > 1 {
		group = otherArgs[1].(*NodeGroup)
	}

	if loader == nil {
		loader = defaultNodeFetcher
	}
	if group == nil {
		group = defaultNodeGroup
	}

	node, present := group.Get(id)
	if present {
		return node
	}
	node = &Node{ID: id, load: loader /*paths: make(map[string][]Path)*/}
	group.Register(node)
	return node
}

// String returns a string representation of the Node.
func (n *Node) String() string {
	return n.ID
}

// Equal defines equality of two Nodes.
// Two graph nodes are equal if their IDs are equal, irrespective of the rest of their state.
func (n *Node) Equal(other *Node) bool {
	return n.ID == other.ID
}

// Connect bidirectionally connects two graph Nodes.
func (n *Node) Connect(other *Node) {
	n.neighbours = appendNodeIfMissing(n.neighbours, other)
	other.neighbours = appendNodeIfMissing(other.neighbours, n)
}

// IsNeighbour returns true if the given node is an immediate neighbour of the current node, false otherwise.
func (n *Node) IsNeighbour(other *Node) bool {
	for _, neighbour := range n.neighbours {
		if other.Equal(neighbour) {
			return true
		}
	}
	return false
}

// PathsTo computes all possible paths from the current node to the target node.
// It returns an empty slice when no paths are available.
func (n *Node) PathsTo(target *Node) []Path {
	if debug {
		fmt.Println()
	}
	chanResults := make(chan []Path)
	go n.pathsTo(target, 0, Path{}, chanResults)
	paths := <-chanResults
	sort.Stable(byPathLength(paths))
	return paths
}

func (n *Node) pathsTo(target *Node, depth int, currentPath Path, chanResults chan []Path) {
	if debug {
		defer func() {
			for i := 0; i <= depth; i++ {
				fmt.Printf("\t")
			}
			fmt.Printf("Returning from pathsTo(%v, %v, %v, >>%v<<)\n", n, target, depth, currentPath)
		}()
	}
	if debug {
		for i := 0; i <= depth; i++ {
			fmt.Printf("\t")
		}
		fmt.Printf("pathsTo(%v, %v, %v, >>%v<<)\n", n, target, depth, currentPath)
	}
	// Lazy load Node
	for n.Data == nil {
		err := n.load(n)
		// Retry loading node after a pause if there was an error while loading
		if err != nil {
			time.Sleep(1 * time.Second)
		}
	}

	// Skip if this node has already been visited in the current run
	if currentPath.Contains(n) {
		chanResults <- []Path{}
		return
	}
	currentPath = append(currentPath, n)

	if n.Equal(target) {
		chanResults <- []Path{currentPath}
		//		n.paths[target.ID] = append(n.paths[target.ID], currentPath)
		return
	}

	// Search for paths from neighbours
	chanNeighbourResults := make(chan []Path)
	if depth < n.group.maxRecursionDepth {
		for _, neighbour := range n.neighbours {
			go neighbour.pathsTo(target, depth+1, currentPath, chanNeighbourResults)
		}
	}

	results := []Path{}
	if debug {
		for i := 0; i <= depth; i++ {
			fmt.Printf("\t")
		}
		fmt.Printf("[%v -> %v] Waiting to read from %v routines\n", n, target, len(n.neighbours))
	}
	for i := 0; i < len(n.neighbours); i++ {
		if debug {
			for i := 0; i <= depth; i++ {
				fmt.Printf("\t")
			}
			fmt.Printf("[%v -> %v] Reading results of %v\n", n, target, i)
		}
		neighbourPaths := <-chanNeighbourResults
		results = append(results, neighbourPaths...)
	}
	//HACK
	results = deDuplicatePaths(results)
	chanResults <- results
	//	n.paths[target.ID] = append(n.paths[target.ID], allPaths...)
	return
}

func appendNodeIfMissing(nodes []*Node, nodeToAppend *Node) []*Node {
	for _, node := range nodes {
		if node.Equal(nodeToAppend) {
			return nodes
		}
	}
	nodes = append(nodes, nodeToAppend)
	return nodes
}

// HACK
func deDuplicatePaths(paths []Path) []Path {
	ddMap := make(map[string]Path)
	for _, path := range paths {
		ddMap[path.String()] = path
	}

	result := []Path{}
	for _, path := range ddMap {
		result = append(result, path)
	}
	return result
}
