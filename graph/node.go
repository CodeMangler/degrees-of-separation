package graph

import "sort"

// Node represents a graph node.
type Node struct {
	ID           string
	Neighbours   []*Node
	paths        map[string][]Path
	visitedNodes map[string]*Node
}

// String returns a string representation of the Node.
func (node *Node) String() string {
	return node.ID
}

// Equal defines equality of two Nodes.
// Two graph nodes are equal if their IDs are equal, irrespective of the rest of their state.
func (node *Node) Equal(other *Node) bool {
	return node.ID == other.ID
}

// Connect bidirectionally connects two graph Nodes.
func (node *Node) Connect(other *Node) {
	node.Neighbours = appendNodeIfMissing(node.Neighbours, other)
	other.Neighbours = appendNodeIfMissing(other.Neighbours, node)
}

// IsNeighbour returns true if the given node is an immediate neighbour of the current node, false otherwise.
func (node *Node) IsNeighbour(other *Node) bool {
	for _, neighbour := range node.Neighbours {
		if other.Equal(neighbour) {
			return true
		}
	}
	return false
}

// PathsTo computes all possible paths from the current node to the target node.
// It returns an empty slice when no paths are available.
func (node *Node) PathsTo(target *Node) []Path {
	paths := node.pathsTo(target, Path{}, []Path{})
	sort.Stable(byPathLength(paths))
	return paths
}

func (node *Node) pathsTo(target *Node, currentPath Path, allPaths []Path) []Path {
	// Skip if this node has already been visited in the current run
	if currentPath.Contains(node) {
		return allPaths
	}
	currentPath = append(currentPath, node)

	// Visiting the destination node, which shouldn't normally happen, unless when starting from the destination node itself
	// Add destination node to path and return
	if node.Equal(target) {
		allPaths = append(allPaths, currentPath)
		return allPaths
	}
	if node.IsNeighbour(target) {
		// Found destination node. Add to path and return.
		currentPath = append(currentPath, target)
		allPaths = append(allPaths, currentPath)
		return allPaths
	}
	// Search for paths from neighbours
	for _, neighbour := range node.Neighbours {
		allPaths = append(allPaths, neighbour.pathsTo(target, currentPath, allPaths)...)
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
