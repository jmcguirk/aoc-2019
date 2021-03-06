package main

import (
	"bufio"
	"math"
	"os"
	"strings"
)

type Problem3B struct {

}

func (this *Problem3B) Solve() {
	Log.Info("Problem 3B solver beginning!")


	file, err := os.Open("source-data/input-day-03a.txt");
	if err != nil {
		Log.FatalError(err);
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)




	lineId := 0;
	allLines := make(map[int]*TwistyLine)

	for scanner.Scan() {
		lineRaw := strings.TrimSpace(scanner.Text());
		if(lineRaw != ""){
			line := &TwistyLine{};
			line.Id = lineId;
			lineId++;
			err := line.Parse(lineRaw);
			if(err != nil){
				Log.FatalError(err);
			}
			allLines[line.Id] = line;
			//Log.Info("Added line with %d segments", len(line.LineSegments));
		}
	}

	grid := &IntegerGrid2D{};
	grid.Init();

	allIntersections := make([]*IntVec2, 0);



	for _, line := range allLines {
		intersects := line.Apply(grid);
		for _, point := range intersects {
			allIntersections = append(allIntersections, point)
		}
	}

	origin := &IntVec2{};
	origin.X = 0;
	origin.Y = 0;

	minDistance := math.MaxInt64;
	for _, intersection := range allIntersections {
		combinedDistance := 0;
		for _, line := range allLines {
			distance := line.StepsToIntersection(intersection);
			if(distance < 0){
				Log.Fatal("Steps to intersection returned a negative number - implies no intersection");
			}
			combinedDistance += distance;
		}
		if(combinedDistance < minDistance){
			minDistance = combinedDistance;
		}

	}
	Log.Info("Finished parsing file %d twisty lines - found %d intersections - fewested combined steps at %d", len(allLines), len(allIntersections), minDistance);
}

