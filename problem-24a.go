package main

import (
	"bufio"
	"os"
)

type Problem24A struct {

}

func (this *Problem24A) Solve() {


	Log.Info("Problem 24A solver beginning!")


	file, err := os.Open("source-data/input-day-24a.txt");
	if err != nil {
		Log.FatalError(err);
	}
	defer file.Close()

	graph := &UndirectedGraph{};
	graph.Init();

	grid := &IntegerGrid2D{};
	grid.Init();

	scanner := bufio.NewScanner(file)
	x := 0;
	y := 0;
	for scanner.Scan() {
		line := scanner.Text();
		if (line != "") {
			//Log.Info(line)
			for _, c := range line{
				grid.SetValue(x, y, int(c));
				x++;
			}
		}
		x = 0;
		y++;
	}

	Log.Info("\nInitial State\n%s", grid.PrintAscii());

	prevGen := grid.Clone();

	xMin := grid.MinRow();
	xMax := grid.MaxRow();

	yMin := grid.MinCol();
	yMax := grid.MaxCol();

	tileBug := int('#');
	tileEmpty := int('.');

	uniqueScores := make(map[int]int);

	maxSteps := 1000000;
	currStep := 0;
	for{
		if(currStep >= maxSteps){
			break;
		}
		grid.CopyTo(prevGen);
		//prevGen = grid.Clone();



		for j := yMin; j<= yMax; j++{
			for i := xMin; i<= xMax; i++{
				if(!prevGen.HasValue(i, j)){
					continue;
				} else{
					val := prevGen.GetValue(i, j);
					nBug := prevGen.GetValue(i, j-1) == tileBug;
					sBug := prevGen.GetValue(i, j+1) == tileBug;
					eBug := prevGen.GetValue(i+1, j) == tileBug;
					wBug := prevGen.GetValue(i-1,j) == tileBug;
					tBug := 0;
					if(nBug){
						tBug++;
					}
					if(sBug){
						tBug++;
					}
					if(eBug){
						tBug++;
					}
					if(wBug){
						tBug++;
					}

					if(val == tileBug){
						if(tBug != 1){
							grid.SetValue(i, j, tileEmpty);
						}
					} else if(val == tileEmpty){
						if(tBug == 1 || tBug == 2){
							grid.SetValue(i, j, tileBug);
						}
					}

				}
			}
		}

		pot := 1;
		score := 0;
		for j := yMin; j<= yMax; j++{
			for i := xMin; i<= xMax; i++{
				val := grid.GetValue(i, j);
				if(val == tileBug){
					score += pot;
				}
				pot *= 2;
			}
		}
		currStep++;

		_, exists := uniqueScores[score];
		if(exists){
			Log.Info("Repeat score at step %d : %d", currStep, score);
			break;
		}
		uniqueScores[score] = currStep;
		//Log.Info("\nAfter %d minute\n%s", currStep, grid.PrintAscii());
	}

	Log.Info("Simulation complete!")
}

