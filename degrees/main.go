package main

import (
	"../graph"
	"../moviebuff"
	"fmt"
	"os"
)

func main() {
	sourceID := os.Args[1]
	targetID := os.Args[2]

	nodeGroup := graph.NewNodeGroup(4)
	sourceNode := graph.NewNode(sourceID, graph.NodeFetcher(moviebuff.Fetch), nodeGroup)
	targetNode := graph.NewNode(targetID, graph.NodeFetcher(moviebuff.Fetch), nodeGroup)

	paths := sourceNode.PathsTo(targetNode, true)
	if len(paths) == 0 {
		fmt.Printf("\nCould not find a connection between %v and %v\n", sourceNode, targetNode)
	} else {
		fmt.Printf("\n!!!!!!!!!!!!!!! Degrees of Separation: %v\n", len(paths[0]))
		for i, node := range paths[0] {
			fmt.Printf("%d. %v\n", i, node)
		}
	}
}
