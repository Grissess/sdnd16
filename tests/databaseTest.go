/*
 * CS350 Team 5 - shortests route for a given network topology
 * Iteration 3 - graph data storage
 * David Josephs and Killian Coddington
 */

package main

import (
	"fmt"
	"github.com/Grissess/sdnd16/database"
)

func main() {

	address := "128.153.144.171:6379"
	dbName := "NATO"

	idMap := make(map[string]int)
	idMap["alpha"] = 0
	idMap["bravo"] = 1
	idMap["charlie"] = 2
	idMap["delta"] = 3

    topologyMap := make(map[int]map[int]int)
    for i := 0; i < 4; i++ {
        topologyMap[i] = make(map[int]int)
        for j := 0; j < 4; j++ {
            topologyMap[i][j] = (i + j)
        }
    }

	pathMap := make(map[int]map[int]string)
	for i := 0; i < 4; i++ {
		pathMap[i] = make(map[int]string)
		for j := 0; j < 4; j++ {
			pathMap[i][j] = "go"
		}
	}

	database.EraseDatabase(dbName, "tcp", address)

	result, err := database.DatabaseExists(dbName, "tcp", address)

	if result {
		panic("DatabaseExists sucks!")
	}

	if err != nil {
		panic(err)
	}

    rdb, err := database.NewRoutingDatabase(dbName, "tcp", address, idMap, pathMap, topologyMap)

	if err != nil {
		panic(err)
	}

	result, err = database.DatabaseExists(dbName, "tcp", address)

	if !result {
		panic("DatabaseExists sucks!")
	}

	/*

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

		for i := 0; i < numberOfNodes; i++ {
			for j := 0; j < numberOfNodes; j++ {
				if i != j {
	                rdb.SetPath(i, j, fmt.Sprintf("%d %d | %d", i, j, i + j))
	            }
	        }
		}

	    rdb.StoreLabelsInDB()
	    rdb.StorePathsInDB()
	*/

	var src, dest, path string

	for {
		fmt.Print("enter the source label > ")
		fmt.Scanln(&src)
        fmt.Print("enter the destination label > ")
		fmt.Scanln(&dest)
        if src == "quit" || dest == "quit" {
            break
        }
		path, err = rdb.GetPath(src, dest)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Printf("the shortest path from %s to %s is: %s\n", src, dest, path)
	}

    err = rdb.CloseConnection()
    if err != nil {
        panic(err)
    }
	database.EraseDatabase(dbName, "tcp", address)

}
