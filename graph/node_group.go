package graph

import "errors"

var defaultNodeGroup = NewNodeGroup()

// NodeGroup represents a group of related Nodes
type NodeGroup struct {
	nodes             map[string]*Node
	maxRecursionDepth int
}

// NewNodeGroup creates a new NodeGroup
// Parameter 1: maxRecursionDepth. Defaults to 6.
func NewNodeGroup(args ...interface{}) *NodeGroup {
	maxRecursionDepth := 6
	if len(args) > 0 {
		maxRecursionDepth = args[0].(int)
	}
	return &NodeGroup{nodes: make(map[string]*Node), maxRecursionDepth: maxRecursionDepth}
}

// Register registers a Node with the current NodeGroup by it's ID
func (g *NodeGroup) Register(node *Node) error {
	if _, exists := g.Get(node.ID); exists {
		return errors.New("Another node has already been registered with the same ID")
	}
	g.nodes[node.ID] = node
	node.group = g
	return nil
}

// Get finds and returns an existing Node in the current NodeGroup matching the given ID
// Returns Node, true if found. Returns nil, false if not found.
func (g *NodeGroup) Get(id string) (*Node, bool) {
	node, present := g.nodes[id]
	return node, present
}
