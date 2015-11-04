package main

import(
	"fmt"
	"github.com/gyuho/goraph/graph"
	"github.com/Grissess/sdnd16/reader"
	"github.com/Grissess/sdnd16/multipath"
)

func main(){
        var g *graph.DefaultGraph
        g, _ = reader.ReadFileToGraph("topology.txt")
        fmt.Println(g.String())

		paths := multipath.Yen(g, "1", "5", 5);
        fmt.Println(paths)

		fmt.Println(g.String())
}
