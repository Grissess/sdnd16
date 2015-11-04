//The reader package contains file reading and other utility functions.
package reader

/*
* CS350 Team 5 - File reading and  related functions
* Michael Fulton
 */

import (
	"bufio"
	"github.com/gyuho/goraph/graph"
	"os"
	"strconv"
	"strings"
)

//Reads in a topology file with the structure src,dst,weight for every edge
func ReadFileToGraph(filename string) (*graph.DefaultGraph, error) {

	f, err := os.Open(filename)

	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanWords)
	var result []string

	for scanner.Scan() {
		x := scanner.Text()
		result = append(result, x)
	}

	var srcs, dests, weights []string

	for i := 0; i < len(result); i = i + 3 {
		srcs = append(srcs, result[i])
	}
	for i := 1; i < len(result); i = i + 3 {
		dests = append(dests, result[i])
	}

	for i := 2; i < len(result); i = i + 3 {
		weights = append(weights, result[i])
	}

	g := graph.NewDefaultGraph()

	for i := 0; i < len(srcs); i = i + 1 {
		g.AddVertex(srcs[i])
		g.AddVertex(dests[i])

		cost, err1 := strconv.Atoi(weights[i])
		if err1 != nil {
			return g, err1
		}

		err2 := g.AddEdge(srcs[i], dests[i], float64(cost))
		err3 := g.AddEdge(dests[i], srcs[i], float64(cost))
		if err2 != nil {
			return g, err2
		}
		if err3 != nil {
			return g, err3
		}
	}

	return g, nil
}

//Returns an array of strings containing the labels of all nodes in the graph.
func GetLabelList(g *graph.DefaultGraph) []string {
	vertices := g.GetVertices()

	labels := make([]string, 0, len(vertices))
	for key := range vertices {
		labels = append(labels, key)
	}

	return labels
}

//Returns a map of strings, mapping the labels of each node in the graph to a unique number
func GetLabelMap(g *graph.DefaultGraph) map[string]int {
	vertices := g.GetVertices()

	labels := make([]string, 0, len(vertices))
	for key := range vertices {
		labels = append(labels, key)
	}

	node_labels := make(map[string]int)
	for i := 0; i < len(labels); i = i + 1 {
		node_labels[labels[i]] = i
	}

	return node_labels
}

//Returns a map of strings mapping the labels of nodes to their neighbors.
func GetNeighborMap(g *graph.DefaultGraph) (map[string]string, error) {
	labels := GetLabelList(g)
	neighborMap := make(map[string]string)
	for i := 0; i < len(labels); i = i + 1 {
		neighbors, err := g.GetParents(labels[i])
		if err != nil {
			return neighborMap, err
		}
		neighborLabels := make([]string, 0, len(neighbors))
		for key := range neighbors {
			neighborLabels = append(neighborLabels, key)
		}

		neighborMap[labels[i]] = strings.Join(neighborLabels, " ")
	}

	return neighborMap, nil
}
