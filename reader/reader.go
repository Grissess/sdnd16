//The reader package contains file reading files and other interfaceing functions for
// CS350 Team 5's Graph Project
package reader

/*
* CS350 Team 5 - File reading and  related functions
* Michael Fulton
 */

import (
	"bufio"
	"github.com/Grissess/sdnd16/network"
	"os"
        "log"
	"strconv"
	"fmt"
)

//Reads in a topology file with the structure src,dst,weight for every edge
func ReadFile(filename string) *network.DsGraph{

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

	g := network.NewGraph()

	for i := 0; i < len(srcs); i = i + 1 {
		s := g.GetOrCreateNode(network.Label(srcs[i]))
		d := g.GetOrCreateNode(network.Label(dests[i]))
		e1, err1 := g.NewEdge(s,d);
		e2, err2 := g.NewEdge(d,s);
		if err1 != nil {
			fmt.Printf("// ERROR: Creating edge: %s", err1);
		}
		if err2 != nil {
			fmt.Printf("// ERROR: Creating edge: %s", err2);
		}
		cost, _ := strconv.Atoi(weights[i]);
		e1.SetAttr("cost", cost);
		e2.SetAttr("cost", cost);
	}

	return g;
}

//Returns an array of strings containing the labels of all nodes in the graph.
func LabelList(g *network.DsGraph) map[string]int {
	nodes := g.GetAllNodes()

	var node_labels map[string]int
	for i := 0; i < len(nodes); i = i + 1 {
		node_labels[nodes[i].GetLabel().String()] = i;
	}

	return node_labels
}
