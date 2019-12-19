package main

type Problem19A struct {

}

func (this *Problem19A) Solve() {
	Log.Info("Problem 19A solver beginning!")

	grid := &IntegerGrid2D{};
	grid.Init();

	robot := &RobotTractorBeam{};
	err := robot.Init("source-data/input-day-19a.txt", grid);
	if(err != nil){
		Log.FatalError(err);
	}

	tractorCells, err := robot.ParseGrid(50, 50);
	if(err != nil){
		Log.FatalError(err);
	}

	Log.Info("Successfully parsed grid, found %d cells" , tractorCells);
	//Log.Info("\n" + robot.PrintGrid());
}
