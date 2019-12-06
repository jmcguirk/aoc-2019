package main

import (
	"math"
	"sort"
)

type Node struct {
	Id 		int;
	Label 	string;
	Edges  []*Edge;
	ContainingGraph Graph;
}



func (this *Node) Init(graph Graph) {
	this.ContainingGraph = graph;
	this.Edges = make([]*Edge, 0);
}


func (this *Node) ReachableNodes() map[int]*Node {

	res := make(map[int]*Node);
	visitedNodes := make(map[int]*Node);

	frontier := make([]*Node, 0);
	frontier = append(frontier, this);

	for {
		if(len(frontier) <= 0){
			break;
		}

		next := frontier[0];
		frontier = frontier[1:];

		for _, edge := range next.Edges{
			_, visited := visitedNodes[edge.To.Id];
			if(!visited) {
				res[edge.To.Id] = edge.To;
				frontier = append(frontier, edge.To);
				visitedNodes[edge.To.Id] = edge.To;
			}
		}

	}
	return res;
}

func (this *Node) ShortestPath(end *Node) []*Node {

	start := this;
	res := make([]*Node, 0);

	visitedNodes := make(map[int]*Node);
	minCostToStart := make(map[int]int);
	nearestToStart := make(map[int]*Node);

	frontier := make([]*Node, 0);
	frontier = append(frontier, start);
	frontierMap := make(map[int]*Node);
	frontierMap[start.Id] = start;
	minCostToStart[start.Id] = 0;

	for {
		if (len(frontier) <= 0) {
			break;
		}
		sort.SliceStable(frontier, func(i, j int) bool {
			return minCostToStart[i] < minCostToStart[j];
		});

		next := frontier[0];
		frontier = frontier[1:];
		delete(frontierMap, next.Id);
		costToHere := minCostToStart[next.Id];
		for _, edge := range next.Edges{
			neighbor := edge.To;
			_, visited := visitedNodes[neighbor.Id];
			if(visited){
				continue;
			}

			bestToHere, bestCostExists := minCostToStart[neighbor.Id];
			if(!bestCostExists){
				bestToHere = int(math.MaxInt32);
			}

			if(costToHere + edge.Weight < bestToHere){
				minCostToStart[neighbor.Id] = costToHere + edge.Weight;
				nearestToStart[neighbor.Id] = next;
				_, alreadyEnqueued := frontierMap[neighbor.Id];
				if(!alreadyEnqueued){
					frontierMap[neighbor.Id] = neighbor;
					frontier = append(frontier, neighbor);
				}
			}

		}
		visitedNodes[next.Id] = next;
		if(next == end){
			break;
		}
	}

	_, exists := minCostToStart[end.Id];
	if(!exists){
		return nil; // No path found
	}

	nextPathStep := end.Id;

	for {
		next := nearestToStart[nextPathStep];
		if(next == start){
			break;
		}
		nextPathStep = next.Id;
		res = append(res, next);
	}

	ReverseSlice(res);
	return res;
}



func (this *Node) Describe() string {
	buff := "";
	buff += this.Label;
	//buff += "(" + strconv.Itoa(this.Id) + ")";
	adjacencyList := ""
	for _, e := range this.Edges {
		if(adjacencyList == ""){
			adjacencyList += " - ";
		} else{
			adjacencyList += ",";
		}
		adjacencyList += e.To.Label;
	}
	buff += adjacencyList;
	return buff;
}

type Edge struct {
	Id 		int;
	Label 	string;
	From 	*Node;
	To 		*Node;
	Weight  int;
}

type Graph interface {
	Init();
	AllNodes() []*Node;
}

