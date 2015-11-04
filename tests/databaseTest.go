/*
 * CS350 Team 5 - shortests route for a given network topology
 * Iteration 3 - graph data storage
 * David Josephs and Killian Coddington
 */

package main

import (
	"github.com/Grissess/sdnd16/database"
	"fmt"
)

func main() {

	numberOfNodes := 3
    address := "128.153.144.171:6379"

	rdb := database.NewRoutingDatabase("test", numberOfNodes)
    err := rdb.Connect("tcp", address)
    if err != nil {
        panic(err)
    }

    // this ip address is that of a virtual machine hosted on
    // phoenix at cosi

    rdb.SetTrivialPaths()
	for i := 0; i < numberOfNodes; i++ {
		for j := 0; j < numberOfNodes; j++ {
			if i != j {
                rdb.SetPath(i, j, fmt.Sprintf("%d %d | %d", i, j, i + j))
            }
        }
	}

	rdb.StorePathsInDB()

	var s, d int
    var path string

	for {
		fmt.Print("enter the pair of nodes you want a path between (s, d) > ")
		fmt.Scanf("(%d, %d)\n", &s, &d)
		if s >= numberOfNodes || s < 0 || d >= numberOfNodes || d < 0 {
			fmt.Printf("at least one specified node is not in the range 0 - %d\n", (numberOfNodes - 1))
			continue
		}
		path, err = rdb.GetPathFromDB(s, d)
        if err != nil {
            panic(err)
        }
		fmt.Printf("the shortest path from %d to %d is: %s\n", s, d, path)
	}

}
