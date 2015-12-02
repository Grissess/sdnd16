/*
* CS350 Team 5 - shortests route for a given network topology
* David Josephs and Killian Coddington
*/

// A package developed to ease interaction between topology information and a redis database.
package database

import (
	"errors"
	"fmt"
    "strings"
    "strconv"
    "github.com/garyburd/redigo/redis"
)

const uninitializedPath = ""

// Structure devoted to storing paths for a topology, as well as a connection to a redis database.
// Functions RoutingDatabase provides a means to store a set of paths for a topology in a redis database.
type RoutingDatabase struct {
	name        string
	connection  redis.Conn
	ids         map[string]int
	labels      map[int]string
	initialized bool
	size        int
}

// Create a new routing database structure with a given name,
// a connection to a redis database specified by a network type and an ip address,
// and a map from strings to ints indicating unique ids for string labels.
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

func (self *RoutingDatabase) setTopology(topology map[int]map[int]int) {
    for i := range topology {
        neighbors := make([]string, 2 * len(topology[i]))
        for j := range topology[i] {
            neighbors[(2*j)] = strconv.Itoa(j)
            neighbors[(2*j)+1] = strconv.Itoa(topology[i][j])
        }
        self.connection.Do("HSET", self.name, fmt.Sprintf("|%d|", i), strings.Join(neighbors, " "))
    }
}

func (self *RoutingDatabase) GetTopology() (map[int]map[int]int, error) {
    topology := make(map[int]map[int]int)
    for i := 0; i < self.size; i++ {
        neighbor, err := redis.String(self.connection.Do("HGET", self.name, fmt.Sprintf("|%d|", i)))
        if err != nil {
            return topology, err
        }
        array := strings.Split(neighbor, " ")
        for j := 0; j < (len(array) / 2); j++ {
            a, _ := strconv.Atoi(array[2*j])
            b, _ := strconv.Atoi(array[(2*j) + 1])
            topology[a][j] = b
        }
    }
    return topology, nil
}

func (self *RoutingDatabase) saveName() {
    self.connection.Do("SADD", "{topologies}", self.name)
}

func GetAllDatabases(network string, address string) ([]string, error) {
    db, err := redis.Dial(network, address)
	if err != nil {
		return nil, err
	}
    databases, err := redis.Strings(db.Do("SMEMBERS", "{topologies}"))
	if err != nil {
		return nil, err
	}
	db.Close()
	return databases, nil
}

func (self *RoutingDatabase) CloseConnection() error {
	err := self.disconnect()
	return err
}

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

func DatabaseExists(dbName string, network string, address string) (bool, error) {
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

func ConnectToDatabase(dbName string, network string, address string) (RoutingDatabase, error) {
    conn, err := redis.Dial(network, address)
    rdb := RoutingDatabase{name: dbName, connection: conn, ids: make(map[string]int), labels: make(map[int]string), initialized: false, size: -1}
    if err != nil {
        return rdb, err
    }
    result, err := DatabaseExists(dbName, network, address)
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

// Set the map of labels in the local graph data.
func (self *RoutingDatabase) setIds(nodeIds map[string]int) {
	self.ids = nodeIds
	self.size = len(nodeIds)
	for k, v := range self.ids {
		self.labels[v] = k
	}
}

// Connect to a redis database specified by a protocol (network) and address.
func (self *RoutingDatabase) connect(network string, address string) error {
	db, err := redis.Dial(network, address)
	if err == nil {
		self.connection = db
		self.initialized = true
	}
	return err
}

// Disconnect from the connected database.
func (self *RoutingDatabase) disconnect() error {
	if self.initialized {
		self.connection.Close()
		return nil
	}
	return errors.New("RoutingDatabase: no initialized connection to disconnect from")
}

// Get the number of nodes for the topology a RoutingDatabase respresents.
func (self *RoutingDatabase) GetSize() int {
	return self.size
}

// Get path information for a specific path from a redis database.
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

// Locally set all paths from a node to itself with cost 0.
func (self *RoutingDatabase) setPaths(paths map[int]map[int]string) {
	for i := range paths {
		for j := range paths[i] {
			self.connection.Do("HSET", self.name, fmt.Sprintf("{%d:%d}", i, j), paths[i][j])
		}
	}
}

// From the connected redis database obtain a local copy of the label map,
// this is necessary to query for paths.
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

// Store the local copy of the label map in the connected redis database.
func (self *RoutingDatabase) setLabels() error {
	self.connection.Do("HSET", self.name, "{size}", fmt.Sprintf("%d", len(self.labels)))
	for key, index := range self.labels {
		self.connection.Do("HSET", self.name, fmt.Sprintf("{%d}", key), index)
	}
	return nil
}
