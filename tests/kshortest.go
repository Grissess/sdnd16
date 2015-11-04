package main

import(
	"fmt"
	"github.com/gyuho/goraph/graph"
	"github.com/Grissess/sdnd16/reader"
	"github.com/Grissess/sdnd16/multipath"
)

func main(){
        var g *graph.DefaultGraph
        g = reader.ReadFileToGraph("topology.txt")
        fmt.Print(g.String())

		paths := multipath.Yen(g, "1", "5", 5);
        fmt.Println(paths)
}