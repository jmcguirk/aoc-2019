package main

type Problem12B struct {

}



func (this *Problem12B) Solve() {
	Log.Info("Problem 12B solver beginning!")


	system := &PlanetarySystem{};
	err := system.Load("source-data/input-day-12b.txt");
	if(err != nil){
		Log.FatalError(err);
	}


	for{
		system.Step();
		if(system.HasRestedOnAllAxises()){
			Log.Info("All axis have been rested on frame %d", system.TimeIndex);
			Log.Info("Period is %d", int64(2)*LCM(int64(system.XPeriod), int64(system.YPeriod), int64(system.ZPeriod)));
			break;
		}
		if(system.TimeIndex % 100000000 == 0){
			Log.Info("Step %d", system.TimeIndex);
		}
	}

	Log.Info("Analysis done");
}



