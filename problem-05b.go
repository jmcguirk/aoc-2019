package main

type Problem5B struct {

}

func (this *Problem5B) Solve() {


	machine := &IntcodeMachineV2{};
	err := machine.Load("source-data/input-day-05b.txt")

	if err != nil {
		Log.FatalError(err);
	}


	machine.SetInputValue(5);
	//machine.PrintContents();

	err = machine.Execute();
	if err != nil {
		Log.FatalError(err);
	}
	//machine.PrintContents();

	Log.Info("Finished executing machine successfully");
}
