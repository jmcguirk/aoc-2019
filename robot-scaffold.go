package main

import (
	"math/big"
	"strings"
)

type RobotScaffold struct {
	Processor *IntcodeMachineV3;
	Grid *IntegerGrid2D;
	InstructionFileName string;
	CurrentPos *IntVec2;
}

const AsciiCodeScaffold = 35;
const AsciiCodeEmpty = 46;
const AsciiCodeNewLine = 10;
const AsciiRobotUp = 94;
const AsciiRobotUnvisited = 100;
const AsciiRobotVisited = 101;

func (this *RobotScaffold) Init(instructionFile string, grid *IntegerGrid2D) error {
	this.Processor = &IntcodeMachineV3{};
	this.Processor.PauseOnOutput = true;
	this.Processor.PauseOnInput = true;
	this.InstructionFileName = instructionFile;
	this.Grid = grid;

	err := this.Processor.Load(this.InstructionFileName)

	if err != nil {
		return err;
	}
	return nil;

}

func (this *RobotScaffold) ParseGrid() error {
	Log.Info("Starting grid parse");
	row:= 0;
	col:= 0;
	for{
		res, err, halted := this.Processor.ReadNextOutput();
		if(err != nil){
			return err;
		}
		if(halted){
			break;
		}
		if(res == AsciiCodeNewLine){
			row++;
			col = 0;
		} else {
			if(res == AsciiRobotUp){
				this.CurrentPos = &IntVec2{};
				this.CurrentPos.X = col;
				this.CurrentPos.Y = row;
			}
			this.Grid.SetValue(col, row, int(res));
			col++;
		}
	}
	return nil;
}



func (this *RobotScaffold) Execute(mainProgram string, function1 string, function2 string, function3 string) int64 {
	this.Processor.Reset();
	this.Processor.SetValueAtRegister(0, big.NewInt(2)); // Turn the robot on
	this.Processor.PauseOnInput = false;
	this.InputInstructionStream(mainProgram);
	this.InputInstructionStream(function1);
	this.InputInstructionStream(function2);
	this.InputInstructionStream(function3);

	this.SetVideoOutput(false);

	for {
		_, err, halted := this.Processor.ReadNextOutput()
		if(err != nil){
			Log.FatalError(err);
		}
		if (halted) {
			return this.Processor.LastOutputValue.Int64();
		}
	}

}

func (this *RobotScaffold) InputInstructionStream(mainProgram string) {
	parts := strings.Split(mainProgram,",");
	for i, v := range parts{
		if(i > 0){
			this.Processor.QueueInput(int64(int(',')));
		}
		for _, c := range v {
			this.Processor.QueueInput(int64(int(c)));
		}
	}
	this.Processor.QueueInput(int64(int(AsciiCodeNewLine)));
}

func (this *RobotScaffold) SetVideoOutput(on bool) {
	if(on){
		this.Processor.QueueInput(int64(int('y')));
	} else{
		this.Processor.QueueInput(int64(int('n')));
	}
	this.Processor.QueueInput(int64(int(AsciiCodeNewLine)));
}

func (this *RobotScaffold) PrintGrid() string {
	xMin := this.Grid.MinRow();
	xMax := this.Grid.MaxRow();

	yMin := this.Grid.MinCol();
	yMax := this.Grid.MaxCol();

	buff := "";
	for j := yMin; j<= yMax; j++{
		for i := xMin; i<= xMax; i++{
			if(!this.Grid.HasValue(i, j)){
				buff += " ";
			} else{
				val := this.Grid.GetValue(i, j);
				if(val == AsciiCodeEmpty){
					buff += "."
				} else if(val == AsciiCodeScaffold) {
					buff += "#";
				} else if (val == AsciiRobotUp){
					buff += "^";
				} else
				{
					Log.Info("unknown value %d", val);
					buff += "?";
				}
			}
		}
		buff += "\n";
	}
	return buff;
}




func PrintScaffoldGrid(grid *IntegerGrid2D) string {
	xMin := grid.MinRow();
	xMax := grid.MaxRow();

	yMin := grid.MinCol();
	yMax := grid.MaxCol();

	buff := "";
	for j := yMin; j<= yMax; j++{
		for i := xMin; i<= xMax; i++{
			if(!grid.HasValue(i, j)){
				buff += " ";
			} else{
				val := grid.GetValue(i, j);
				if(val == AsciiCodeEmpty){
					buff += "."
				} else if(val == AsciiCodeScaffold) {
					buff += "#";
				} else if (val == AsciiRobotUp){
					buff += "^";
				} else if(val == AsciiRobotUnvisited){
					buff += "?";
				} else if(val == AsciiRobotVisited){
					buff += "*";
				} else
				{
					Log.Info("unknown value %d", val);
					buff += "?";
				}
			}
		}
		buff += "\n";
	}
	return buff;
}

func IsScaffoldIntersection(i int, j int, grid *IntegerGrid2D) bool {
	val := grid.GetValue(i, j);
	if(val == AsciiCodeScaffold) {
		if(grid.GetValue(i-1, j) == AsciiCodeScaffold && grid.GetValue(i+1, j) == AsciiCodeScaffold && grid.GetValue(i, j+1) == AsciiCodeScaffold && grid.GetValue(i, j-1) == AsciiCodeScaffold) {
			return true;
		}
	}
	return false;
}

func (this *RobotScaffold) CalculateCheckSum() int {
	sum := 0;
	xMin := this.Grid.MinRow();
	xMax := this.Grid.MaxRow();

	yMin := this.Grid.MinCol();
	yMax := this.Grid.MaxCol();

	for j := yMin; j<= yMax; j++{
		for i := xMin; i<= xMax; i++{
			val := this.Grid.GetValue(i, j);
			if(val == AsciiCodeScaffold) {
				if(this.Grid.GetValue(i-1, j) == AsciiCodeScaffold && this.Grid.GetValue(i+1, j) == AsciiCodeScaffold && this.Grid.GetValue(i, j+1) == AsciiCodeScaffold && this.Grid.GetValue(i, j-1) == AsciiCodeScaffold) {
					sum += i * j;
				}
			}
		}

	}
	return sum;
}