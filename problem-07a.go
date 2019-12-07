package main

import (
	"log"
	"math"
	"strconv"
)

type Problem7A struct {

}

func (this *Problem7A) Solve() {


	machine := &IntcodeMachineV2{};

	err := machine.Load("source-data/input-day-07a.txt")

	if err != nil {
		Log.FatalError(err);
	}

	vals := make([]int64, 0);
	vals = append(vals, 0);
	vals = append(vals, 1);
	vals = append(vals, 2);
	vals = append(vals, 3);
	vals = append(vals, 4);

	permutations := make([][]int64, 0);

	greatestVal := int64(math.MinInt64);
	var greatestInput []int64;

	Perm(vals, func(a []int64) {
		cpy := make([]int64, len(a));
		copy(cpy, a);
		permutations = append(permutations, cpy);
	})

	for _, perm := range permutations{
		currentInput := int64(0);

		for _, v := range perm {
			machine.Reset();
			machine.QueueInput(v);
			machine.QueueInput(currentInput);
			err := machine.Execute();
			if(err != nil){
				log.Fatal("Encountered error");
			}
			currentInput = int64(machine.LastOutputValue);
		}
		//Log.Info("processed %s - %d", toString(perm), currentInput)

		if(currentInput > greatestVal){
			greatestVal = currentInput;
			greatestInput = perm;

		}
	}



	Log.Info("Finished executing machine successfully - greatest output is are %d using inputs %s", greatestVal, toString(greatestInput));
}

func toString(arr []int64) string{
	buff := "";
	for _, v := range arr{
		buff += strconv.FormatInt(v, 10);
	}
	return buff;
}