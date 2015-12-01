package main

import(
	"fmt"
	// "github.com/gonum/graph"
	"github.com/gonum/graph/encoding/dot"
	"github.com/Grissess/sdnd16/utils"
	// "github.com/Grissess/sdnd16/algorithms"
)

func main(){
	g, labels, err1 := utils.ReadFileToGraph("topology.txt")
	if err1 != nil {
		fmt.Println("ReadFileToGraph:");
		fmt.Println(err1);
		return;
	}
	labels = labels;
	// fmt.Println(g);
	bytes, err2 := dot.Marshal(g, "", "", "  ", false)
	if err2 != nil {
		fmt.Println("dot.Marshall:");
		fmt.Println(err1);
		return;
	}
	fmt.Println(string(bytes));
}
