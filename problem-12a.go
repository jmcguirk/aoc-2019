package main

type Problem12A struct {

}



func (this *Problem12A) Solve() {
	Log.Info("Problem 12A solver beginning!")


	system := &PlanetarySystem{};
	err := system.Load("source-data/input-day-12a.txt");
	if(err != nil){
		Log.FatalError(err);
	}
	targetSteps := 1000;

	for i:= 0; i < targetSteps; i++{
		system.Step();
		Log.Info(system.Print());
	}


	Log.Info("Analysis done");
}



