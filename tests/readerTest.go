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

        vertices := g.GetVertices()
        labels := make([]string, 0, len(vertices))
        for key := range vertices{
                labels = append(labels, key);
        }

        path, distance, _:= graph.Dijkstra(g, "1", "5")
        fmt.Println(path)
        fmt.Println(distance)
}
