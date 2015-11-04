package main

import(
	"fmt"
	"github.com/gyuho/goraph/graph"
	"github.com/Grissess/sdnd16/reader"
)

func main(){
        var g *graph.DefaultGraph
        g, _ = reader.ReadFileToGraph("topology.txt")
        fmt.Print(g.String())

        labels := reader.GetLabelList(g)
        fmt.Println(labels)

        neighbors , _ :=reader.GetNeighborMap(g)
        fmt.Println(neighbors)
}
