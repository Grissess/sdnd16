package reader

import(
	"bufio"
	"fmt"
	"os"
	"github.com/Grissess/sdnd16/network"
)


func ReadFile(filename string) network.DsGraph{

	f, _ := os.Open(filename)

	scanner := bufio.NewScanner(f);
	scanner.Split(bufio.ScanWords);
	var result []string;

	for scanner.Scan(){
		x := scanner.Text()
		result = append(result, x);
	}

	var srcs, dests, weights []string
	var i int

	for i =0 ; i < len(result); i= i+3{
		srcs = append(srcs,result[i]);
	}
	for i =1; i < len(result); i= i+3{
		dests = append(dests, result[i]);
	}

	for i = 2; i< len(result); i = i+3{
		weights = append(weights, result[i]);
	}

	g := network.NewGraph();

	for i = 0; i<len(srcs); i= i +1{
		s := g.GetOrCreateNode(network.Label(fmt.Sprintf("%d", srcs[i])))
		d := g.GetOrCreateNode(network.Label(fmt.Sprintf("%d", dests[i])))
		e1, _ := g.NewEdge(s,d);
		e2, _ := g.NewEdge(d,s);
		e1.SetAttr("cost", weights[i]);
		e2.SetAttr("cost", weights[i]);
	}

        return g;
}
