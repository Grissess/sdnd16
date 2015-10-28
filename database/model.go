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
	dsNumNodes int
}

func NewPathData(numNodes int) struct {
	paths := make([][]string, numNodes)
	for i := range paths {
		paths[i] = make([]string, numNodes)
	}
	return DsPathData {dsPaths: paths, dsNumNodes: numNodes}
}

func (self *DsPathData) SetPath(source int, destination int, path string) {
	if source < self.dsNumNodes && destination < self.dsNumNodes && source >= 0 && destination >= 0
		self.dsPaths[source][destination] = path
}

func (self *DsPathData) GetPaths() [][]string {
	return self.dsPaths
}

func (self *DsPathData) GetPath(source int, destination int) string {
	if source < self.dsNumNodes && destination < self.dsNumNodes && source >= 0 && destination >= 0
		return self.dsPaths[source][destination]
	else
		return "invalid source and/or destination"
}

func (self *DsPathData) GetNumNodes() int {
	return self.dsNumNodes
}

type DsTopologyData struct {
	dsGraph []string
	dsNumNodes int
}

func NewTopologyData(numNodes int) struct {
	graph := make([]string, numNodes)
	return DsTopologyData {dsGraph: graph, dsNumNodes: numNodes}
}

func (self *DsTopologyData) SetNodeNeighbors(node int, neighbors string) {
	if node < self.dsNumNodes && node >= 0
		self.dsGraph[node] = neighbors
}

func (self *DsTopologyData) GetNodeNeighbors(node int) string {
	if node < self.dsNumNodes && node >= 0
		return self.dsGraph[node]
	else
		return "invalide node"
}

type DsTopologyDatabase struct {
	dsConnection Conn
	dsPaths DsPathData
	dsTopology DsTopologyData
}

func NewTopologyDatabase(numNodes int, network string, address string) struct {
	db, err := redis.Dial(network, address)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	newTopology := NewTopologyData(numNodes)
	newPaths := NewPathData(numNodes)
	return DsTopologyDatabase{dsConnection: db, dsPaths: newPaths, dsTopology: newTopology}
}

/*
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
*/
