package graph

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNodeConstruction(t *testing.T) {
	node := NewNode("one")
	assert.Equal(t, node.load, defaultNodeFetcher)
	assert.Equal(t, node.group, defaultNodeGroup)
	//	assert.NotNil(t, node.paths)

	var nodeFetcher NodeFetcher = func(n *Node) {}
	node = NewNode("two", nodeFetcher)
	assert.Equal(t, node.load, nodeFetcher)
	assert.Equal(t, node.group, defaultNodeGroup)
	//	assert.NotNil(t, node.paths)

	nodeGroup := NewNodeGroup()
	node = NewNode("two", nil, nodeGroup)
	assert.Equal(t, node.load, defaultNodeFetcher)
	assert.Equal(t, node.group, nodeGroup)
	//	assert.NotNil(t, node.paths)
}

func TestGraphConstruction(t *testing.T) {
	n := &Node{ID: "a"}
	n.Connect(&Node{ID: "b", neighbours: []*Node{&Node{ID: "c"}}})
	if len(n.neighbours) != 1 {
		t.Errorf("Found %d neighbours on 'a'. Expected 'a' to have exactly one neighbour.", len(n.neighbours))
	}

	if n.neighbours[0].ID != "b" {
		t.Errorf("Found %v to be first neighbour of 'a'. Expected it to be 'b'", n.neighbours[0].ID)
	}
}

func TestNodeStringRepresentation(t *testing.T) {
	nodeOne := &Node{ID: "One"}
	nodeTwo := &Node{ID: "Two", neighbours: []*Node{nodeOne}}
	if nodeOne.String() != "One" {
		t.Errorf("String representation of node incorrect. Got: %v, Expected: One", nodeOne.String())
	}
	if nodeTwo.String() != "Two" {
		t.Errorf("String representation of node incorrect. Got: %v, Expected: Two", nodeTwo.String())
	}
}

func TestNodeEquality(t *testing.T) {
	nodeOne := &Node{ID: "A"}
	nodeTwo := &Node{ID: "A"}
	nodeThree := nodeTwo
	nodeFour := &Node{ID: "A", neighbours: []*Node{&Node{ID: "C"}}}
	nodeFive := &Node{ID: "B"}

	if !nodeOne.Equal(nodeTwo) {
		t.Errorf("%v was not equal to %v. Expected them to be equal.", nodeOne, nodeTwo)
	}
	if !nodeOne.Equal(nodeThree) {
		t.Errorf("%v was not equal to %v. Expected them to be equal.", nodeOne, nodeThree)
	}
	if !nodeOne.Equal(nodeFour) {
		t.Errorf("%v was not equal to %v. Nodes should be equal if their IDs match, regardless of the rest of the state.", nodeOne, nodeFour)
	}
	if nodeOne.Equal(nodeFive) {
		t.Errorf("%v was equal to %v. Expected them not to be equal since their IDs don't match.", nodeOne, nodeFive)
	}
}

func TestNodeNeighbours(t *testing.T) {
	a := &Node{ID: "A"}
	b := &Node{ID: "B"}
	c := &Node{ID: "C"}
	d := &Node{ID: "D"}
	e := &Node{ID: "E"}
	f := &Node{ID: "F"}
	g := &Node{ID: "G"}

	a.Connect(b)
	a.Connect(c)
	a.Connect(g)
	b.Connect(c)
	b.Connect(d)
	d.Connect(e)
	e.Connect(c)
	e.Connect(f)
	f.Connect(g)

	if !a.IsNeighbour(b) {
		t.Errorf("Expected 'A' to have 'B' as it's neighbour")
	}
	if !b.IsNeighbour(a) {
		t.Errorf("Expected 'B' to have 'A' as it's neighbour")
	}
	if !a.IsNeighbour(g) {
		t.Errorf("Expected 'A' to have 'G' as it's neighbour")
	}
	if !a.IsNeighbour(c) {
		t.Errorf("Expected 'A' to have 'C' as it's neighbour")
	}
	if a.IsNeighbour(d) {
		t.Errorf("'A' and 'D' were not expected to be neighbours")
	}
	if a.IsNeighbour(e) {
		t.Errorf("'A' and 'E' were not expected to be neighbours")
	}
	if a.IsNeighbour(f) {
		t.Errorf("'A' and 'F' were not expected to be neighbours")
	}
}

