//The reader package contains file reading files and other interfaceing functions for
// CS350 Team 5's Graph Project
package reader

/*
* CS350 Team 5 - File reading and  related functions
* Michael Fulton
 */

import (
	"bufio"
	"github.com/gyuho/goraph/graph"
	"os"
        "log"
	"strconv"
)

//Reads in a topology file with the structure src,dst,weight for every edge
func ReadFile(filename string) *graph.DefaultGraph {

	f, err := os.Open(filename)

        if err != nil {
                log.Fatal(err)
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
	        g.AddVertex(srcs[i]);
		g.AddVertex(dests[i]);

		cost, _ := strconv.Atoi(weights[i]);

		g.AddEdge(srcs[i], dests[i], float64(cost));
		g.AddEdge(dests[i], srcs[i], float64(cost));
	}

	return g;
}

//Returns an array of strings containing the labels of all nodes in the graph.
func LabelList(g *graph.DefaultGraph) map[string]int {
	nodes := g.GetVertices()

	var node_labels map[string]int
	for i := 0; i < len(nodes); i = i + 1 {
	}

	return node_labels
}
