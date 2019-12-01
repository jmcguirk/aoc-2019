package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type Problem1A struct {

}

func (this *Problem1A) Solve() {
	Log.Info("Problem 1A solver beginning!")


	file, err := os.Open("source-data/input-day-01a.txt");
	if err != nil {
		Log.FatalError(err);
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var totalFuelRequired int64 = 0;
	for scanner.Scan() {             // internally, it advances token based on sperator
		line := strings.TrimSpace(scanner.Text());
		if(line != ""){
			mass, err := strconv.ParseInt(line, 10, 64);
			if(err != nil){
				Log.FatalError(err);
			}
			fuelRequired := mass /3;
			fuelRequired = fuelRequired -2;
			//Log.Info("Fuel Required For %d is %d", mass, fuelRequired);
			totalFuelRequired += fuelRequired;
		}
	}
	Log.Info("Finished parsing file - total fuel required is %d", totalFuelRequired);
}
