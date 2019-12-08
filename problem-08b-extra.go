package main

import "github.com/otiai10/gosseract"

type Problem8BExtra struct {

}

func (this *Problem8BExtra) Solve() {
	Log.Info("Problem 8B (Bonus stuff) solver beginning!")

	image := &IntImage{};

	//err := image.Parse(2, 2,"source-data/input-day-08b-test.txt");
	err := image.Parse(25, 6,"source-data/input-day-08b.txt");

	if err != nil {
		Log.FatalError(err);
	}
	fileName := "artifacts/problem-08.png";
	image.FlattenAndRenderToFile(fileName, 512);

	client := gosseract.NewClient()
	client.SetPageSegMode(gosseract.PSM_SINGLE_WORD);
	client.Variables["tessedit_char_whitelist"] = UpperAlphaCharacters();
	defer client.Close()
	client.SetImage(fileName);
	text, err := client.Text();
	if(err != nil){
		Log.FatalError(err);
	}
	Log.Info("Tesseract found image contents - %s", text);
	//Log.Info("Finished parsing image - total layers %d", len(image.Layers));
}
