package main

type Problem13A struct {

}

func (this *Problem13A) Solve() {
	Log.Info("Problem 13A solver beginning!")

	grid := &IntegerGrid2D{};
	grid.Init();

	game := &IntCodeVideoGame{};
	err := game.Init("source-data/input-day-13a.txt", grid);
	if(err != nil){
		Log.FatalError(err);
	}
	err = game.ParseMap()
	if(err != nil){
		Log.FatalError(err);
	}
	Log.Info("Finished parsing game - has %d tiles", game.TileCount);
	Log.Info("\n%s", game.Print())
	Log.Info("%d Blocks Remain", game.CountBlocks());
}
