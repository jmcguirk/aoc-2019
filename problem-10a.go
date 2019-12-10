package main

import (
	"bufio"
	"math"
	"os"
	"strings"
)

type Problem10A struct {

}



func (this *Problem10A) Solve() {
	Log.Info("Problem 10A solver beginning!")


	file, err := os.Open("source-data/input-day-10a.txt");
	if err != nil {
		Log.FatalError(err);
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)




	points := make([]*IntVec2, 0);

	j := 0;
	for scanner.Scan() {
		lineRaw := strings.TrimSpace(scanner.Text());
		if(lineRaw != ""){
			for i, letter := range lineRaw{
				if(string(letter) == "#"){
					p := &IntVec2{};
					p.X = i;
					p.Y = j;
					points = append(points, p);
				}
			}
			j++;
		}
	}
	Log.Info("Finished parsing file - considering %d points", len(points));

	bestPointReachable := int(math.MinInt32);
	var bestPoint *IntVec2;
	for _, candidate := range points {
		reachable := 0;
		for _, neighbor := range points {
			if(neighbor == candidate){
				continue;
			}
			isOccluded := false;
			slopeN := candidate.Slope(neighbor);
			distN := candidate.Distance(neighbor);
			for _, occluder := range points {
				if(occluder == neighbor || occluder == candidate){
					continue;
				}
				slopeO := candidate.Slope(occluder);


				if(math.Abs(float64(slopeN - slopeO)) <= slopeEpsilon){

					if(math.Abs(float64((candidate.Distance(occluder) + neighbor.Distance(occluder)) - distN)) <= distEpsilon){
						isOccluded = true;
					}
				}
			}
			if(!isOccluded){
				reachable++;
			}
		}
		if(reachable > bestPointReachable){
			bestPointReachable = reachable;
			bestPoint = candidate;
		}
	}

	Log.Info("Analysis done - best point is at %d,%d with reachability of %d", bestPoint.X, bestPoint.Y, bestPointReachable)
}
