/*
* CS350 Team 5 - shortests route for a given network topology
* David Josephs and Killian Coddington
*/

// A package developed to ease interaction between topology information and a redis database
package database

import (
	"errors"
	"fmt"
    "strings"
    "strconv"
    "github.com/garyburd/redigo/redis"
)

const uninitializedPath = ""

// Structure devoted to managing a connection to a redis database
// Functions RoutingDatabase provides a means to store a set of paths for a topology in a redis database
type RoutingDatabase struct {
	name        string
	connection  redis.Conn
	ids         map[string]int
	labels      map[int]string
	initialized bool
	size        int
}

// Create a new routing database structure with a given name
// a connection to a redis database specified by a network type and an ip address
// a map from strings to ints indicatating ids
// a map from int to map of ints to strings representing paths
// a map similar to the last but to ints indicating a topology
func NewRoutingDatabase(dbName string, network string, address string, ids map[string]int, paths map[int]map[int]string, topology map[int]map[int]int) (RoutingDatabase, error) {
	rdb := RoutingDatabase{name: dbName, connection: nil, ids: make(map[string]int), labels: make(map[int]string), initialized: false, size: -1}
	err := rdb.connect(network, address)
	if err != nil {
		return rdb, err
	}
	rdb.saveName()
    rdb.setIds(ids)
	rdb.setLabels()
	rdb.setPaths(paths)
    rdb.setTopology(topology)
	rdb.initialized = true
	return rdb, nil
}

// Private function to set a databases topology information
func (self *RoutingDatabase) setTopology(topology map[int]map[int]int) {
    for i, neighmap := range topology {
        neighbors := make([]string, 0, 2 * len(neighmap))
        for neighbor, cost := range neighmap {
			neighbors = append(neighbors, strconv.Itoa(neighbor))
			neighbors = append(neighbors, strconv.Itoa(cost))
        }
        self.connection.Do("HSET", self.name, fmt.Sprintf("|%d|", i), strings.Join(neighbors, " "))
    }
}

// Function to return a topology in the format it is provided in
func (self *RoutingDatabase) GetTopology() (map[int]map[int]int, error) {
    topology := make(map[int]map[int]int)
    for i := 0; i < self.size; i++ {
        neighbor, err := redis.String(self.connection.Do("HGET", self.name, fmt.Sprintf("|%d|", i)))
        if err != nil {
            return topology, err
        }
		topology[i] = make(map[int]int)
        array := strings.Split(neighbor, " ")
        for j := 0; j < (len(array) / 2); j++ {
            a, _ := strconv.Atoi(array[2*j])
            b, _ := strconv.Atoi(array[(2*j) + 1])
            topology[i][a] = b
        }
    }
    return topology, nil
}

// Private function to add the current topology's name so all topologies may be enumerated
func (self *RoutingDatabase) saveName() {
    self.connection.Do("SADD", "{topologies}", self.name)
}

// Enumerate a list of all topologies by name
func GetAllTopologies(network string, address string) ([]string, error) {
    db, err := redis.Dial(network, address)
	if err != nil {
		return nil, err
	}
    topologies, err := redis.Strings(db.Do("SMEMBERS", "{topologies}"))
	if err != nil {
		return nil, err
	}
	db.Close()
	return topologies, nil
}

// Safely close a connection to the current redis database
func (self *RoutingDatabase) CloseConnection() error {
	err := self.disconnect()
	return err
}

// Erase all keys associated with a topology and remove its name from the topology list
func EraseDatabase(dbName string, network string, address string) error {
	db, err := redis.Dial(network, address)
	if err != nil {
		return err
	}
	_, err = db.Do("DEL", dbName)
    if err != nil {
		return err
	}
    _, err = db.Do("SREM", "{topologies}", dbName)
	db.Close()
	return err
}

