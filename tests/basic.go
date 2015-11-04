//Test cases used for developing and debuging this project
package main

import(
        "github.com/Grissess/sdnd16/reader"
        "github.com/Grissess/sdnd16/network"
        "fmt"
)

func main(){
		g := reader.ReadFile("topology.txt");
		start := g.GetAllNodes()[0];
		sg, err := g.Search(start, func(edge *network.DsEdge) int { return edge.GetAttr("cost").(int); });
		sg.ToDot();
        fmt.Println(g.ToDot());
		//fmt.Printf("(%p)\n", g);
		//fmt.Println(start);
		//fmt.Printf("(%p of %p)\n", start, start.GetGraph());
		if(err == nil) {
			//fmt.Println(sg.ToDot());
		} else {
			//fmt.Println("Error: ", err);
		}
}
