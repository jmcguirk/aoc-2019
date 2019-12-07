package main

import (
	"math"
)

type Problem7B struct {

}

func (this *Problem7B) Solve() {






	vals := make([]int64, 0);
	vals = append(vals, 5);
	vals = append(vals, 6);
	vals = append(vals, 7);
	vals = append(vals, 8);
	vals = append(vals, 9);

	permutations := make([][]int64, 0);

	greatestVal := int64(math.MinInt64);
	var greatestInput []int64;

	Perm(vals, func(a []int64) {
		cpy := make([]int64, len(a));
		copy(cpy, a);
		permutations = append(permutations, cpy);
	})

	machines := make([]*IntcodeMachineV2, 0);
	for i := 0; i < 5; i++{
		machine := &IntcodeMachineV2{};
		err := machine.Load("source-data/input-day-07b.txt")

		if err != nil {
			Log.FatalError(err);
		}
		machine.PauseOnOutput = true;
		machines = append(machines, machine);
	}



	for _, perm := range permutations{
		currentInput := int64(0);
		for _, machine := range machines{
			machine.Reset();
		}
		for i, v := range perm{
			machines[i].QueueInput(v);
		}

		//machines[0].QueueInput(currentInput);
		for {
			for _, machine := range machines{
				machine.QueueInput(currentInput);
				err := machine.Execute();
				if(err != nil){
					Log.Info(err.Error());
				}
				currentInput = machine.LastOutputValue;
			}
			if(machines[4].HasHalted){ // See if our last machine has halted and break if so
				break;
			}

		}
		Log.Info("processed %s - %d", toString(perm),  machines[4].LastOutputValue)
		if(currentInput > greatestVal){
			greatestVal = currentInput;
			greatestInput = perm;
		}
	}



	Log.Info("Finished executing machine successfully - greatest output is are %d using inputs %s", greatestVal, toString(greatestInput));
}