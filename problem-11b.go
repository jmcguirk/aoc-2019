package main

import "github.com/otiai10/gosseract"

type Problem11B struct {

}

func (this *Problem11B) Solve() {
	Log.Info("Problem 11B solver beginning!")

	grid := &IntegerGrid2D{};
	grid.Init();

	grid.SetValue(0, 0, 1);

	robot := &RobotPainter{};
	err := robot.Init("source-data/input-day-11b.txt", "Raphael", &IntVec2{}, grid);
	if(err != nil){
		Log.FatalError(err);
	}
	robot.PrintState();

	for {
		if(robot.IsComplete()){
			Log.Info("Robot completed - work done %d", robot.WorkDone);
			fileName := "artifacts/problem-11b.png"
			grid.SetValue(0, 0, 0); // Clean up origin
			grid.PrintToFile(fileName, 512);
			//Log.Info("\n" + grid.Print());
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
			break;
		}
		err := robot.Step();
		if(err != nil){
			Log.FatalError(err);
		}
	}
}
