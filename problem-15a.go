package main

type Problem15A struct {

}

func (this *Problem15A) Solve() {
	Log.Info("Problem 15A solver beginning!")

	grid := &IntegerGrid2D{};
	grid.Init();

	robot := &RobotRepair{};
	err := robot.Init("source-data/input-day-15a.txt", &IntVec2{}, grid);
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
	Log.Info("Found oxygen system at %s - best path is %d steps", robot.SystemLocation.ToString(), len(grid.ShortestPath(robot.SystemLocation, &IntVec2{},RepairStateWall)));
}
