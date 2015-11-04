/*
 * CS350 Team 5 - shortests route for a given network topology
 * Iteration 3 - graph data storage
 * David Josephs and Killian Coddington
 */

package main

import (
	"../database"
	"fmt"
)

func main() {

    address := "128.153.144.171:6379"

    labelMap := make(map[string]int)
    labelMap["alpha"] = 0
    labelMap["bravo"] = 1
    labelMap["charlie"] = 2
    labelMap["delta"] = 3

	rdb, err := database.NewRoutingDatabase("NATO", "tcp", address, labelMap)
    if err != nil {
        panic(err)
    }

    // this ip address is that of a virtual machine hosted on
    // phoenix at cosi

    rdb.SetTrivialPaths()
    for source, a := range labelMap {
        for destination, b := range labelMap {
            if source != destination {
                rdb.SetPath(source, destination, fmt.Sprintf("%s %s | %d", source, destination, a + b))
            }
        }
    }
/*
	for i := 0; i < numberOfNodes; i++ {
		for j := 0; j < numberOfNodes; j++ {
			if i != j {
                rdb.SetPath(i, j, fmt.Sprintf("%d %d | %d", i, j, i + j))
            }
        }
	}
*/
    rdb.StoreLabelsInDB()
    rdb.StorePathsInDB()

	var src, dest, path string

	for {
		fmt.Print("enter the source label > ")
        fmt.Scanln(&src)
        fmt.Print("enter the destination label > ")
        fmt.Scanln(&dest)
        path, err = rdb.GetPathFromDB(src, dest)
        if err != nil {
            fmt.Println(err)
            continue
        }
		fmt.Printf("the shortest path from %s to %s is: %s\n", src, dest, path)
	}

}
