package main

type Problem21B struct {

}

func (this *Problem21B) Solve() {
	Log.Info("Problem 21B solver beginning!")

	robot := &RobotSpringDroid{};
	err := robot.Init("source-data/input-day-21b.txt");
	if(err != nil){
		Log.FatalError(err);
	}

	err = robot.LoadProgram("source-data/input-day-21b-program01.txt");
	if(err != nil){
		Log.FatalError(err);
	}

	Log.Info("Successfully initialized droid");
	//Log.Info(robot.DescribeLoadedProgram());
	err = robot.Execute(true);
	if(err != nil){
		Log.FatalError(err);
	}
}
