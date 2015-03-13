package graph

import "errors"

// NodeGroup represents a group of related Nodes
type NodeGroup struct {
	nodes map[string]*Node
}

// NewNodeGroup creates a new NodeGroup
func NewNodeGroup() *NodeGroup {
	return &NodeGroup{nodes: make(map[string]*Node)}
}

// Register registers a Node with the current NodeGroup by it's ID
func (g *NodeGroup) Register(node *Node) error {
	if _, exists := g.nodes[node.ID]; exists {
		return errors.New("Another node has already been registered with the same ID")
	}
	g.nodes[node.ID] = node
	return nil
}
