package graph

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConstruction(t *testing.T) {
	assert.NotNil(t, NewNodeGroup().nodes)
}

func TestNodeRegistration(t *testing.T) {
	group := NewNodeGroup()
	err := group.Register(&Node{ID: "one"})
	err = group.Register(&Node{ID: "two"})

	assert.Nil(t, err)
	assert.Equal(t, len(group.nodes), 2)

	err = group.Register(&Node{ID: "one"})
	assert.NotNil(t, err)
	assert.Equal(t, len(group.nodes), 2)
}

func TestGetNode(t *testing.T) {
	group := NewNodeGroup()
	node := &Node{ID: "one"}
	group.Register(node)

	foundNode, present := group.Get("one")
	assert.True(t, present)
	if node != foundNode {
		t.Errorf("Expected original and found nodes to be identical, but they weren't.")
	}
}
