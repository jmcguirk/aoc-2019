package main

type Problem18A struct {

}

func (this *Problem18A) Solve() {
	Log.Info("Problem 18A solver beginning!")

	system := &TunnelSearchSystem{};
	err := system.Init("source-data/input-day-18a.txt");
	if(err != nil){
		Log.FatalError(err);
	}
	Log.Info("Parsed grid system - beginning path finding");
	system.FindShortestPath();

}
