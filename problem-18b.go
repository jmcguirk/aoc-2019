package main

type Problem18B struct {

}

func (this *Problem18B) Solve() {
	Log.Info("Problem 18B solver beginning!")

	system := &TunnelSearchSystemMultiAgent{};
	err := system.Init("source-data/input-day-18b.txt");
	if(err != nil){
		Log.FatalError(err);
	}
	//Log.Info("Parsed grid system - beginning path finding");
	//system.FindShortestPath();

}
