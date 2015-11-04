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
    "strconv"
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

// Return the number of nodes in a given Paths structure.
func (self *Paths) GetNumNodes() int {
    return self.numNodes
}

// Structure devoted to storing paths for a topology, as well as a connection to a redis database.
// Functions RoutingDatabase provides a means to store a set of paths for a topology in a redis database.
type RoutingDatabase struct {
    name       string
    connection redis.Conn
    paths      Paths
    labels     map[string]int
    topology   map[string]string
    labelsInitialized     bool
    connectionInitialized bool
    numNodes              int
}

// Create a new routing database structure with a given name,
// a connection to a redis database specified by a network type and an ip address,
// and a map from strings to ints indicating unique ids for string labels.
func NewRoutingDatabase(dbName string, network string, address string, nodeLabels map[string]int) (RoutingDatabase, error) {
    rdb := RoutingDatabase{name: dbName, connection: nil, paths: NewPaths(0), labels: make(map[string]int), labelsInitialized: false, connectionInitialized: false, numNodes: -1}
    err := rdb.Connect(network, address)
    if err != nil {
        return rdb, err
    }
    rdb.SetLabels(nodeLabels)
    rdb.StoreLabelsInDB()
    rdb.SetTrivialPaths()
    return rdb, nil
}

// Create a new routing database structure for a database with the specified name, and a connection to a database specified as in NewRoutingDatabase.
// This function is intended to be used when grabbing from a databse and not intending to modify the topology.
func NewRoutingDatabaseFromDB(dbName string, network string, address string) (RoutingDatabase, error) {
    rdb := RoutingDatabase{name: dbName, connection: nil, paths: NewPaths(0), labels: make(map[string]int), labelsInitialized: false, connectionInitialized: false, numNodes: -1}
    err := rdb.Connect(network, address)
    if err != nil {
        return rdb, err
    }
    rdb.GetLabelsFromDB()
    rdb.SetTrivialPaths()
    return rdb, nil
}

// Set the map of labels in the local graph data.
func (self *RoutingDatabase) SetLabels(nodeLabels map[string]int) {
    self.labels = nodeLabels
    self.numNodes = len(nodeLabels)
    self.paths = NewPaths(len(nodeLabels))
    self.labelsInitialized = true
}

// Set a path in the local paths data.
func (self *RoutingDatabase) SetPath(source string, destination string, path string) error {
    if !self.labelsInitialized {
        return errors.New("RoutingDatabase: no labels for topology")
    }
    s, okSource := self.labels[source]
    d, okDestination := self.labels[destination]
    if okSource && okDestination {
        self.paths.SetPath(s, d, path)
        return nil
    } else {
        if !okSource {
            return errors.New(fmt.Sprintf("RoutingDatabase: invalid source label provided (%s)", source))
        }
        if !okDestination {
            return errors.New(fmt.Sprintf("RoutingDatabase: invalid destination label provided (%s)", destination))
        }
        return errors.New("RoutingDatabase: fatal error (you broke boolean algebra)")
    }
}

