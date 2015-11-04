package main

import(
	"fmt"
	"github.com/gyuho/goraph/graph"
        "github.com/Grissess/sdnd16/reader"
)


func main(){
        var g *graph.DefaultGraph
        g = reader.ReadFile("topology.txt") 
        fmt.Print(g.String())
}
