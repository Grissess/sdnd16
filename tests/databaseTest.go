/*
 * CS350 Team 5 - shortests route for a given network topology
 * Iteration 2 - graph data storage
 * David Josephs and Killian Coddington
 */

package main

import (
	"../database"
	"fmt"
)

func main() {

	numberOfNodes := 3

	rdb := database.NewRoutingDatabase("test", numberOfNodes)
	rdb.Connect("tcp", ":6379")

	for i := 0; i < numberOfNodes; i++ {
		for j := 0; j < numberOfNodes; j++ {
			rdb.SetPath(i, j, fmt.Sprintf("%d %d | %d", (i + 1), (j + 1), (i + j) + 2))
		}
	}

	rdb.StorePathsData()

	var s, d int

	for {
		fmt.Print("enter the pair of nodes you want a path between (s, d) > ")
		fmt.Scanf("(%d, %d)\n", &s, &d)
		if s > numberOfNodes || s <= 0 || d > numberOfNodes || d <= 0 {
			fmt.Printf("at least one specified node is not in the range 1 - %d\n", numberOfNodes)
			continue
		}
		path := rdb.DBGetPath((s - 1), (d - 1))
		fmt.Printf("the shortest path from %d to %d is: %s\n", s, d, path)
	}

}
