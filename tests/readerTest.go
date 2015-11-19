package main

import(
	"fmt"
	"github.com/gyuho/goraph/graph"
	"github.com/Grissess/sdnd16/utils"
)

func main(){
        var g *graph.DefaultGraph
        g, _ = utils.ReadFileToGraph("topology.txt")
        g = utils.ReadFile("topology.txt")
        fmt.Print(g.String())

        labels := utils.GetLabelList(g)
        fmt.Println(labels)

        neighbors , _ :=utils.GetNeighborMap(g)
        fmt.Println(neighbors)
}
