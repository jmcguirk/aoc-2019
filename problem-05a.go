package main

type Problem5A struct {

}

func (this *Problem5A) Solve() {


	machine := &IntcodeMachineV2{};
	err := machine.Load("source-data/input-day-05a.txt")

	if err != nil {
		Log.FatalError(err);
	}


	machine.SetInputValue(1);
	machine.PrintContents();

	err = machine.Execute();
	if err != nil {
		Log.FatalError(err);
	}
	machine.PrintContents();

	Log.Info("Finished executing machine successfully - contents of first register are %d", machine.GetValueAtRegister(0));
}
