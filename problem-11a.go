package main

type Problem11A struct {

}

func (this *Problem11A) Solve() {
	Log.Info("Problem 11A solver beginning!")

	grid := &IntegerGrid2D{};
	grid.Init();

	robot := &RobotPainter{};
	err := robot.Init("source-data/input-day-11a.txt", "Raphael", &IntVec2{}, grid);
	if(err != nil){
		Log.FatalError(err);
	}
	robot.PrintState();

	for {
		if(robot.IsComplete()){
			Log.Info("Robot completed - work done %d", robot.WorkDone);
			break;
		}
		err := robot.Step();
		if(err != nil){
			Log.FatalError(err);
		}
	}
}
