//The reader package contains file reading and other utility functions.
package utils

/*
* CS350 Team 5 - File reading and  related functions
* Michael Fulton
 */

import (
	"bufio"
	"github.com/gonum/graph"
	"github.com/gonum/graph/simple"
	"os"
	"math"
	"strconv"
)

//Reads in a topology file with the structure src,dst,weight for every edge
func ReadFileToGraph(filename string) (*simple.UndirectedGraph, map[int]string, error) {

	f, err := os.Open(filename)

	if err != nil {
		return nil, nil, err
	}

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanWords)
	var result []string

	for scanner.Scan() {
		x := scanner.Text()
		result = append(result, x)
	}

	var srcs, dsts, weights []string

	for i := 0; i < len(result); i = i + 3 {
		srcs = append(srcs, result[i])
	}
	for i := 1; i < len(result); i = i + 3 {
		dsts = append(dsts, result[i])
	}

	for i := 2; i < len(result); i = i + 3 {
		weights = append(weights, result[i])
	}

	g := simple.NewUndirectedGraph(0.0, math.Inf(1))
	labels := make(map[int]string)
	revlabels := make(map[string]int)
	ctr := 0;

	for i := 0; i < len(srcs); i++ {
		sid, sok := revlabels[srcs[i]];
		did, dok := revlabels[dsts[i]];
		if !sok {
			sid = ctr;
			g.AddNode(simple.Node(sid));
			labels[sid] = srcs[i];
			revlabels[srcs[i]] = sid;
			ctr++;
		}
		if !dok {
			did = ctr;
			g.AddNode(simple.Node(did));
			labels[did] = dsts[i];
			revlabels[dsts[i]] = did;
			ctr++;
		}

		cost, err1 := strconv.ParseFloat(weights[i], 64)
		if err1 != nil {
			return g, labels, err1
		}

		g.SetEdge(simple.Edge{F: simple.Node(sid), T: simple.Node(did), W: cost});
	}

	return g, labels, nil
}

// Reverses the label map (id est, reverse lookups from name to ID)
func GetRevLabels(labels map[int]string) map[string]int {
	revlabels := make(map[string]int);
	for k, v := range(labels) {
		revlabels[v] = k;
	}
	return revlabels;
}

//Returns an array of strings containing the labels of all nodes in the graph.
// This should be a (possibly improper) subset of the values of labels.
func GetLabelList(g graph.Graph, labels map[int]string) []string {
	nodes := g.Nodes();

	ret := make([]string, 0, len(nodes))
	for _, node := range nodes {
		ret = append(ret, labels[node.ID()])
	}

	return ret
}

//Returns a map of strings mapping the labels of nodes to their neighbors.
func GetNeighborMap(g graph.Graph) map[int]map[int]int {
	nodes := g.Nodes();
	neighborMap := make(map[int]map[int]int);
	for _, node := range(nodes) {
		neighborMap[node.ID()] = make(map[int]int);
		neighbors := g.From(node);
		for _, neighbor := range neighbors {
			neighborMap[node.ID()][neighbor.ID()] = int(g.Edge(node, neighbor).Weight());
		}
	}

	return neighborMap
}

// Reconstruct a graph from adjacency maps
func GraphFromNeighborMap(adjmap map[int]map[int]int) graph.Graph {
	g := simple.NewUndirectedGraph(0.0, math.Inf(1));
	for srcid, neighmap := range(adjmap) {
		if !g.Has(simple.Node(srcid)) {
			g.AddNode(simple.Node(srcid))
		}
		for dstid, cost := range(neighmap) {
			if !g.Has(simple.Node(dstid)) {
				g.AddNode(simple.Node(dstid))
			}
			g.SetEdge(simple.Edge{F: simple.Node(srcid), T: simple.Node(dstid), W: float64(cost)})
		}
	}
	return g
}
