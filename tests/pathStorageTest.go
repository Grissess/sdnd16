/*
 * CS350 Team 5 - shortests route for a given network topology
 * Iteration 2 - graph data storage
 * David Josephs and Killian Coddington
 */

package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
)

func main() {

	// read in the size of the network topology
	fmt.Print("enter the number of nodes in the network topology > ")
	var n int
	_, err := fmt.Scanf("%d", &n)

	// create the 2D array storing nodes and the shortest path to each other node
	paths := make([][]string, n)
	for i := range paths {
		paths[i] = make([]string, n)
	}

	// set each value in the shortest paths data structure to be a temporary value
	for i := range paths {
		for j := range paths[i] {
			paths[i][j] = fmt.Sprintf("%d %d | %d", j, i, (i + j))
		}
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
