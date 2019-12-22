package main

import "os"

type Problem21BExtra struct {

}

func (this *Problem21BExtra) Solve() {
	Log.Info("Problem 21B solver (bonus round) beginning!")

	robot := &RobotSpringDroid{};
	err := robot.Init("source-data/input-day-21b.txt");
	if(err != nil){
		Log.FatalError(err);
	}

	allInstrux := make([]SpringCodeInstruction, 0);
	allLetters := make([]int, 0);
	allLetters = append(allLetters, int('A'));
	allLetters = append(allLetters, int('B'));
	allLetters = append(allLetters, int('C'));
	allLetters = append(allLetters, int('D'));
	allLetters = append(allLetters, int('E'));
	allLetters = append(allLetters, int('F'));
	allLetters = append(allLetters, int('G'));
	allLetters = append(allLetters, int('H'));
	allLetters = append(allLetters, int('I'));
	allLetters = append(allLetters, int('J'));
	allLetters = append(allLetters, int('T'));

	for _, l := range allLetters{
		andJ := &SpringCodeAND{};
		andJ.Register1 = l;
		andJ.Register2 = int('J');

		allInstrux = append(allInstrux, andJ);

		andT := &SpringCodeAND{};
		andT.Register1 = l;
		andT.Register2 = int('T');

		allInstrux = append(allInstrux, andT);

		orJ := &SpringCodeOR{};
		orJ.Register1 = l;
		orJ.Register2 = int('J');

		allInstrux = append(allInstrux, orJ);

		orT := &SpringCodeOR{};
		orT.Register1 = l;
		orT.Register2 = int('T');

		allInstrux = append(allInstrux, orT);

		notJ := &SpringCodeNOT{};
		notJ.Register1 = l;
		notJ.Register2 = int('J');

		allInstrux = append(allInstrux, notJ);

		notT := &SpringCodeNOT{};
		notT.Register1 = l;
		notT.Register2 = int('T');

		allInstrux = append(allInstrux, notT);
	}

	maxOd := len(allInstrux) - 1;

	Log.Info("Total instructions %d", len(allInstrux));

	run := &SpringCodeRun{};
	maxLen := 15;

	for i := 1; i <= maxLen; i++{
		program := make([]SpringCodeInstruction, i+1);
		program[i] = run;
		indexArr := make([]int, i);
		for i, _ := range indexArr{
			indexArr[i] = 0;
		}

		for {
			atLim := false;
			for j := len(indexArr) - 1; j >= 0; j--{
				if(indexArr[j] + 1 < maxOd){
					indexArr[j]++;
					break;
				} else{
					if(j == 0){
						atLim = true;
						break;
					}
					indexArr[j] = 0;
				}
			}
			if(atLim){
				break;
			}
			for j, pos := range indexArr{
				program[j] = allInstrux[pos];
			}
			success, _ := robot.ExecuteRaw(program);
			if(success){
				Log.Info("Found a working solution with length %d", len(program));
				Log.Info("\n" + robot.DescribeLoadedProgram());
				os.Exit(0);
			}
		}

		Log.Info("No viable programs of len %d found", i);

	}
}
