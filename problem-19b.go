package main

import (
	"bufio"
	"os"
	"strings"
)

type Problem19B struct {

}

func (this *Problem19B) Solve() {
	Log.Info("Problem 19B solver beginning!")


	/*
	grid := &IntegerGrid2D{};
	grid.Init();

	robot := &RobotTractorBeam{};
	err := robot.Init("source-data/input-day-19a.txt", grid);
	if(err != nil){
		Log.FatalError(err);
	}

	tractorCells, err := robot.ParseGrid(2500, 2500);
	if(err != nil){
		Log.FatalError(err);
	}

	Log.Info("Successfully parsed grid, found %d cells" , tractorCells);


	xMin := grid.MinRow();
	xMax := grid.MaxRow();

	yMin := grid.MinCol();
	yMax := grid.MaxCol();

	buff := "";
	for j := yMin; j<= yMax; j++{
		for i := xMin; i<= xMax; i++{
			if(!grid.HasValue(i, j)){
				buff += ".";
			} else{
				val := grid.GetValue(i, j);
				if(val > 0){
					buff += "#"
				} else{
					buff += ".";
				}
			}
		}
		buff += "\n";
	}

	Log.Info("Writing to file");
	//Log.Info("\n" + buff);



	data := []byte(buff)
	ioutil.WriteFile("source-data/input-day-19b-big.txt", data, 0644)
	Log.Info("Wrote to file!");*/



	grid := &IntegerGrid2D{};
	grid.Init();

	file, err := os.Open("source-data/input-day-19b-big.txt");
	if err != nil {
		Log.FatalError(err);
	}
	defer file.Close()

	row := 0;
	col := 0;
	cells := 0;
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text());
		if(line != ""){
			for _, c := range line{
				if int(c) == int('#'){
					grid.SetValue(col, row, 1);
					cells++;
				} else{
					grid.SetValue(col, row, 0);
				}
				col++;
			}

			row++;
			col = 0;
		}
	}

	Log.Info("Parsed reading contains %d cells", cells);

	targetSize := 100;
	maxY := grid.MaxRow();
	maxX := grid.MaxCol();

	for row := 0; row <= maxY; row++{

		tracking := 0;
		for x := 0; x <= maxX; x++ {
			if (grid.GetValue(x, row) == 1){
				tracking++;
			} else{
				if(tracking > 0){
					//Log.Info("row %d - width %d", row, tracking);
					tracking = 0;
				}
			}
		}
	}




	for x := 0; x <= maxX; x++{
		for y := 0; y <= maxY; y++{
			if(grid.GetValue(x, y) == 1){
				allSet := true;
				for i := 0; i < targetSize; i++{
					for j := 0; j < targetSize; j++{
						if(grid.GetValue(x+i, y+j) != 1){
							allSet = false;
							break;
						}
					}
					if(!allSet){
						break;
					}
				}
				if(allSet){
					Log.Info("Found a suitable window of size %d at %d, %d", targetSize, x, y);
					os.Exit(0);
				}

			}
		}
	}

	Log.Info("failed to find a window of size %d", targetSize);
	//Log.Info("\n" + robot.PrintGrid());*/
}
