package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type Problem20A struct {

}

const MazeSpaceChar = int(' ');
const MazeWallChar = int('#');
const MazePathChar = int('.');
const MazePortalCharStart = int('A');
const MazePortalCharEnd = int('Z');


type Portal struct{
	Label string;
	IsMazeStart bool;
	IsMazeExit bool;
	StartPosition *IntVec2;
	EndPosition *IntVec2;
	StartLabelPosition *IntVec2;
	EndLabelPosition *IntVec2;
	PortalStartIsDescending bool;
	PortalEndIsDescending bool;
}

func (this *Problem20A) Solve() {


	Log.Info("Problem 20A solver beginning!")


	file, err := os.Open("source-data/input-day-20a.txt");
	if err != nil {
		Log.FatalError(err);
	}
	defer file.Close()

	graph := &UndirectedGraph{};
	graph.Init();

	grid := &IntegerGrid2D{};
	grid.Init();

	scanner := bufio.NewScanner(file)
	x := 0;
	y := 0;
	for scanner.Scan() {
		line := scanner.Text();
		if (line != "") {
			//Log.Info(line)
			for _, c := range line{
				if(int(c) != MazeSpaceChar){
					grid.SetValue(x, y, int(c));
				}


				x++;
			}
		}
		x = 0;
		y++;
	}

	portals := make(map[string]*Portal)

	portalGrid := grid;//grid.Clone();

	maxX := portalGrid.MaxX();
	maxY := portalGrid.MaxY();


	//Log.Info("Max X: %d, MaxY: %d", maxX, maxY );

	var entrance *Portal;
	var exit *Portal;

	for y := 0; y <= maxY; y++{

		for x := 0; x <= maxX; x++ {
			val := portalGrid.GetValue(x, y);
			if(val >= MazePortalCharStart && val <= MazePortalCharEnd){
				// Vertical oriented
				sVal := portalGrid.GetValue(x, y+1);
				if(sVal >= MazePortalCharStart && sVal <= MazePortalCharEnd){
					label := fmt.Sprintf("%c%c", val, sVal);


					portal := &Portal{};
					portal.Label = label;

					anchor := IntVec2{};
					anchor.X = x;
					if(portalGrid.GetValue(x, y+2) == MazePathChar){
						anchor.Y = y+2;
					} else{
						anchor.Y = y-1;
					}

					if(label == "AA"){
						portal.IsMazeStart = true;
						portal.StartPosition = &anchor;
						entrance = portal;
					} else if(label == "ZZ"){
						portal.StartPosition = &anchor;
						portal.IsMazeExit = true;
						exit = portal;
					} else {
						existing, exists := portals[label];
						if(!exists){
							portal.StartPosition = &anchor;
							portals[label] = portal;
						} else{
							//Log.Info("Joined portal %s", label);
							existing.EndPosition = &anchor;
						}
					}
					continue;
				}
				// Horiz oriented
				eVal := portalGrid.GetValue(x+1, y);
				if(eVal >= MazePortalCharStart && eVal <= MazePortalCharEnd){
					label := fmt.Sprintf("%c%c", val, eVal);
					portal := &Portal{};
					portal.Label = label;

					anchor := IntVec2{};
					anchor.Y = y;
					if(portalGrid.GetValue(x+2, y) == MazePathChar){
						anchor.X = x+2;
					} else{
						anchor.X = x-1;
					}

					if(label == "AA"){
						portal.IsMazeStart = true;
						portal.StartPosition = &anchor;
						entrance = portal;
					} else if(label == "ZZ"){
						portal.IsMazeExit = true;
						portal.StartPosition = &anchor;
						exit = portal;
					} else {
						existing, exists := portals[label];
						if (!exists) {
							portal.StartPosition = &anchor;
							portals[label] = portal;
						} else {
							//Log.Info("Joined portal %s", label);
							existing.EndPosition = &anchor;
						}
					}
					continue;
				}
			}
		}
	}

	Log.Info("Finished parsing portals - Entrance at %d,%d Exit at %d,%d", entrance.StartPosition.X, entrance.StartPosition.Y, exit.StartPosition.X, exit.StartPosition.Y);

	for y := 0; y <= maxY; y++{ // Load the normal topology into the grid

		for x := 0; x <= maxX; x++ {
			val := portalGrid.GetValue(x, y);
			if(val == MazePathChar){
				node := graph.GetOrCreateNode(fmt.Sprintf("%d,%d", x, y));
				nVal := grid.GetValue(x, y-1);
				if(nVal == MazePathChar){
					neighbor := graph.GetOrCreateNode(fmt.Sprintf("%d,%d", x, y-1));
					graph.CreateEdge(node, neighbor);
				}
				sVal := grid.GetValue(x, y+1);
				if(sVal == MazePathChar){
					neighbor := graph.GetOrCreateNode(fmt.Sprintf("%d,%d", x, y+1));
					graph.CreateEdge(node, neighbor);
				}
				wVal := grid.GetValue(x-1, y);
				if(wVal == MazePathChar){
					neighbor := graph.GetOrCreateNode(fmt.Sprintf("%d,%d", x-1, y));
					graph.CreateEdge(node, neighbor);
				}
				eVal := grid.GetValue(x+1, y);
				if(eVal == MazePathChar){
					neighbor := graph.GetOrCreateNode(fmt.Sprintf("%d,%d", x+1, y));
					graph.CreateEdge(node, neighbor);
				}
			}
		}
	}

	Log.Info("Finished loading grid");


	for _, portal := range portals{
		if(portal.EndPosition == nil){
			log.Fatal("Portal " + portal.Label + " had no exit");
		}
		startAnchor := graph.GetOrCreateNode(fmt.Sprintf("%d,%d", portal.StartPosition.X, portal.StartPosition.Y));
		endAnchor := graph.GetOrCreateNode(fmt.Sprintf("%d,%d", portal.EndPosition.X, portal.EndPosition.Y));
		graph.CreateEdge(startAnchor, endAnchor);
		//Log.Info("Linked %s to %s using portal %s", startAnchor.Label, endAnchor.Label, portal.Label);
	}

	startNode := graph.GetOrCreateNode(fmt.Sprintf("%d,%d", entrance.StartPosition.X, entrance.StartPosition.Y));
	endNode := graph.GetOrCreateNode(fmt.Sprintf("%d,%d", exit.StartPosition.X, exit.StartPosition.Y));

	Log.Info("starting pathfinding from %s to %s", startNode.Label, endNode.Label);
	path := startNode.ShortestPath(endNode);
	if(path == nil){
		Log.Info("Failed to find a path!");
	} else{
		Log.Info("Shortest path from %s to %s is %d steps", startNode.Label, endNode.Label, len(path) +1);
	}
}
