package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Problem25A struct {

}

func (this *Problem25A) Solve() {
	Log.Info("Problem 25A solver beginning!")


	robot := &RobotAdventurer{};
	err := robot.Init("source-data/input-day-25a.txt");
	if(err != nil){
		Log.FatalError(err);
	}

	items := make([]string, 0);
	items = append(items, "klein bottle");
	items = append(items, "loom");
	items = append(items, "mutex");
	items = append(items, "pointer");
	items = append(items, "polygon");
	items = append(items, "hypercube");
	items = append(items, "mug");
	items = append(items, "manifold");

	indicies := make([]int, 0);
	for i, _ := range items{
		indicies = append(indicies, i);
	}

	combinations := Combinations(len(indicies), indicies);
	/*
	for i := range combinations{
		for j := range combinations[i] {
			fmt.Println(combinations[i][j], " ")
		}
	}*/

	const TooLightFail = "A loud, robotic voice says \"Alert! Droids on this ship are heavier than the detected value!\" and you are ejected back to the checkpoint.";
	const TooHeavyFail = "A loud, robotic voice says \"Alert! Droids on this ship are lighter than the detected value!\" and you are ejected back to the checkpoint."

	for i := range combinations{
		for _, c := range combinations[i] {
			robot.LoadSaveState("source-data/input-day-25a-save-state.txt");
			//robot.PrintGrid = false;
			_, _, err := robot.ReadState();
			if(err != nil){
				Log.FatalError(err);
			}
			Log.Info("Building combination " + fmt.Sprintln(c, " "));
			itemStr := make([]string, 0);
			for _, v := range c {
				item := items[v];
				itemStr = append(itemStr, item);
				robot.ProcessCommand("take " + item);
			}
			Log.Info("Trying with %s", itemStr);
			res, _, err :=robot.ProcessCommand("west");

			if(!strings.Contains(res, TooLightFail) && !strings.Contains(res, TooHeavyFail)){
				Log.Info("%s got result \n%s", itemStr, res);
				os.Exit(1);
			}

		}

	}


	prompt, _, err := robot.ReadState();
	if(err != nil){
		Log.FatalError(err);
	}
	Log.Info("Robot ready to explore! %s\nCommand?", prompt);

	for{
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text);
		if(text == "quit"){
			break;
		}
		res, _, err := robot.ProcessCommand(text);
		if(err != nil){
			Log.FatalError(err);
		}
		if(res != ""){
			fmt.Println(res);
		}
	}
	Log.Info("Exiting");
	for _, j := range robot.CommandHistory{
		fmt.Println(j);
	}
}
