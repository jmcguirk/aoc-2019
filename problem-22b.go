package main

type Problem22B struct {

}

func (this *Problem22B) Solve() {
	Log.Info("Problem 22B solver beginning!")

	exploreValue := 10007;
	//exploreValue := 10;

	deck := &Deck{};
	deck.InitSlim(exploreValue);

	instructions, err := deck.ParseShuffleInstructionSet("source-data/input-day-22b.txt");
	if(err != nil){
		Log.FatalError(err);
	}

	//Log.Info("Applied %d instructions deck is now: %s", len(instructions), deck.Print());

	//cardOfInterest := 2019;
	//index should be 4684
	cardOfInterest := 4684;

	composite := &DeckOperationComposite{};
	composite.CompressInverse(instructions, deck);
	Log.Info("Finished compressing operations y = %dx+%d", composite.M, composite.B);
	deck.SetCardOfInterest(cardOfInterest);
	composite.ApplySlim(deck);

	Log.Info("Applied compressed instructions card of interest %d is at position %d", deck.CardOfInterest, deck.CardOfInterestIndex);
}