// Get a path from the locally stored paths data.
func (self *RoutingDatabase) GetPath(source string, destination string) (string, error) {
    if !self.labelsInitialized {
        return uninitializedPath, errors.New("RoutingDatabase: no labels for topology")
    }
    s, okSource := self.labels[source]
    d, okDestination := self.labels[destination]
    if okSource && okDestination {
        path := self.paths.GetPath(s, d)
        if path == uninitializedPath {
            return uninitializedPath, errors.New("RoutingDatabase: local path requested is uninitialized")
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

// Connect to a redis database specified by a protocol (network) and address.
func (self *RoutingDatabase) Connect(network string, address string) error {
    db, err := redis.Dial(network, address)
    if err == nil {
        self.connection = db
        self.connectionInitialized = true
    }
    return err
}

// Disconnect from the connected database.
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

// Get path information for a specific path from a redis database.
func (self *RoutingDatabase) GetPathFromDB(source string, destination string) (string, error) {
    if !self.connectionInitialized {
        return uninitializedPath, errors.New("RoutingDatabase: no initialized connection to get path from")
    }
    if !self.labelsInitialized {
        return uninitializedPath, errors.New("RoutingDatabase: no labels for topology")
    }
    s, okSource := self.labels[source]
    d, okDestination := self.labels[destination]
    if okSource && okDestination {
        path, err := redis.String(self.connection.Do("HGET", self.name, fmt.Sprintf("P:S%d:D%d", s, d)))
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
func (self *RoutingDatabase) SetTrivialPaths() {
    for key, index := range self.labels {
        self.paths.SetPath(index, index, fmt.Sprintf("%s | 0", key))
    }
}

// Check if a specific label exists in the label map.
func (self *RoutingDatabase) IsValidLabel(label string) bool {
    if !self.labelsInitialized {
        return false
    }
    _, ok := self.labels[label]
    return ok
}

// From the connected redis database obtain a local copy of the label map,
// this is necessary to query for paths.
func (self *RoutingDatabase) GetLabelsFromDB() error {
    if !self.connectionInitialized {
        return errors.New("RoutingDatabase: no connected database to store paths in")
    }
    var sizeStr, nodeLabel string
    self.connection.Do("HGET", self.name, "L:SIZE", sizeStr)
    size, _ := strconv.Atoi(sizeStr)
    for i := 0; i < size; i++ {
        self.connection.Do("HSET", self.name, fmt.Sprintf("L:%d", i), nodeLabel)
        self.labels[nodeLabel] = i
    }
    self.labelsInitialized = true
    return nil
}

// Store the local copy of the label map in the connected redis database.
func (self *RoutingDatabase) StoreLabelsInDB() error {
    if !self.connectionInitialized {
        return errors.New("RoutingDatabase: no connected database to store paths in")
    }
    if !self.labelsInitialized {
        return errors.New("RoutingDatabase: no labels for topology")
    }
    self.connection.Do("HSET", self.name, "L:SIZE", fmt.Sprintf("%d", len(self.labels)))
    for key, index := range self.labels {
        self.connection.Do("HSET", self.name, fmt.Sprintf("L:%d", index), key)
    }
    return nil
}

// Set the local topology info to a specific map from a node label to its neighbors labels.
func (self *RoutingDatabase) SetTopology(topologyMap map[string]string) {
    self.topology = topologyMap
}

// Store all local topology info in the connected redis database.
func (self *RoutingDatabase) StoreTopologyInDB() error {
    if !self.connectionInitialized {
        return errors.New("RoutingDatabase: no connected database to store paths in")
    }
    if !self.labelsInitialized {
        return errors.New("RoutingDatabase: no labels for topology")
    }
    self.connection.Do("HSET", self.name, "T:SIZE", fmt.Sprintf("%d", len(self.topology)))
    for node, neighbors := range self.topology {
        index, _ := self.labels[node]
        self.connection.Do("HSET", self.name, fmt.Sprintf("T:%d", index), neighbors)
    }
    return nil
}

// Store all local topology info in the connected redis database.
func (self *RoutingDatabase) GetNeighborsFromDB(node string) (string, error) {
    if !self.connectionInitialized {
        return "", errors.New("RoutingDatabase: no connected database to store paths in")
    }
    if !self.labelsInitialized {
        return "", errors.New("RoutingDatabase: no labels for topology")
    }
    var neighbors string
    index, _ := self.labels[node]
    self.connection.Do("HGET", self.name, fmt.Sprintf("T:%d", index), neighbors)
    return neighbors, nil
}

// Store all local path information to a redis database.
func (self *RoutingDatabase) StorePathsInDB() error {
    if !self.connectionInitialized {
        return errors.New("RoutingDatabase: no connected database to store paths in")
    }
    if !self.labelsInitialized {
        return errors.New("RoutingDatabase: no labels for topology")
    }
    uninitializedPaths := 0
    for _, i := range self.labels {
        for _, j := range self.labels {
            path := self.paths.GetPath(i, j)
            if path == uninitializedPath {
                uninitializedPaths++
            }
            self.connection.Do("HSET", self.name, fmt.Sprintf("P:S%d:D%d", i, j), path)
        }
    }
    if uninitializedPaths != 0 {
        return errors.New(fmt.Sprintf("RoutingDatabase: %d uninitialized paths stored in database", uninitializedPaths))
    }
    return nil
}
