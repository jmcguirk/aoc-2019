package main

type Problem23A struct {

}

func (this *Problem23A) Solve() {
	Log.Info("Problem 23A solver beginning!")

	network := &IntcodeNetwork{};
	network.Init();

	err := network.AddTerminals("source-data/input-day-23a.txt", 50);

	if(err != nil){
		Log.FatalError(err);
	}

	network.AddressOfInterest = 255;
	err = network.Simulate();
	if(err != nil){
		Log.FatalError(err);
	}
}

