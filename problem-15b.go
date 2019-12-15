package main

type Problem15B struct {

}

func (this *Problem15B) Solve() {
	Log.Info("Problem 15B solver beginning!")

	grid := &IntegerGrid2D{};
	grid.Init();

	robot := &RobotRepair{};
	err := robot.Init("source-data/input-day-15b.txt", &IntVec2{}, grid);
	if(err != nil){
		Log.FatalError(err);
	}
	err = robot.Explore();
	if(err != nil){
		Log.FatalError(err);
	}
	Log.Info("Exploration finished");
	if(robot.SystemLocation == nil){
		Log.Info("Failed to find oxygen system");
	}
	Log.Info("Found oxygen system at %s beginning flood", robot.SystemLocation.ToString());
	steps := robot.SimulateOxygenFill();
	Log.Info("Found completed in %d steps", steps);
}
