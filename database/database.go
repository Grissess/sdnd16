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
	"strconv"
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
func NewRoutingDatabase(dbName string, network string, address string, ids map[string]int, paths map[int]map[int]string) (RoutingDatabase, error) {
    rdb := RoutingDatabase{name: dbName, connection: nil, ids: make(map[string]int), labels: make(map[int]string), initialized: false, size: -1}
	err := rdb.Connect(network, address)
	if err != nil {
		return rdb, err
	}
	rdb.setIds(ids)
    rdb.setPaths(paths)
    rdb.initialized = true
    return rdb, nil
}

func EraseRoutingDatabase(dbName string, network string, address string) error {
    db, err := redis.Dial(network, address)
    if err != nil {
        return err
    }
    _, err = db.Do("DEL", dbName)
    if err != nil {
        return err
    }
    db.Close()
    return nil
}

func DatabaseExists(dbName string, network string, address string) (bool, error) {
    db, err := redis.Dial(network, address)
    if err != nil {
        return false, err
    }
    var reply string
    reply, err = redis.String(db.Do("EXISTS", dbName))
    result := (reply == "1")
    if err != nil {
        return false, err
    }
    db.Close()
    return result, nil
}
/*
func ConnectToRoutingDatabase(dbName string, network string, address string, ) (RoutingDatabase, error) {
}
*/
// Set the map of labels in the local graph data.
func (self *RoutingDatabase) setIds(nodeIds map[string]int) {
	self.ids = nodeIds
	self.size = len(nodeIds)
    for k, v := range self.ids {
        self.labels[v] = k
    }
}

// Connect to a redis database specified by a protocol (network) and address.
func (self *RoutingDatabase) Connect(network string, address string) error {
	db, err := redis.Dial(network, address)
	if err == nil {
		self.connection = db
		self.initialized = true
	}
	return err
}

// Disconnect from the connected database.
func (self *RoutingDatabase) Disconnect() error {
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
	if !self.initialized {
		return uninitializedPath, errors.New("RoutingDatabase: no initialized connection to get path from")
	}
	if !self.initialized {
		return uninitializedPath, errors.New("RoutingDatabase: no labels for topology")
	}
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
	if !self.initialized {
		return errors.New("RoutingDatabase: not initialized")
	}
	var sizeStr, nodeLabel string
	var size int
	sizeStr, err := redis.String(self.connection.Do("HGET", self.name, "SIZE"))
	if err != nil {
		return err
	}
	size, _ = strconv.Atoi(sizeStr)
	for i := 0; i < size; i++ {
		nodeLabel, err = redis.String(self.connection.Do("HGET", self.name, fmt.Sprintf("{%d}", i)))
		if err != nil {
			return err
		}
		self.ids[nodeLabel] = i
	}
	self.initialized = true
	return nil
}

// Store the local copy of the label map in the connected redis database.
func (self *RoutingDatabase) storeLabels() error {
	if !self.initialized {
		return errors.New("RoutingDatabase: no connected database to store paths in")
	}
	if !self.initialized {
		return errors.New("RoutingDatabase: no labels for topology")
	}
	self.connection.Do("HSET", self.name, "SIZE", fmt.Sprintf("%d", len(self.labels)))
	for key, index := range self.labels {
		self.connection.Do("HSET", self.name, fmt.Sprintf("{%d}", index), key)
	}
	return nil
}
