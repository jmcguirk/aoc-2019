package main

type Problem22B struct {

}

func (this *Problem22B) Solve() {
	Log.Info("Problem 22B solver beginning!")

	exploreValue := 79357;
	//exploreValue := 10;

	deck := &Deck{};
	deck.InitSlim(exploreValue);

	instructions, err := deck.ParseShuffleInstructionSet("source-data/input-day-22b.txt");
	if(err != nil){
		Log.FatalError(err);
	}

	//Log.Info("Applied %d instructions deck is now: %s", len(instructions), deck.Print());

	cardOfInterest := 2020;
	deck.SetCardOfInterest(cardOfInterest);
	cycles := make(map[int]int);


	i:=0;
	for {
		deck.ApplySlim(instructions);
		i++;
		pos := deck.CardOfInterestIndex;
		orig, exists := cycles[pos];
		if(!exists){
			cycles[pos] = i;
		} else{
			delta := i - orig;
			cycles[pos] = i;
			Log.Info("Cycle detected %d - %d", delta, pos);
			break;
		}
	}

	Log.Info("Applied %d instructions card of interest %d is at position %d", len(instructions), deck.CardOfInterest, deck.CardOfInterestIndex);
}

