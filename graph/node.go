package graph

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

const (
	debug           = false
	maxLoadAttempts = 3
)

// NodeFetcher is a function that can lazily load Node data.
type NodeFetcher func(*Node) error

var defaultNodeFetcher NodeFetcher = func(n *Node) error {
	n.SetData(true)
	return nil
}

// Node represents a graph node.
type Node struct {
	ID         string
	data       interface{}
	neighbours []*Node
	load       NodeFetcher
	group      *NodeGroup
	lock       sync.Mutex
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

// SetData sets Node data in a thread-safe manner.
func (n *Node) SetData(data interface{}) {
	n.lock.Lock()
	n.data = data
	n.lock.Unlock()
}

// HasData checks for presence of Node data in a thread-safe manner.
func (n *Node) HasData() bool {
	result := false
	n.lock.Lock()
	result = (n.data != nil)
	n.lock.Unlock()
	return result
}

// PathsTo computes all possible paths from the current node to the target node.
// It returns an empty slice when no paths are available.
func (n *Node) PathsTo(target *Node, args ...interface{}) []Path {
	stopAtFirstPath := false
	if len(args) > 0 {
		stopAtFirstPath = args[0].(bool)
	}

	chanResults := make(chan []Path)
	go n.pathsTo(target, 0, Path{n, target}.String(), stopAtFirstPath, Path{}, chanResults)
	paths := <-chanResults
	sort.Stable(byPathLength(paths))
	return paths
}

func (n *Node) pathsTo(target *Node, depth int, pathID string, stopAtFirstPath bool, currentPath Path, chanResults chan []Path) {
	if debug {
		tabs(depth)
		fmt.Printf("pathsTo(%v, %v, %v, >>%v<<)\n", n, target, depth, currentPath)
	}

	if stopAtFirstPath && n.group.pathsFound[pathID] {
		chanResults <- []Path{}
		return
	}
	// Lazy load Node
	loadAttempt := 0
	for !n.HasData() {
		if debug {
			tabs(depth)
			fmt.Printf("Loading %v. Attempt %v\n", n.ID, loadAttempt)
		}
		err := n.load(n)
		// Retry loading node after a pause if there was an error while loading
		if err != nil {
			if loadAttempt > maxLoadAttempts {
				if debug {
					tabs(depth)
					fmt.Printf(">>>>>>>>>>>>>>>>> Failed to load %v. Bailing out.\n", n.ID)
				}
				chanResults <- []Path{}
				return
			}
			loadAttempt++
			time.Sleep(1 * time.Second)
		}
	}

	// Skip if this node has already been visited in the current run
	n.group.lock.Lock()
	loop := currentPath.Contains(n)
	n.group.lock.Unlock()
	if loop {
		chanResults <- []Path{}
		return
	}
	n.group.lock.Lock()
	currentPath = append(currentPath, n)
	n.group.lock.Unlock()

	if n.Equal(target) {
		if debug {
			tabs(depth)
			fmt.Printf("$$$$$$$$$$$$$ [%v] %v -> %v Found: %v\n", pathID, n, target, currentPath)
		}
		n.group.lock.Lock()
		n.group.pathsFound[pathID] = true
		n.group.lock.Unlock()
		chanResults <- []Path{currentPath}
		//		n.paths[target.ID] = append(n.paths[target.ID], currentPath)
		return
	}

	// Search for paths from neighbours
	chanNeighbourResults := make(chan []Path)
	if depth < n.group.maxRecursionDepth {
		for _, neighbour := range n.neighbours {
			go neighbour.pathsTo(target, depth+1, pathID, stopAtFirstPath, currentPath, chanNeighbourResults)
		}
	}

	results := []Path{}
	for i := 0; i < len(n.neighbours); i++ {
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

func tabs(count int) {
	for i := 0; i <= count; i++ {
		fmt.Printf("\t")
	}
}
