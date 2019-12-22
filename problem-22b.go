package main

type Problem22B struct {

}

func (this *Problem22B) Solve() {
	Log.Info("Problem 22B solver beginning!")

	deckSize := 119315717514047;
	iterationCount := 101741582076661;

	deck := &Deck{};
	deck.InitSlim(deckSize);

	instructions, err := deck.ParseShuffleInstructionSet("source-data/input-day-22b.txt");
	if(err != nil){
		Log.FatalError(err);
	}

	//Log.Info("Applied %d instructions deck is now: %s", len(instructions), deck.Print());

	//cardOfInterest := 2019;
	//index should be 4684
	cardOfInterest := 2020;

	composite := &DeckOperationComposite{};
	composite.CompressInverse(instructions, deck);
	Log.Info("Finished compressing operations y = %dx+%d", composite.M, composite.B);
	deck.SetCardOfInterest(cardOfInterest);
	// Incorrect answers 90444333371210
	// Incorrect answers 70053774363234 - too high
	// Incorrect answers 64470772223464 - too high
	// 					 11061385851898 - Unknown
	//					 112201712869651 - Unknown
	//					 51481617374172 - Unknown
	composite.ApplyMulti(deck, iterationCount);
	//composite.ApplySlim(deck);

	Log.Info("Applied inverse compressed instructions. We are interested in slot %d which has value %d", cardOfInterest, deck.CardOfInterestIndex);
}

