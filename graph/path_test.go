package graph

import "testing"

func TestPathStringRepresentation(t *testing.T) {
	a := &Node{ID: "A"}
	b := &Node{ID: "B"}
	c := &Node{ID: "C"}
	pOne := Path{a, b, c}
	pTwo := Path{a, c, b}
	pThree := Path{b, a}
	pFour := Path{c}
	pFive := Path{}

	if pOne.String() != "A -> B -> C" {
		t.Errorf("Path string representation incorrect. Got: %v, Expected: A -> B -> C", pOne.String())
	}
	if pTwo.String() != "A -> C -> B" {
		t.Errorf("Path string representation incorrect. Got: %v, Expected: A -> C -> B", pTwo.String())
	}
	if pThree.String() != "B -> A" {
		t.Errorf("Path string representation incorrect. Got: %v, Expected: B -> A", pThree.String())
	}
	if pFour.String() != "C" {
		t.Errorf("Path string representation incorrect. Got: %v, Expected: C", pFour.String())
	}
	if pFive.String() != "<EMPTY PATH>" {
		t.Errorf("Path string representation incorrect. Got: %v, Expected: <EMPTY PATH>", pFive.String())
	}
}

func TestPathEquality(t *testing.T) {
	a := &Node{ID: "A"}
	b := &Node{ID: "B"}
	c := &Node{ID: "C"}
	pOne := Path{a, b, c}
	pTwo := Path{a, b, c}
	pThree := Path{a, b}
	pFour := Path{a, c, b}

	if !pOne.Equal(pTwo) {
		t.Errorf("%v should have been Equal to %v", pOne, pTwo)
	}
	if !pTwo.Equal(pOne) {
		t.Errorf("%v should have been Equal to %v", pTwo, pOne)
	}
	if pOne.Equal(pThree) {
		t.Errorf("%v should not have been Equal to %v", pOne, pThree)
	}
	if pOne.Equal(pFour) {
		t.Errorf("%v should not have been Equal to %v", pOne, pFour)
	}
}

func TestPathContainsNode(t *testing.T) {
	a := &Node{ID: "A"}
	b := &Node{ID: "B"}
	c := &Node{ID: "C"}
	d := &Node{ID: "D"}
	path := Path{a, b, c}
	if !path.Contains(a) {
		t.Errorf("%v should contain %v, but it didn't", path, a)
	}
	if !path.Contains(b) {
		t.Errorf("%v should contain %v, but it didn't", path, b)
	}
	if !path.Contains(c) {
		t.Errorf("%v should contain %v, but it didn't", path, c)
	}
	if path.Contains(d) {
		t.Errorf("%v should not contain %v, but it didn't", path, d)
	}
}
