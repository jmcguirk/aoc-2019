package main

import "math"

type Problem8A struct {

}

func (this *Problem8A) Solve() {
	Log.Info("Problem 8A solver beginning!")

	image := &IntImage{};

	//err := image.Parse(3, 2,"source-data/input-day-08a-test.txt");
	err := image.Parse(25, 6,"source-data/input-day-08a.txt");

	if err != nil {
		Log.FatalError(err);
	}


	var fewestZeroLayer *IntImageLayer;
	fewestedZeros := math.MaxInt64;
	for _, layer := range image.Layers {
		zeroCount := layer.GetHistValue(0);
		if(zeroCount < fewestedZeros){
			fewestZeroLayer = layer;
			fewestedZeros = zeroCount;
		}
	}

	check := fewestZeroLayer.GetHistValue(1) * fewestZeroLayer.GetHistValue(2);

	//Log.Info(image.Draw());

	Log.Info("Finished parsing image - total layers %d - output value %d", len(image.Layers), check);
}
