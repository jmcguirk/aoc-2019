package main

type Problem16A struct {

}

func (this *Problem16A) Solve() {
	Log.Info("Problem 16A solver beginning!")

	system := &FrequencySystem{};
	err := system.Parse("source-data/input-day-16a.txt");
	if(err != nil){
		Log.FatalError(err);
	}

	system.Step(100);
	Log.Info(system.ToShortString());
}
