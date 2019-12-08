package main



type Problem8B struct {

}

func (this *Problem8B) Solve() {
	Log.Info("Problem 8B solver beginning!")

	image := &IntImage{};

	//err := image.Parse(2, 2,"source-data/input-day-08b-test.txt");
	err := image.Parse(25, 6,"source-data/input-day-08b.txt");

	if err != nil {
		Log.FatalError(err);
	}


	res := image.FlattenAndDraw();

	Log.Info(res);

	Log.Info("Finished parsing image - total layers %d", len(image.Layers));
}
