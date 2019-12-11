package main

type Problem11B struct {

}

func (this *Problem11B) Solve() {
	Log.Info("Problem 11B solver beginning!")

	grid := &IntegerGrid2D{};
	grid.Init();

	grid.SetValue(0, 0, 1);

	robot := &RobotPainter{};
	err := robot.Init("source-data/input-day-11b.txt", "Raphael", &IntVec2{}, grid);
	if(err != nil){
		Log.FatalError(err);
	}
	robot.PrintState();

	for {
		if(robot.IsComplete()){
			Log.Info("Robot completed - work done %d", robot.WorkDone);
			Log.Info("\n" + grid.Print());
			break;
		}
		err := robot.Step();
		if(err != nil){
			Log.FatalError(err);
		}
	}
}
