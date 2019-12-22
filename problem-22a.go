package main

type Problem22A struct {

}

func (this *Problem22A) Solve() {
	Log.Info("Problem 22A solver beginning!")

	deck := &Deck{};
	deck.Init(10007);

	instructions, err := deck.ParseShuffleInstructionSet("source-data/input-day-22a.txt");
	if(err != nil){
		Log.FatalError(err);
	}
	deck.Apply(instructions);
	Log.Info("Applied %d instructions deck is now: %s", len(instructions), deck.Print());

	cardOfInterest := 2019;
	Log.Info("Card %d is at position %d in deck", cardOfInterest, deck.IndexOf(cardOfInterest));
}

