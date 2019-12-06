package main

import (
	"bufio"
	"os"
	"strings"
)

type Problem6B struct {

}

func (this *Problem6B) Solve() {


	Log.Info("Problem 6B solver beginning!")


	file, err := os.Open("source-data/input-day-06b.txt");
	if err != nil {
		Log.FatalError(err);
	}
	defer file.Close()

	graph := &UndirectedGraph{};
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

	us := graph.GetOrCreateNode("YOU");
	santa := graph.GetOrCreateNode("SAN");

	path := us.ShortestPath(santa);
	if(path == nil){
		Log.Info("Failed to find a path to santa :(");
	} else{
		Log.Info("Shorted path to santa is %d transfers ", len(path) - 1 ); // Desired count is edges, not nodes
	}

	//Log.Info(us.Describe());
	//Log.Info(santa.Describe());

}
