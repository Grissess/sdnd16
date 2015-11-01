/*
 * CS350 Team 5 - shortests route for a given network topology
 * Iteration 2 - graph data storage
 * David Josephs and Killian Coddington
 */

// A package developed to ease interaction between topology information and a redis database.
package database

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
)

// Data structure to hold information about the paths between nodes in a topology.
type PathsData struct {
	paths []string
	numNodes int
}

// Create a new PathsData structure with empty path information and a specified number of nodes.
func NewPathsData(numberOfNodes int) PathsData {
	return PathsData{paths: make([]string, numberOfNodes * numberOfNodes), numNodes: numberOfNodes}
}

// Set a specific path from source to destination nodes.
func (self *PathsData) SetPath(source int, destination int, path string) {
	if source < self.numNodes && destination < self.numNodes && source >= 0 && destination >= 0 {
		self.paths[(source * self.numNodes) + destination] = path
	}
}

// Get a specific path between source and destination nodes.
func (self *PathsData) GetPath(source int, destination int) string {
	if source < self.numNodes && destination < self.numNodes && source >= 0 && destination >= 0 {
		return self.paths[(source * self.numNodes) + destination]
	} else {
		return "invalid source and/or destination"
	}
}

// Return the number of nodes in a given topology PathsData is storing paths for.
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

// Structure devoted to storing a PathsData for a topology, as well as a connection to a redis database.
// Functions RoutingDatabase provides aim to ease storing a set of paths or a topology in a redis database.
type RoutingDatabase struct {
	name string
	connection redis.Conn
	paths PathsData
//	topology DsTopologyData
	connectionInitialized bool
	numNodes int
}

// Create a new routing database structure with a given name and number of nodes.
func NewRoutingDatabase(dbName string, numberOfNodes int) RoutingDatabase {
	return RoutingDatabase{name: dbName, connection: nil, paths: NewPathsData(numberOfNodes), /*topology: NewTopologyData(numberOfNodes),*/ connectionInitialized: false, numNodes: numberOfNodes}
}

// Set a path in the RoutingDatabase's corresponding PathsData structure. (local)
func (self *RoutingDatabase) SetPath(source int, destination int, path string) {
	self.paths.SetPath(source, destination, path)
}

// Get a path from a RoutingDatabase's corresponding PathsData structure. (local)
func (self *RoutingDatabase) GetPath(source int, destination int) string {
	return self.paths.GetPath(source, destination)
}

// Connect to a redis database specified by a protocol (network) and address.
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

// Get the number of nodes for the topology a RoutingDatabase respresents.
func (self *RoutingDatabase) GetNumNodes() int {
	return self.numNodes
}

// Get path information for a specific path from a redis database. (remote) 
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

// Store all local PathsData information to a redis database. (remote)
func (self *RoutingDatabase) StorePathsData() {
	for i := 0; i < self.numNodes; i++ {
		for j := 0; j < self.numNodes; j++ {
			self.connection.Do("HSET", self.name, fmt.Sprintf("s%d:d%d", i, j), self.paths.GetPath(i, j))
		}
	}
}
