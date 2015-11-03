package main

import(
        "github.com/Grissess/sdnd16/reader"
        "github.com/Grissess/sdnd16/network"
        "fmt"
)

func main(){
        var g network.DsGraph;
        g= reader.ReadFile("topology.txt");

        fmt.Print(g.String());
        fmt.Println();

        node_labels := reader.LabelList(g);
        for i:= 0; i< len(node_labels); i = i + 1 {
                fmt.Println(node_labels[i]);
        }
}
