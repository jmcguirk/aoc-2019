package main

type Problem14A struct {

}



func (this *Problem14A) Solve() {
	Log.Info("Problem 14A solver beginning!")


	system := &ChemicalReactionSystem{};
	err := system.Load("source-data/input-day-14a.txt");
	if(err != nil){
		Log.FatalError(err);
	}

	outputProduct := "FUEL";
	outputQuantity := 1;
	inputDesired := "ORE";

	res, err := system.GetTotalInputRequired(outputProduct, outputQuantity, inputDesired);
	if(err != nil){
		Log.FatalError(err);
	}

	Log.Info("Analysis done - we require %d %s to produce %d %s", res, inputDesired, outputQuantity, outputProduct);
}