func TestSimplePathComputation(t *testing.T) {
	/*
	      H---I
	   A      |
	          J
	*/
	a := NewNode("A")
	h := NewNode("H")
	i := NewNode("I")
	j := NewNode("J")

	h.Connect(i)
	i.Connect(j)

	hToJ := h.PathsTo(j)
	if len(hToJ) != 1 {
		t.Fatalf("Found %d path(s) from H to J. Expected exactly one path from H to J.", len(hToJ))
	}
	if !hToJ[0].Equal(Path{h, i, j}) {
		t.Errorf("Path from H to J was: %v. Expected Path to be H -> I -> J", hToJ[0])
	}

	jToI := j.PathsTo(i)
	if len(jToI) != 1 {
		t.Fatalf("Found %d path(s) from J to I. Expected exactly one path from J to I.", len(jToI))
	}
	if !jToI[0].Equal(Path{j, i}) {
		t.Errorf("Path from J to I was: %v. Expected Path to be J -> I", jToI[0])
	}

	jToJ := j.PathsTo(j)
	if len(jToJ) != 1 {
		t.Fatalf("Found %d path(s) from J to J. Expected exactly one path from J to J.", len(jToJ))
	}
	if !jToJ[0].Equal(Path{j}) {
		t.Errorf("Path from J to J was: %v. Expected Path to be J", jToJ[0])
	}

	jToA := j.PathsTo(a)
	if len(jToA) != 0 {
		t.Fatalf("Found %d path(s) from J to A. Expected no paths to be available from J to A.", len(jToA))
	}
}

func TestMultiplePathComputation(t *testing.T) {
	/*
	     G-----F
	    / |    |
	   /  |    |    H---I
	   A--C----E        |
	   \  |    /        J
	    \ |   /
	     B---D
	*/
	a := NewNode("A")
	b := NewNode("B")
	c := NewNode("C")
	d := NewNode("D")
	e := NewNode("E")
	f := NewNode("F")
	g := NewNode("G")

	a.Connect(b)
	a.Connect(c)
	a.Connect(g)
	b.Connect(c)
	b.Connect(d)
	c.Connect(g)
	d.Connect(e)
	e.Connect(c)
	e.Connect(f)
	f.Connect(g)

	aToF := a.PathsTo(f)
	assert.Equal(t, 10, len(aToF))
	if len(aToF) == 10 {
		assert.Equal(t, Path{a, g, f}.String(), aToF[0].String())
		assert.Equal(t, Path{a, c, g, f}.String(), aToF[1].String())
		assert.Equal(t, Path{a, c, e, f}.String(), aToF[2].String())
		assert.Equal(t, Path{a, g, c, e, f}.String(), aToF[3].String())
		assert.Equal(t, Path{a, b, d, e, f}.String(), aToF[4].String())
		assert.Equal(t, Path{a, b, c, g, f}.String(), aToF[5].String())
		assert.Equal(t, Path{a, b, c, e, f}.String(), aToF[6].String())
		assert.Equal(t, Path{a, c, b, d, e, f}.String(), aToF[7].String())
		assert.Equal(t, Path{a, g, c, b, d, e, f}.String(), aToF[8].String())
		assert.Equal(t, Path{a, b, d, e, c, g, f}.String(), aToF[9].String())
	}
}

func TestNodeLazyLoading(t *testing.T) {
	a := NewNode("A")
	b := NewNode("B")
	c := NewNode("C")

	loaderWasCalled := false
	a.loaded = false
	a.load = func(n *Node) {
		loaderWasCalled = true
		n.Connect(b)
		b.Connect(c)
	}

	aToC := a.PathsTo(c)

	if !loaderWasCalled {
		t.Fatalf("Loader should've been called on A")
	}
	if len(aToC) != 1 {
		t.Fatalf("Paths were not calculated from data loaded by the loader. Got %v paths. Expected 1", len(aToC))
	}
	if !aToC[0].Equal(Path{a, b, c}) {
		t.Errorf("Path from A to C was incorrectly calculated. Got: %v. Expected A -> B -> C", aToC)
	}
}

/*
func TestPathCaching(t *testing.T) {
//		A--B----C
//		   \   /
//		    \ /
//		     D

	a := NewNode("A")
	b := NewNode("B")
	c := NewNode("C")
	d := NewNode("D")

	a.Connect(b)
	b.Connect(c)
	b.Connect(d)
	c.Connect(d)

	a.PathsTo(d)
	assert.Equal(t, 1, len(a.paths))
	var cachedPaths []Path = a.paths["D"]
	sort.Stable(byPathLength(cachedPaths))

	assert.Equal(t, 2, len(cachedPaths))
	if len(cachedPaths) == 2 {
		assert.Equal(t, Path{a, b, d}, cachedPaths[0])
		assert.Equal(t, Path{a, b, c, d}, cachedPaths[1])
	}
}
*/
