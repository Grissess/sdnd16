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

type PathsData struct {
	paths []string
	numNodes int
}

func NewPathsData(numberOfNodes int) PathsData {
	return PathsData{paths: make([]string, numberOfNodes * numberOfNodes), numNodes: numberOfNodes}
}

func (self *PathsData) SetPath(source int, destination int, path string) {
	if source < self.numNodes && destination < self.numNodes && source >= 0 && destination >= 0 {
		self.paths[(source * self.numNodes) + destination] = path
	}
}

func (self *PathsData) GetPath(source int, destination int) string {
	if source < self.numNodes && destination < self.numNodes && source >= 0 && destination >= 0 {
		return self.paths[(source * self.numNodes) + destination]
	} else {
		return "invalid source and/or destination"
	}
}

func (self *PathsData) GetNumNodes() int {
	return self.numNodes
}

/*
type DsTopologyData struct {
	graph []string
	numNodes int
}

func NewTopologyData(numberOfNodes int) DsTopologyData {
	return DsTopologyData {graph: make([]string, (numberOfNodes * numberOfNodes)), numNodes: numberOfNodes}
}

func (self *DsTopologyData) SetNodeNeighbors(node int, neighbors string) {
	if node < self.numNodes && node >= 0 {
		self.graph[node] = neighbors
	}
}

func (self *DsTopologyData) GetNodeNeighbors(node int) string {
	if node < self.numNodes && node >= 0 {
		return self.graph[node]
	} else {
		return "invalide node"
	}
}
*/

type RoutingDatabase struct {
	name string
	connection redis.Conn
	paths PathsData
//	topology DsTopologyData
	connectionInitialized bool
	numNodes int
}

func NewRoutingDatabase(dbName string, numberOfNodes int) RoutingDatabase {
	return RoutingDatabase{name: dbName, connection: nil, paths: NewPathsData(numberOfNodes), /*topology: NewTopologyData(numberOfNodes),*/ connectionInitialized: false, numNodes: numberOfNodes}
}

func (self *RoutingDatabase) SetPath(source int, destination int, path string) {
	self.paths.SetPath(source, destination, path)
}

func (self *RoutingDatabase) GetPath(source int, destination int) string {
	return self.paths.GetPath(source, destination)
}

func (self *RoutingDatabase) Connect(network string, address string) {
	db, err := redis.Dial(network, address)
	if err != nil {
		panic(err)
	}
	//defer db.Close() 
	//idk what this does, probably a good idea to use but i need to research it further
	self.connection = db
	self.connectionInitialized = true
}

func (self *RoutingDatabase) GetNumNodes() int {
	return self.numNodes
}

func (self *RoutingDatabase) DBGetPath(source int, destination int) string {
	if !self.connectionInitialized {
		panic("attempting to request from uninitialized database")
	}
	path, err := redis.String(self.connection.Do("HGET", self.name, fmt.Sprintf("s%d:d%d", source, destination)))
	if err != nil {
		panic("key not found")
	}
	return path
}

func (self *RoutingDatabase) StorePathsData() {
	for i := 0; i < self.numNodes; i++ {
		for j := 0; j < self.numNodes; j++ {
			self.connection.Do("HSET", self.name, fmt.Sprintf("s%d:d%d", i, j), self.paths.GetPath(i, j))
		}
	}
}
