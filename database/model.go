/*
 * CS350 Team 5 - shortests route for a given network topology
 * Iteration 2 - graph data storage
 * David Josephs and Killian Coddington
 */

package database

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
)

type DsPathData struct {
	dsPaths [][]string
	dsNodes int
}

func NewPathData(nodes int) struct {
	paths := make([][]string, nodes)
	for i := range paths {
		paths[i] = make([]string, nodes)
	}
	return DsPathData {dsPaths: paths, dsNodes: nodes}
}

	// connect to running redis server via TCP on port 6379
	db, err := redis.Dial("tcp", ":6379")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// loop through all stored path info and store in redis database in a hashmap
	for i := range paths {
		for j := range paths[i] {
			db.Do("HSET", "paths", fmt.Sprintf("s%d:d%d", i, j), paths[i][j])
		}
	}

	// variables for the source and destination nodes when asking for a path
	var s, d int

	// ask the user to specify valid source and destination nodes to retrieve the shortest path
	for {
		fmt.Print("enter the pair of nodes you want a path between (s, d) > ")
		fmt.Scanf("(%d, %d)\n", &s, &d)
		if s > n || s < 0 || d > n || d < 0 {
			fmt.Printf("at least one specified node is not in the range 1 - %d\n", n)
			continue
		}

		path, err := redis.String(db.Do("HGET", "paths", fmt.Sprintf("s%d:d%d", s, d)))
		if err != nil {
			fmt.Println("key not found")
		}
		fmt.Printf("the shortest path from %d to %d is: %s\n", s, d, path)
	}

}
