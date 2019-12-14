package main

type Problem14B struct {

}



func (this *Problem14B) Solve() {
	Log.Info("Problem 14B solver beginning!")



	system := &ChemicalReactionSystem{};
	err := system.Load("source-data/input-day-14b.txt");
	if(err != nil){
		Log.FatalError(err);
	}

	reserves := 1000000000000;


	outputProduct := "FUEL";
	inputDesired := "ORE";

	// Pretty naive approach, but it's fast enough
	for i := 1; i < reserves; i++{
		res, err := system.GetTotalInputRequired(outputProduct, i, inputDesired);
		if(err != nil){
			Log.FatalError(err);
		}
		if(res >= reserves){
			Log.Info("Exhausted reserves - we can produce a max of %d fuel", i - 1 );
			break;
		} else{
			Log.Info("%d - %d", i, res);
		}
	}



	//Log.Info("Analysis done - we require %d %s to produce %d %s", res, inputDesired, outputQuantity, outputProduct);
}



