package main

import (
	"bufio"
	"math"
	"os"
	"sort"
	"strings"
)

type Problem10B struct {

}

const rotationEpsilon = 0.0001;

func (this *Problem10B) Solve() {
	Log.Info("Problem 10B solver beginning!")


	file, err := os.Open("source-data/input-day-10b.txt");
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

	Log.Info("Analysis done - best point is at %d,%d with reachability of %d out of %d points", bestPoint.X, bestPoint.Y, bestPointReachable, len(points));

	tA := &IntVec2{};
	tA.X = 11;
	tA.Y = 13;
	tB := &IntVec2{};

	tB.X = 11;
	tB.Y = 12;

	Log.Info("angle between two points %.2f", tA.Angle(tB));

	remainingAsteroids := Filter(bestPoint, points);

	destroyIndex := 1;
	currOrientation := -float32(90 * (math.Pi / 180));
	targetStep := 200;
	Log.Info("Starting cycle with %d points at orientation %.2f", len(remainingAsteroids), currOrientation);
	for{
		if(len(remainingAsteroids) <= 0) {
			break;
		}
		nextTarget := remainingAsteroids[0];
		if(len(remainingAsteroids) > 1){
			reachable := bestPoint.GetVisiblePoints(remainingAsteroids);
			sort.SliceStable(reachable, func(i, j int) bool {
				angleA := bestPoint.Angle(reachable[i]);
				angleB := bestPoint.Angle(reachable[j]);
				behindA := angleA < currOrientation;
				behindB := angleB < currOrientation;
				if(behindA != behindB){
					if(behindA){
						return false;
					}
					return true;
				}
				return angleA - currOrientation < angleB - currOrientation;
				//return math.Abs(float64(angleA - currOrientation)) < math.Abs(float64(angleB - currOrientation));
			});

			nextTarget = reachable[0];
			currOrientation = bestPoint.Angle(nextTarget) + rotationEpsilon;
		}


		if(destroyIndex == targetStep){
			Log.Info("Step %d - Blowing up %d,%d - rotation is %.2f - Answer is %d", destroyIndex, nextTarget.X, nextTarget.Y, currOrientation, 100*nextTarget.X + nextTarget.Y);
		}

		destroyIndex++;
		remainingAsteroids = Filter(nextTarget, remainingAsteroids);
	}
	Log.Info("All asteroids obliterated!")
}
