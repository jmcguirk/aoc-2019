package main

import (
	"bufio"
	"fmt"
	"os"
)

type Problem20B struct {

}

func (this *Problem20B) Solve() {


	Log.Info("Problem 20A solver beginning!")


	file, err := os.Open("source-data/input-day-20b.txt");
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
				grid.SetValue(x, y, int(c));
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

	var interiorNW *IntVec2;
	var interiorNE *IntVec2;
	var interiorSW *IntVec2;
	var interiorSE *IntVec2;


	for y := 0; y <= maxY; y++{
		hasEncounteredWall := false;
		for x := 0; x <= maxX; x++ {
			val := portalGrid.GetValue(x, y);
			if(hasEncounteredWall && val == MazeSpaceChar){
				n := portalGrid.HasValue(x, y-1) && portalGrid.GetValue(x, y-1) == MazeWallChar;
				s := portalGrid.HasValue(x, y+1) && portalGrid.GetValue(x, y+1) == MazeWallChar;
				w := portalGrid.HasValue(x-1, y) && portalGrid.GetValue(x-1, y) == MazeWallChar;
				e := portalGrid.HasValue(x+1, y) && portalGrid.GetValue(x+1, y) == MazeWallChar;
				if(interiorNW == nil && n && w){
					interiorNW = &IntVec2{};
					interiorNW.X = x;
					interiorNW.Y = y;
				}
				if(interiorNE == nil && n && e){
					interiorNE = &IntVec2{};
					interiorNE.X = x;
					interiorNE.Y = y;
				}
				if(interiorSE == nil && s && e){
					interiorSE = &IntVec2{};
					interiorSE.X = x;
					interiorSE.Y = y;
				}
				if(interiorSW == nil && s && w){
					interiorSW = &IntVec2{};
					interiorSW.X = x;
					interiorSW.Y = y;
				}
			}
			if(val == MazeWallChar){
				hasEncounteredWall = true;
			}

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
							portal.StartLabelPosition = &IntVec2{};
							portal.StartLabelPosition.X = x;
							portal.StartLabelPosition.Y = y;
							portal.StartPosition = &anchor;
							portals[label] = portal;
						} else{
							existing.EndLabelPosition = &IntVec2{};
							existing.EndLabelPosition.X = x;
							existing.EndLabelPosition.Y = y;
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
							portal.StartLabelPosition = &IntVec2{};
							portal.StartLabelPosition.X = x;
							portal.StartLabelPosition.Y = y;
							portal.StartPosition = &anchor;
							portals[label] = portal;
						} else {
							existing.EndLabelPosition = &IntVec2{};
							existing.EndLabelPosition.X = x;
							existing.EndLabelPosition.Y = y;
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



	for _, portal := range portals{
		portal.PortalStartIsDescending = IsInterior(portal.StartLabelPosition,  interiorNW, interiorNE, interiorSW, interiorSE);
		portal.PortalEndIsDescending = IsInterior(portal.EndLabelPosition,  interiorNW, interiorNE, interiorSW, interiorSE);
		//Log.Info("%s start is descending %t, end is descending %t", portal.Label, portal.PortalStartIsDescending, portal.PortalEndIsDescending);
	}



	maxDepth := 100;

	for i := 0; i <= maxDepth; i++{
		for y := 0; y <= maxY; y++{ // Load the normal topology into the grid
			for x := 0; x <= maxX; x++ {
				if(i > 0 && x == entrance.StartPosition.X && y == entrance.StartPosition.Y){
					continue;
				}
				if(i > 0 && x == exit.StartPosition.X && y == exit.StartPosition.Y){
					continue;
				}
				val := portalGrid.GetValue(x, y);
				if(val == MazePathChar){
					node := graph.GetOrCreateNode(fmt.Sprintf("%d,%d,%d", x, y, i));
					nVal := grid.GetValue(x, y-1);
					if(nVal == MazePathChar){
						neighbor := graph.GetOrCreateNode(fmt.Sprintf("%d,%d,%d", x, y-1, i));
						graph.CreateEdge(node, neighbor);
					}
					sVal := grid.GetValue(x, y+1);
					if(sVal == MazePathChar){
						neighbor := graph.GetOrCreateNode(fmt.Sprintf("%d,%d,%d", x, y+1, i));
						graph.CreateEdge(node, neighbor);
					}
					wVal := grid.GetValue(x-1, y);
					if(wVal == MazePathChar){
						neighbor := graph.GetOrCreateNode(fmt.Sprintf("%d,%d,%d", x-1, y, i));
						graph.CreateEdge(node, neighbor);
					}
					eVal := grid.GetValue(x+1, y);
					if(eVal == MazePathChar){
						neighbor := graph.GetOrCreateNode(fmt.Sprintf("%d,%d,%d", x+1, y, i));
						graph.CreateEdge(node, neighbor);
					}
				}
			}
		}
	}

	for i := 0; i <= maxDepth; i++{
		for _, portal := range portals{
			if(!portal.PortalStartIsDescending && i == 0){
				continue;
			}

			nextDepth := i-1;
			if(portal.PortalStartIsDescending){
				nextDepth = i+1;
			}

			startAnchor := graph.GetOrCreateNode(fmt.Sprintf("%d,%d,%d", portal.StartPosition.X, portal.StartPosition.Y, i));
			endAnchor := graph.GetOrCreateNode(fmt.Sprintf("%d,%d,%d", portal.EndPosition.X, portal.EndPosition.Y, nextDepth));
			graph.CreateEdge(startAnchor, endAnchor);
			//Log.Info("Linked %s to %s using portal %s", startAnchor.Label, endAnchor.Label, portal.Label);
		}
	}


	Log.Info("Finished loading grid");

	startNode := graph.GetOrCreateNode(fmt.Sprintf("%d,%d,%d", entrance.StartPosition.X, entrance.StartPosition.Y, 0));
	endNode := graph.GetOrCreateNode(fmt.Sprintf("%d,%d,%d", exit.StartPosition.X, exit.StartPosition.Y, 0));

	Log.Info("starting pathfinding from %s to %s", startNode.Label, endNode.Label);
	path := startNode.ShortestPath(endNode);
	if(path == nil){
		Log.Info("Failed to find a path!");
	} else{
		Log.Info("Shortest path from %s to %s is %d steps", startNode.Label, endNode.Label, len(path) +1);
	}
}

func IsInterior(point *IntVec2, NW *IntVec2, NE *IntVec2, SW *IntVec2, SE *IntVec2) bool {
	if(point.X < NW.X || point.Y < NW.Y){
		return false;
	}
	if(point.X > NE.X || point.Y < NE.Y){
		return false;
	}
	if(point.X < SW.X || point.Y > SW.Y){
		return false;
	}
	if(point.X > SE.X || point.Y > SE.Y){
		return false;
	}
	return true;
}