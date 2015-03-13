package graph

import (
	"fmt"
	"sort"
)

const debug = false

const maxDepth = 4

// NodeFetcher is a function that can lazily load Node data.
type NodeFetcher func(*Node)

var defaultNodeFetcher NodeFetcher = func(n *Node) {}
var defaultNodeGroup = NewNodeGroup()

// Node represents a graph node.
type Node struct {
	ID         string
	neighbours []*Node
	loaded     bool
	load       NodeFetcher
	group      *NodeGroup
}

// NewNode constructs a new node with an ID and a lazy loader, and returns a pointer to the newly constructed Node.
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

	return &Node{ID: id, load: loader, group: group}
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
	paths := n.pathsTo(target, 0, Path{}, []Path{})
	sort.Stable(byPathLength(paths))
	return paths
}

func (n *Node) pathsTo(target *Node, depth int, currentPath Path, allPaths []Path) []Path {
	if debug {
		for i := 0; i <= depth; i++ {
			fmt.Printf("\t")
		}
		fmt.Printf("pathsTo(%v, %v, %v, >>%v<<, ##%v##)\n", n, target, depth, currentPath, allPaths)
	}
	// Lazy load Node if required
	if !n.loaded {
		n.load(n)
		n.loaded = true
	}
	// Skip if this node has already been visited in the current run
	if currentPath.Contains(n) {
		return allPaths
	}
	currentPath = append(currentPath, n)

	// Visiting the destination node, which shouldn't normally happen, unless when starting from the destination node itself
	// Add destination node to path and return
	if n.Equal(target) {
		allPaths = append(allPaths, currentPath)
		return allPaths
	}
	if n.IsNeighbour(target) {
		// Found destination node. Add to path and return.
		currentPath = append(currentPath, target)
		allPaths = append(allPaths, currentPath)
		return allPaths
	}
	// Search for paths from neighbours
	for _, neighbour := range n.neighbours {
		if depth < maxDepth {
			allPaths = append(allPaths, neighbour.pathsTo(target, depth+1, currentPath, allPaths)...)
		}
	}
	// HACK
	return deDuplicatePaths(allPaths)
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
