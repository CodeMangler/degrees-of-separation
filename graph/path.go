package graph

// Path represents a collection of Nodes in a certain order.
type Path []*Node

// String returns a string representation of the path
func (p Path) String() string {
	pathLength := len(p)
	if pathLength == 0 {
		return "<EMPTY PATH>"
	}
	result := ""
	for i, node := range p {
		result += node.String()
		if i < pathLength-1 {
			result += " -> "
		}
	}
	return result
}

// Equal defines equality of two Paths.
// Two paths are equal if they have the same number of nodes, specified in the same order.
func (p Path) Equal(other Path) bool {
	if len(p) != len(other) {
		return false
	}
	for i, otherNode := range other {
		if !p[i].Equal(otherNode) {
			return false
		}
	}
	return true
}

// Contains returns true if the current path contains the given Node, false otherwise.
func (p Path) Contains(node *Node) bool {
	for _, pathNode := range p {
		if pathNode.Equal(node) {
			return true
		}
	}
	return false
}

type byPathLength []Path

func (a byPathLength) Len() int      { return len(a) }
func (a byPathLength) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byPathLength) Less(i, j int) bool {
	if len(a[i]) == len(a[j]) {
		return a[j].String() < a[i].String()
	}
	return len(a[i]) < len(a[j])
}
