package main

import (
	"bufio"
	"os"
)

type Problem24B struct {

}

const tileBug = int('#');
func (this *Problem24B) Solve() {


	Log.Info("Problem 24B solver beginning!")


	file, err := os.Open("source-data/input-day-24b.txt");
	if err != nil {
		Log.FatalError(err);
	}
	defer file.Close()



	tileEmpty := int('.');

	grid := &LayeredGrid{};
	grid.Init(tileEmpty);

	scanner := bufio.NewScanner(file)
	x := 0;
	y := 0;
	for scanner.Scan() {
		line := scanner.Text();
		if (line != "") {
			//Log.Info(line)
			for _, c := range line{
				grid.SetValue(x, y, 0, int(c));
				x++;
			}
		}
		x = 0;
		y++;
	}

	//Log.Info("\nInitial State\n%s", grid.PrintAscii());



	prevGen := grid.Clone();

	xMin := grid.GetOrCreateLayer(0).MinRow();
	xMax := grid.GetOrCreateLayer(0).MaxRow();

	yMin := grid.GetOrCreateLayer(0).MinCol();
	yMax := grid.GetOrCreateLayer(0).MaxCol();


	//uniqueScores := make(map[int]int);

	maxSteps := 200;
	for i := 1; i <= maxSteps; i++{
		grid.GetOrCreateLayer(i).SetValue(2, 2, int('?')); // Force expand
		grid.GetOrCreateLayer(i * -1).SetValue(2, 2, int('?'));
	}

	currStep := 0;
	for{


		if(currStep >= maxSteps){
			break;
		}

		minLayer := grid.MinLayer();
		maxLayer := grid.MaxLayer();
		prevGen = grid.Clone();



		for z := minLayer; z <= maxLayer; z++{
			for j := yMin; j<= yMax; j++{
				for i := xMin; i<= xMax; i++{
					if(i == 2 && j == 2) {
						continue;
					}
					if(!prevGen.HasValue(i, j, z)){
						continue;
					} else{
						val := prevGen.GetValue(i, j, z);
						tBug := prevGen.GetNeighborBugCount(i, j, z);
						if(val == tileBug){
							if(tBug != 1){
								grid.SetValue(i, j, z, tileEmpty);
							}
						} else if(val == tileEmpty){
							if(tBug == 1 || tBug == 2){
								grid.SetValue(i, j, z, tileBug);
							}
						}

					}
				}
			}
		}
		currStep++;
		Log.Info("%d bugs after %d minutes - min %d - max - %d", grid.CountValue(tileBug), currStep, minLayer, maxLayer);
	}
	minLayer := grid.MinLayer();
	maxLayer := grid.MaxLayer();
	total := 0;
	for z := minLayer; z <= maxLayer; z++{
		for j := yMin; j<= yMax; j++{
			for i := xMin; i<= xMax; i++{
				if(grid.GetValue(i, j, z) == tileBug){
					total++;
				}


			}
		}
	}

	Log.Info("Simulation complete! %d", total);
}

func (this *LayeredGrid) GetNeighborBugCount(x int, y int, z int) int {
	sum := 0;
	sum += this.CountIfBug(x+1, y, z);
	sum += this.CountIfBug(x-1, y, z);
	sum += this.CountIfBug(x, y+1, z);
	sum += this.CountIfBug(x, y-1, z);
	// Cell 12 in Drawing
	if(x == 0){
		sum += this.CountIfBug(1, 2, z-1);
	}
	// Cell 8 in Drawing
	if(y == 0){
		sum += this.CountIfBug(2, 1, z-1);
	}
	// Cell 14 in drawing
	if(x == 4){
		sum += this.CountIfBug(3, 2, z-1);
	}
	// Cell 18 in drawing
	if(y == 4){
		sum += this.CountIfBug(2, 3, z-1);
	}


	if(x == 1 && y == 2){
		for i := 0; i < 5; i++{
			sum += this.CountIfBug(0, i, z+1); // Roll up left side of board
		}
	}
	if(x == 3 && y == 2){
		for i := 0; i < 5; i++{
			sum += this.CountIfBug(4, i, z+1); // Roll up right side of board
		}
	}
	if(x == 2 && y == 1){
		for i := 0; i < 5; i++{
			sum += this.CountIfBug(i, 0, z+1); // Roll up top side of board
		}
	}
	if(x == 2 && y == 3){
		for i := 0; i < 5; i++{
			sum += this.CountIfBug(i, 4, z+1); // Roll up bottom side of board
		}
	}

	return sum;
}

func (this *LayeredGrid) CountIfBug(x int, y int, z int) int {
	if(this.GetValue(x, y, z) == tileBug){
		return 1;
	}
	return 0;
}