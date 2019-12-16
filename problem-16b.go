package main

type Problem16B struct {

}

func (this *Problem16B) Solve() {
	Log.Info("Problem 16B solver beginning!")

	system := &FrequencySystem{};
	err := system.Parse("source-data/input-day-16b.txt");
	if(err != nil){
		Log.FatalError(err);
	}

	system.Log = true;
	system.StepMulti(10000, 100);
	Log.Info(system.ToShortString());
}
