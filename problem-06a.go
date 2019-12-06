package main

import (
	"bufio"
	"os"
	"strings"
)

type Problem6A struct {

}

func (this *Problem6A) Solve() {


	Log.Info("Problem 6A solver beginning!")


	file, err := os.Open("source-data/input-day-06a.txt");
	if err != nil {
		Log.FatalError(err);
	}
	defer file.Close()

	graph := &DirectedGraph{};
	graph.Init();

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text());
		if(line != ""){
			lineParts := strings.Split(line, ")");
			labelB := lineParts[0];
			labelA := lineParts[1];
			nodeA := graph.GetOrCreateNode(labelA);
			nodeB := graph.GetOrCreateNode(labelB);
			graph.CreateEdge(nodeA, nodeB);

		}
	}


	totalOrbits := 0;

	for _, node := range graph.Nodes{
		totalOrbits +=  len(node.ReachableNodes());
	}

	Log.Info("Finished parsing file. total orbits is %d", totalOrbits);
}
