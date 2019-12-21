package main

type Problem21A struct {

}

func (this *Problem21A) Solve() {
	Log.Info("Problem 21A solver beginning!")

	robot := &RobotSpringDroid{};
	err := robot.Init("source-data/input-day-21a.txt");
	if(err != nil){
		Log.FatalError(err);
	}

	err = robot.LoadProgram("source-data/input-day-21a-program01.txt");
	if(err != nil){
		Log.FatalError(err);
	}

	Log.Info("Successfully initialized droid");
	//Log.Info(robot.DescribeLoadedProgram());
	err = robot.Execute(false);
	if(err != nil){
		Log.FatalError(err);
	}
}
