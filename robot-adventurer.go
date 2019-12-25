package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
	"unicode"
)

const AdventureCommandNorth = "north";
const AdventureCommandSouth = "south";
const AdventureCommandWest = "west";
const AdventureCommandEast = "east";
const AdventureCommandInv = "inv";
const AdventureCommandTake = "take";
const AdventureCommandDrop = "drop";

type RobotAdventurer struct {
	Processor *IntcodeMachineV3;
	Grid *IntegerGrid2D;
	InstructionFileName string;
	CurrentPos *IntVec2;
	CommandHistory []string;
	PrintGrid bool;
}

func (this *RobotAdventurer) Init(instructionFile string) error {
	this.Grid = &IntegerGrid2D{};
	this.Grid.Init();
	this.Processor = &IntcodeMachineV3{};
	this.Processor.PauseOnOutput = true;
	this.Processor.PauseOnInput = true;
	this.InstructionFileName = instructionFile;

	err := this.Processor.Load(this.InstructionFileName)

	this.CurrentPos = &IntVec2{};
	this.CurrentPos.X = 0;
	this.CurrentPos.Y = 0;

	this.CommandHistory = make([]string, 0);


	this.MarkCurrentLocationVisited();

	if err != nil {
		return err;
	}
	return nil;

}

func (this *RobotAdventurer) MarkCurrentLocationVisited() {
	this.Grid.SetValue(this.CurrentPos.X, this.CurrentPos.Y, int('#'));
}

func (this *RobotAdventurer) ReadOutput() (string, bool, error) {
	outputBuff := make([]int, 0);
	for{
		res, err, hasHalted := this.Processor.ReadNextOutput();
		if(err != nil){
			return "", false, err;
		}
		if(hasHalted){
			break;
		}
		outputBuff = append(outputBuff, int(res));
		if(this.Processor.PendingInput){
			return this.RenderOutput(outputBuff), true, nil;
		}
	}
	return this.RenderOutput(outputBuff), false, nil;
}

func (this *RobotAdventurer) ReadState() (string, bool, error) {
	raw, pendingInput, err := this.ReadOutput();
	if(err != nil){
		return "", pendingInput, err;
	}
	xMin := this.Grid.MinRow();
	xMax := this.Grid.MaxRow();

	yMin := this.Grid.MinCol();
	yMax := this.Grid.MaxCol();

	buff := "";
	for j := yMin; j<= yMax; j++{
		for i := xMin; i<= xMax; i++{
			if(!this.Grid.HasValue(i, j)){
				buff += " ";
			} else if(i == this.CurrentPos.X && j == this.CurrentPos.Y){
				buff += fmt.Sprintf("%c", int('■'));
			} else{
				val := this.Grid.GetValue(i, j);
				buff += fmt.Sprintf("%c", val);
			}
		}
		buff += "\n";
	}

	return fmt.Sprintf("\n%s\n%s", buff, raw), pendingInput, err;
}

func (this *RobotAdventurer) LoadSaveState(fileName string) (error) {

	file, err := os.Open(fileName);
	if err != nil {
		Log.FatalError(err);
	}


	this.Grid = &IntegerGrid2D{};
	this.Grid.Init();
	this.CurrentPos.X = 0;
	this.CurrentPos.Y = 0;
	this.MarkCurrentLocationVisited();

	this.Processor = &IntcodeMachineV3{};
	this.Processor.PauseOnOutput = true;
	this.Processor.PauseOnInput = true;

	err = this.Processor.Load(this.InstructionFileName)
	if(err != nil){
		return err;
	}

	this.CommandHistory = make([]string, 0);

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text());
		if (line != "") {
			this.ProcessCommand(line);
		}
	}
	return err;
}

func (this *RobotAdventurer) ProcessCommand(command string) (string, bool, error) {
	this.CommandHistory = append(this.CommandHistory, command);
	switch(command){
		case AdventureCommandEast:
			this.CurrentPos.X++;
			this.MarkCurrentLocationVisited();
			this.WriteInstruction(AdventureCommandEast);
			return this.ReadState();
			break;
		case AdventureCommandWest:
			this.CurrentPos.X--;
			this.MarkCurrentLocationVisited();
			this.WriteInstruction(AdventureCommandWest);
			return this.ReadState();
			break;
		case AdventureCommandNorth:
			this.CurrentPos.Y--;
			this.MarkCurrentLocationVisited();
			this.WriteInstruction(AdventureCommandNorth);
			return this.ReadState();
			break;
		case AdventureCommandSouth:
			this.CurrentPos.Y++;
			this.MarkCurrentLocationVisited();
			this.WriteInstruction(AdventureCommandSouth);
			return this.ReadState();
			break;
		case AdventureCommandInv:
			this.WriteInstruction(AdventureCommandInv);
			return this.ReadState();
			break;
		default:
			parts := strings.Split(command, " ");
			if(parts[0] == AdventureCommandTake){
				this.WriteInstruction(strings.TrimSpace(command));
				return this.ReadState();
			} else if(parts[0] == AdventureCommandDrop){
				this.WriteInstruction(strings.TrimSpace(command));
				return this.ReadState();
			}
			break;
	}
	Log.Info("Unknown command %s", command);
	return "", true, nil;
}



func (this *RobotAdventurer) RenderOutput(output []int) string {
	var str strings.Builder;
	for _, v := range output{
		if(v < unicode.MaxASCII){
			str.WriteByte(byte(v));
		} else{
			str.WriteString(fmt.Sprintf("%d", v));
		}
	}
	val := str.String();

	lines := strings.Split(val, "\n");
	doorsIndex := int(math.MaxInt64);
	x := this.CurrentPos.X;
	y := this.CurrentPos.Y;
	for i, l := range lines{
		trimmed := strings.TrimSpace(l);
		if(trimmed == "Doors here lead:"){
			doorsIndex = i;
			continue;
		}
		if(i > doorsIndex){
			if(trimmed == "- north"){
				if(!this.Grid.HasValue(x, y-1)){
					this.Grid.SetValue(x, y-1, int('?'));
				}
			}
			if(trimmed == "- south"){
				if(!this.Grid.HasValue(x, y+1)){
					this.Grid.SetValue(x, y+1, int('?'));
				}
			}
			if(trimmed == "- west"){
				if(!this.Grid.HasValue(x-1, y)){
					this.Grid.SetValue(x-1, y, int('?'));
				}
			}
			if(trimmed == "- east"){
				if(!this.Grid.HasValue(x+1, y)){
					this.Grid.SetValue(x+1, y, int('?'));
				}
			}
		}
	}
	return val;
}


func (this *RobotAdventurer) WriteInstruction(instructionLiteral string) {
	for _, c := range instructionLiteral{
		this.Processor.QueueInput(int64(c));
	}
	this.Processor.QueueInput(int64('\n'));
}
