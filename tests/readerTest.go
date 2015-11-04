package main

import(
	"fmt"
	"github.com/gyuho/goraph/graph"
        "github.com/Grissess/sdnd16/reader"
)

func main(){
        var g *graph.DefaultGraph
        g = reader.ReadFileToGraph("topology.txt")
        fmt.Print(g.String())

        labels := reader.LabelList(g);
        fmt.Print(labels);

        path, distance, _:= graph.Dijkstra(g, "1", "5")
        fmt.Println(path)
        fmt.Println(distance)
}
