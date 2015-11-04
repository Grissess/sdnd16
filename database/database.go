/*
 * CS350 Team 5 - shortests route for a given network topology
 * Iteration 4 - graph data storage
 * David Josephs and Killian Coddington
 */

// A package developed to ease interaction between topology information and a redis database.
package database

import (
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
)

const uninitializedPath = ""

// Data structure to hold information about the paths between nodes in a topology.
type Paths struct {
	paths    []string
	numNodes int
}

// Create a new Path structure with empty path information and a specified number of nodes.
func NewPaths(numberOfNodes int) Paths {
	return Paths{paths: make([]string, numberOfNodes*numberOfNodes), numNodes: numberOfNodes}
}

// Set a specific path from source to destination nodes.
func (self *Paths) SetPath(source int, destination int, path string) {
	if source < self.numNodes && destination < self.numNodes && source >= 0 && destination >= 0 {
		self.paths[(source*self.numNodes)+destination] = path
	}
}

// Get a specific path between source and destination nodes.
func (self *Paths) GetPath(source int, destination int) string {
	if source < self.numNodes && destination < self.numNodes && source >= 0 && destination >= 0 {
		return self.paths[(source*self.numNodes)+destination]
	} else {
		return uninitializedPath
	}
}

// Return the number of nodes in a given topology Paths is storing paths for.
func (self *Paths) GetNumNodes() int {
	return self.numNodes
}

func (self *Paths) SetTrivial() {
	for i := 0; i < self.numNodes; i++ {
		self.SetPath(i, i, fmt.Sprintf("%d | 0", i))
	}
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

// Structure devoted to storing a Paths for a topology, as well as a connection to a redis database.
// Functions RoutingDatabase provides aim to ease storing a set of paths or a topology in a redis database.
type RoutingDatabase struct {
	name       string
	connection redis.Conn
	paths      Paths
	//	topology DsTopologyData
	connectionInitialized bool
	numNodes              int
}

// Create a new routing database structure with a given name and number of nodes.
func NewRoutingDatabase(dbName string, numberOfNodes int) RoutingDatabase {
	return RoutingDatabase{name: dbName, connection: nil, paths: NewPaths(numberOfNodes) /*topology: NewTopologyData(numberOfNodes),*/, connectionInitialized: false, numNodes: numberOfNodes}
}

// Set a path in the RoutingDatabase's corresponding Paths structure. (local)
func (self *RoutingDatabase) SetPath(source int, destination int, path string) {
	self.paths.SetPath(source, destination, path)
}

// Get a path from a RoutingDatabase's corresponding Paths structure. (local)
func (self *RoutingDatabase) GetPath(source int, destination int) (string, error) {
	path := self.paths.GetPath(source, destination)
	if path == uninitializedPath {
		return path, errors.New("RoutingDatabase: local path requested is uninitialized")
	}
	return path, nil
}

// Connect to a redis database specified by a protocol (network) and address.
func (self *RoutingDatabase) Connect(network string, address string) error {
	db, err := redis.Dial(network, address)
	if err == nil {
		self.connection = db
		self.connectionInitialized = true
	}
	return err
}

func (self *RoutingDatabase) Disconnect() error {
	if self.connectionInitialized {
		self.connection.Close()
		return nil
	}
	return errors.New("RoutingDatabase: no initialized connection to disconnect from")
}

// Get the number of nodes for the topology a RoutingDatabase respresents.
func (self *RoutingDatabase) GetNumNodes() int {
	return self.numNodes
}

// Get path information for a specific path from a redis database. (remote)
func (self *RoutingDatabase) GetPathFromDB(source int, destination int) (string, error) {
	if !self.connectionInitialized {
		return uninitializedPath, errors.New("RoutingDatabase: no initialized connection to get path from")
	}
	path, err := redis.String(self.connection.Do("HGET", self.name, fmt.Sprintf("S%d:D%d", source, destination)))
	if err != nil {
		return uninitializedPath, errors.New("RoutingDatabase: path not found")
	}
	return path, nil
}

func (self *RoutingDatabase) SetTrivialPaths() {
	self.paths.SetTrivial()
}

// Store all local Paths information to a redis database. (remote)
func (self *RoutingDatabase) StorePathsInDB() error {
	if !self.connectionInitialized {
		return errors.New("RoutingDatabase: no connected database to store paths in")
	}
	uninitializedPaths := 0
	for i := 0; i < self.numNodes; i++ {
		for j := 0; j < self.numNodes; j++ {
			path := self.paths.GetPath(i, j)
			if path == uninitializedPath {
				uninitializedPaths++
			}
			self.connection.Do("HSET", self.name, fmt.Sprintf("S%d:D%d", i, j), path)
		}
	}
	if uninitializedPaths != 0 {
		return errors.New(fmt.Sprintf("RoutingDatabase: %d uninitialized paths stored in database", uninitializedPaths))
	}
	return nil
}
