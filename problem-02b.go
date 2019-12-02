package main

type Problem2B struct {

}

func (this *Problem2B) Solve() {


	machine := &IntcodeMachineMutable{};
	err := machine.Load("source-data/input-day-02a.txt")

	target := int64(19690720);

	if err != nil {
		Log.FatalError(err);
	}

	for noun := 0; noun < 100000; noun++ {
		for verb := 0; verb < 100000; verb++ {
			machine.Reset();
			machine.SetValueAtRegister(1, int64(noun));
			machine.SetValueAtRegister(2, int64(verb));
			err := machine.Execute();
			if(err != nil){
				continue;
			}
			generated := machine.GetValueAtRegister(0);
			if(generated == target){
				Log.Info("Found magic pair of inputs at Noun: %d, Verb: %d - Solution : %d", noun, verb, 100 * noun + verb);
				return;
				break;
			}
		}
	}
	Log.Info("failed to find a pair of inputs");
}
