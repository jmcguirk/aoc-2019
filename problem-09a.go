package main

type Problem9A struct {

}

func (this *Problem9A) Solve() {

	//legacyMachine := &IntcodeMachineV2{};
	//err := legacyMachine.Load("source-data/input-day-05a.txt")

	//if err != nil {
	//	Log.FatalError(err);
	//}
	//legacyMachine.QueueInput(1)
	//legacyMachine.PrintContents();
	//err = legacyMachine.Execute();
	//if err != nil {
		//Log.FatalError(err);
	//}


	machine := &IntcodeMachineV3{};
	err := machine.Load("source-data/input-day-09a.txt")

	if err != nil {
		Log.FatalError(err);
	}


	machine.QueueInput(1)

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
