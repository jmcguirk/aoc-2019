package main

type Problem9B struct {

}

func (this *Problem9B) Solve() {

	Log.Info("Problem 9B starting");

	machine := &IntcodeMachineV3{};
	err := machine.Load("source-data/input-day-09a.txt")

	if err != nil {
		Log.FatalError(err);
	}


	machine.QueueInput(2)

	err = machine.Execute();
	if err != nil {
		Log.FatalError(err);
	}
	//machine.PrintContents();
	lastOutput := "No Output";
	if(machine.LastOutputValue != nil){
		lastOutput = machine.LastOutputValue.String();
	}

	Log.Info("Finished executing machine successfully - output is %s", lastOutput);
}
