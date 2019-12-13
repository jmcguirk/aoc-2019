package main

type Problem13B struct {

}

func (this *Problem13B) Solve() {
	Log.Info("Problem 13B solver beginning!")

	grid := &IntegerGrid2D{};
	grid.Init();

	game := &IntCodeVideoGame{};
	err := game.Init("source-data/input-day-13b.txt", grid);
	if(err != nil){
		Log.FatalError(err);
	}
	game.EnableFreePlay();
	err = game.Play()
	if(err != nil){
		Log.FatalError(err);
	}
	Log.Info("Finished playing game - has %d blocks remain", game.CountBlocks());
}
