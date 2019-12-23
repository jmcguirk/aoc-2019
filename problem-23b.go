package main

type Problem23B struct {

}

func (this *Problem23B) Solve() {
	Log.Info("Problem 23B solver beginning!")

	network := &IntcodeNetwork{};
	network.Init();
	network.EnableNAT = true;

	err := network.AddTerminals("source-data/input-day-23a.txt", 50);

	if(err != nil){
		Log.FatalError(err);
	}
	err = network.Simulate();
	if(err != nil){
		Log.FatalError(err);
	}
}

