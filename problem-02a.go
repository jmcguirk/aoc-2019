package main

type Problem2A struct {

}

func (this *Problem2A) Solve() {


	machine := &IntcodeMachineMutable{};
	err := machine.Load("source-data/input-day-02a.txt")

	if err != nil {
		Log.FatalError(err);
	}
	// Per problem formulation - set the first two registers
	machine.SetValueAtRegister(1, 12);
	machine.SetValueAtRegister(2, 2);

	Log.Info("Parsed machine successfully, initial state is");
	machine.PrintContents();

	err = machine.Execute();
	if err != nil {
		Log.FatalError(err);
	}
	machine.PrintContents();

	Log.Info("Finished executing machine successfully - contents of first register are %d", machine.GetValueAtRegister(0));
}
