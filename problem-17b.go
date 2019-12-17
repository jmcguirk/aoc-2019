package main

import (
	"strconv"
	"strings"
)

type Problem17B struct {

}

func (this *Problem17B) Solve() {
	Log.Info("Problem 17B solver beginning!")


	grid := &IntegerGrid2D{};
	grid.Init();

	robot := &RobotScaffold{};
	err := robot.Init("source-data/input-day-17b.txt", grid);
	if(err != nil){
		Log.FatalError(err);
	}

	err = robot.ParseGrid();
	if(err != nil){
		Log.FatalError(err);
	}
	Log.Info("Successfully parsed grid");


	workingGrid := robot.Grid.Clone();
	startingPos := robot.CurrentPos;
	visitGrid := workingGrid.Clone();

	xMin := visitGrid.MinRow();
	xMax := visitGrid.MaxRow();

	yMin := visitGrid.MinCol();
	yMax := visitGrid.MaxCol();

	for j := yMin; j<= yMax; j++{
		for i := xMin; i<= xMax; i++{
			val := workingGrid.GetValue(i, j);
			if(val == AsciiCodeScaffold) {
				visitGrid.SetValue(i,j,AsciiRobotUnvisited)
			}
			if(val == AsciiRobotUp){
				visitGrid.SetValue(i,j,AsciiRobotVisited);
			}
		}

	}


	//Log.Info("\n" + PrintScaffoldGrid(visitGrid));

	currentOrientation := OrientationNorth;

	instructions := make([]string, 0);


	curPos := startingPos.Clone();

	for{
		west := workingGrid.GetValue(curPos.X-1, curPos.Y) == AsciiCodeScaffold;
		east := workingGrid.GetValue(curPos.X+1, curPos.Y) == AsciiCodeScaffold;
		north := workingGrid.GetValue(curPos.X, curPos.Y-1) == AsciiCodeScaffold;
		south := workingGrid.GetValue(curPos.X, curPos.Y+1) == AsciiCodeScaffold;

		westVisited := visitGrid.GetValue(curPos.X-1, curPos.Y) == AsciiRobotVisited;
		eastVisited := visitGrid.GetValue(curPos.X+1, curPos.Y) == AsciiRobotVisited;
		northVisited := visitGrid.GetValue(curPos.X, curPos.Y-1) == AsciiRobotVisited;
		southVisited := visitGrid.GetValue(curPos.X, curPos.Y+1) == AsciiRobotVisited;

		westIntersection := IsScaffoldIntersection(curPos.X-1, curPos.Y, workingGrid);
		eastIntersection := IsScaffoldIntersection(curPos.X+1, curPos.Y, workingGrid);
		northIntersection := IsScaffoldIntersection(curPos.X, curPos.Y-1, workingGrid);
		SouthIntersection := IsScaffoldIntersection(curPos.X, curPos.Y+1, workingGrid);



			if(currentOrientation == OrientationNorth){
				if((!northVisited || northIntersection) && north){
					currentOrientation = OrientationNorth;
					instructions = append(instructions, "M");
					curPos.Y--;
					visitGrid.SetValue(curPos.X, curPos.Y, AsciiRobotVisited);
					continue;
				} else if(!eastVisited && east){
					currentOrientation = OrientationEast;
					instructions = append(instructions, "R");
					instructions = append(instructions, "M");
					curPos.X++;
					visitGrid.SetValue(curPos.X, curPos.Y, AsciiRobotVisited);
					continue;
				} else if(!westVisited && west){
					currentOrientation = OrientationWest;
					instructions = append(instructions, "L");
					instructions = append(instructions, "M");
					curPos.X--;
					visitGrid.SetValue(curPos.X, curPos.Y, AsciiRobotVisited);
					continue;
				} else{
					break;
				}
			}

			if(currentOrientation == OrientationEast) {
				if ((!eastVisited || eastIntersection) && east) {
					currentOrientation = OrientationEast;
					instructions = append(instructions, "M");
					curPos.X++;
					visitGrid.SetValue(curPos.X, curPos.Y, AsciiRobotVisited);
					continue;
				} else if (!northVisited && north) {
					currentOrientation = OrientationNorth;
					instructions = append(instructions, "L");
					instructions = append(instructions, "M");
					curPos.Y--;
					visitGrid.SetValue(curPos.X, curPos.Y, AsciiRobotVisited);
					continue;
				} else if (!southVisited && south) {
					currentOrientation = OrientationSouth;
					instructions = append(instructions, "R");
					instructions = append(instructions, "M");
					curPos.Y++;
					visitGrid.SetValue(curPos.X, curPos.Y, AsciiRobotVisited);
					continue;
				} else {
					break;
				}
			}

			if(currentOrientation == OrientationSouth) {
				if ((!southVisited || SouthIntersection) && south) {
					currentOrientation = OrientationSouth;
					instructions = append(instructions, "M");
					curPos.Y++;
					visitGrid.SetValue(curPos.X, curPos.Y, AsciiRobotVisited);
					continue;
				} else if (!eastVisited && east) {
					currentOrientation = OrientationEast;
					instructions = append(instructions, "L");
					instructions = append(instructions, "M");
					curPos.X++;
					visitGrid.SetValue(curPos.X, curPos.Y, AsciiRobotVisited);
					continue;
				} else if (!westVisited && west) {
					currentOrientation = OrientationWest;
					instructions = append(instructions, "R");
					instructions = append(instructions, "M");
					curPos.X--;
					visitGrid.SetValue(curPos.X, curPos.Y, AsciiRobotVisited);
					continue;
				} else {
					break;
				}
			}

			if(currentOrientation == OrientationWest) {
				if ((!westVisited || westIntersection) && west) {
					currentOrientation = OrientationWest;
					instructions = append(instructions, "M");
					curPos.X--;
					visitGrid.SetValue(curPos.X, curPos.Y, AsciiRobotVisited);
					continue;
				} else if (!northVisited && north) {
					currentOrientation = OrientationNorth;
					instructions = append(instructions, "R");
					instructions = append(instructions, "M");
					curPos.Y--;
					visitGrid.SetValue(curPos.X, curPos.Y, AsciiRobotVisited);
					continue;
				} else if (!southVisited && south) {
					currentOrientation = OrientationSouth;
					instructions = append(instructions, "L");
					instructions = append(instructions, "M");
					curPos.Y++;
					visitGrid.SetValue(curPos.X, curPos.Y, AsciiRobotVisited);
					continue;
				}
			}
		}


	//Log.Info("\n" + PrintScaffoldGrid(workingGrid));
	//tour := strings.Join(instructions, ",");


	moveCount := 0;
	compressedInstructions := make([]string, 0);
	for _, c := range instructions{
		if(c != "M"){
			if(moveCount > 0){
				compressedInstructions = append(compressedInstructions, strconv.Itoa(moveCount));
				moveCount = 0;
			}
			compressedInstructions = append(compressedInstructions, c);
		} else{
			moveCount++;
		}
	}
	if(moveCount > 0){
		compressedInstructions = append(compressedInstructions, strconv.Itoa(moveCount));
		moveCount = 0;
	}

	// Movement compressed instructions
	// R,6,L,10,R,8,R,8,R,12,L,8,L,8,R,6,L,10,R,8,R,8,R,12,L,8,L,8,L,10,R,6,R,6,L,8,R,6,L,10,R,8,R,8,R,12,L,8,L,8,L,10,R,6,R,6,L,8,R,6,L,10,R,8,L,10,R,6,R,6,L,8
	stream := strings.Join(compressedInstructions, ",");

	Log.Info("Completed tour - \n%s", stream);
	combos := AllSubstrings(stream, len(stream));

	uniqueCombos := make(map[string]int);
	allValidCombos := make([]string, 0);

	for _, combo := range combos{
		if(isValidSubstring(combo)){
			_, exists := uniqueCombos[combo];
			if(!exists){
				uniqueCombos[combo] = 1;
				allValidCombos = append(allValidCombos, combo);
			}
		}
	}

	Log.Info("Considering %d valid combinations ", len(allValidCombos));

	programSet := make([]string, 3);
	foundProgram := false;
	mainProgram := "";
	for i, A := range allValidCombos{
		for j, B := range allValidCombos {
			for k, C := range allValidCombos {
				if (j > i && k > j) {
					programSet[0] = A;
					programSet[1] = B;
					programSet[2] = C;
					foundProgram, mainProgram = checkCandidate(stream, programSet);
				}
				if (foundProgram) {
					break;
				}
			}
			if (foundProgram) {
				break;
			}
		}
		if(foundProgram){
			break;
		}
	}

	Log.Info("Executing program");
	Log.Info("Main - %s", mainProgram);
	Log.Info("Program A - %s ", programSet[0]);
	Log.Info("Program B - %s ", programSet[1]);
	Log.Info("Program C - %s ", programSet[2]);

	output := robot.Execute(mainProgram, programSet[0], programSet[1], programSet[2]);
	Log.Info("Finished executing - output is %d",output);
}

func checkCandidate(instructionStream string, programs []string) (bool, string){
	res := strings.Replace(instructionStream, programs[0], "A", -1);
	res = strings.Replace(res, programs[1], "B", -1);
	res = strings.Replace(res, programs[2], "C", -1);

	//Log.Info(res);
	if(len(res) > 20){
		return false, "";
	}
	split := strings.Split(res, ",");
	for _, v := range split {
		if(v != "A" && v != "B" && v != "C"){
			return false, "";
		}
	}

	return true, res;
}

func isValidSubstring(instructionStream string) bool{
	if(len(instructionStream) == 0){
		return false;
	}
	if(len(instructionStream) > 20){
		return false;
	}
	if(string(instructionStream[0]) == ","){
		return false;
	}
	if(len(instructionStream) > 1) {
		if(string(instructionStream[len(instructionStream) - 1]) == ",") {
			return false
		}
	}
	return true;
}