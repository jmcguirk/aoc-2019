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
	outputProduct := "FUEL";
	inputDesired := "ORE";
	reserves := 1000000000000;

	search := 1;
	stepSize := 2;
	for{
		res, err := system.GetTotalInputRequired(outputProduct, search, inputDesired);
		if(err != nil){
			Log.FatalError(err);
		}
		if(search < 0){
			break;
		}
		if(res < reserves){
			if(stepSize == 1){
				Log.Info("Exhausted reserves - max we can produce is %d", search);
				break;
			}
			stepSize *= 2;
			search += stepSize;
		} else {
			stepSize = stepSize / 2;
			search -= stepSize;
		}
	}

	/*

	for i := 1; i < reserves; i++{
		res, err := system.GetTotalInputRequired(outputProduct, i, inputDesired);
		if(err != nil){
			Log.FatalError(err);
		}
		if(res >= reserves){
			Log.Info("Exhausted reserves - we can product max is %d", i - 1 );
			break;
		} else{
			Log.Info("%d - %d", i, res);
		}
	}*/



	//Log.Info("Analysis done - we require %d %s to produce %d %s", res, inputDesired, outputQuantity, outputProduct);
}