// Verify if a given topology exists within the specified database
func TopologyExists(dbName string, network string, address string) (bool, error) {
	db, err := redis.Dial(network, address)
	if err != nil {
		return false, err
	}
	var reply int
	reply, err = redis.Int(db.Do("EXISTS", dbName))
	result := (reply == 1)
	if err != nil {
		return false, err
	}
	db.Close()
	return result, nil
}

// Connect to a preexisting database
func ConnectToDatabase(dbName string, network string, address string) (RoutingDatabase, error) {
    conn, err := redis.Dial(network, address)
    rdb := RoutingDatabase{name: dbName, connection: conn, ids: make(map[string]int), labels: make(map[int]string), initialized: false, size: -1}
    if err != nil {
        return rdb, err
    }
    result, err := TopologyExists(dbName, network, address)
    if err != nil {
        return rdb, err
    }
    if !result {
		return rdb, errors.New("RoutingDatabase: Specified database does not exist")
	}
    err = rdb.getLabels()
    rdb.initialized = true
	return rdb, err
}

// Set the map of ids in the database
func (self *RoutingDatabase) setIds(nodeIds map[string]int) {
	self.ids = nodeIds
	self.size = len(nodeIds)
	for k, v := range self.ids {
		self.labels[v] = k
	}
}

// Connect to a redis database specified by a protocol (network) and address
func (self *RoutingDatabase) connect(network string, address string) error {
	db, err := redis.Dial(network, address)
	if err == nil {
		self.connection = db
		self.initialized = true
	}
	return err
}

// Disconnect from the connected database
func (self *RoutingDatabase) disconnect() error {
	if self.initialized {
		self.connection.Close()
		return nil
	}
	return errors.New("RoutingDatabase: no initialized connection to disconnect from")
}

// Get the number of nodes for the topology a RoutingDatabase respresents
func (self *RoutingDatabase) GetSize() int {
	return self.size
}

// Get path information for a specific path from a redis database and topology
func (self *RoutingDatabase) GetPath(source string, destination string) (string, error) {
	s, okSource := self.ids[source]
	d, okDestination := self.ids[destination]
	if okSource && okDestination {
		path, err := redis.String(self.connection.Do("HGET", self.name, fmt.Sprintf("{%d:%d}", s, d)))
		if err != nil {
			return uninitializedPath, errors.New("RoutingDatabase: path not found")
		}
		return path, nil
	} else {
		if !okSource {
			return uninitializedPath, errors.New(fmt.Sprintf("RoutingDatabase: invalid source label provided (%s)", source))
		}
		if !okDestination {
			return uninitializedPath, errors.New(fmt.Sprintf("RoutingDatabase: invalid destination label provided (%s)", destination))
		}
		return uninitializedPath, errors.New("RoutingDatabase: fatal error (you broke boolean algebra)")
	}
}

// Set paths using a given map, private function to be run when creating a topology in a database
func (self *RoutingDatabase) setPaths(paths map[int]map[int]string) {
	for i := range paths {
		for j := range paths[i] {
			self.connection.Do("HSET", self.name, fmt.Sprintf("{%d:%d}", i, j), paths[i][j])
		}
	}
}

// Get the labels from the connected database (also ids) private function called when connecting to an existing topology 
func (self *RoutingDatabase) getLabels() error {
	var nodeLabel string
	size, err := redis.Int(self.connection.Do("HGET", self.name, "{size}"))
	if err != nil {
		return err
	}
	for i := 0; i < size; i++ {
		nodeLabel, err = redis.String(self.connection.Do("HGET", self.name, fmt.Sprintf("{%d}", i)))
		if err != nil {
			return err
		}
		self.ids[nodeLabel] = i
		self.labels[i] = nodeLabel
	}
	return nil
}

// Set all of the label information for a topology in the database, again private for making new topologies
func (self *RoutingDatabase) setLabels() error {
	self.connection.Do("HSET", self.name, "{size}", fmt.Sprintf("%d", len(self.labels)))
	for key, index := range self.labels {
		self.connection.Do("HSET", self.name, fmt.Sprintf("{%d}", key), index)
	}
	return nil
}
