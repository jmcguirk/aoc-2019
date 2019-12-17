package main

type Problem17A struct {

}

func (this *Problem17A) Solve() {
	Log.Info("Problem 17A solver beginning!")

	grid := &IntegerGrid2D{};
	grid.Init();

	robot := &RobotScaffold{};
	err := robot.Init("source-data/input-day-17a.txt", grid);
	if(err != nil){
		Log.FatalError(err);
	}

	err = robot.ParseGrid();
	if(err != nil){
		Log.FatalError(err);
	}

	Log.Info("Successfully parsed grid, checksum is %d", robot.CalculateCheckSum());
	Log.Info("\n" + robot.PrintGrid());
}
