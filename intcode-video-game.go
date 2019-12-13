package main

import (
	"errors"
	"math/big"
)

const VideoGameEmptyTile = 0;
const VideoGameWallTile = 1;
const VideoGameBlockTile = 2;
const VideoGamePaddleTile = 3;
const VideoGameBallTile = 4;

const VideoGameNeutralPos = 0;
const VideoGameLeftPos = -1;
const VideoGameRightPos = 1;

const VideoGameScoreSignifierX = -1;
const VideoGameScoreSignifierY = 0;


type IntCodeVideoGame struct {
	Processor *IntcodeMachineV3;
	Grid *IntegerGrid2D;
	InstructionFileName string;
	TileCount int;
	CurrPaddlePos *IntVec2;
	LastPaddlePos *IntVec2;

	CurrBallPos *IntVec2;
	LastBallPos *IntVec2;
	CurrentScore int;
	GameOver bool;
	TickCount int;
}

func (this *IntCodeVideoGame) Init(instructionFile string, grid *IntegerGrid2D) error {
	this.Processor = &IntcodeMachineV3{};
	this.Processor.PauseOnOutput = true;
	this.Processor.PauseOnInput = true;
	this.CurrentScore = -1;
	this.InstructionFileName = instructionFile;
	err := this.Processor.Load(this.InstructionFileName)
	this.Grid = grid;
	if err != nil {
		return err;
	}
	return nil;

}

func (this *IntCodeVideoGame) EnableFreePlay() {
	this.Processor.SetValueAtRegister(0, big.NewInt(2));
}

func (this *IntCodeVideoGame) MoveJoystick()  {
	move := VideoGameNeutralPos;
	if(this.LastBallPos != nil){
		if(this.LastBallPos.X < this.LastPaddlePos.X){
			move = VideoGameLeftPos;
		}
		if(this.LastBallPos.X > this.LastPaddlePos.X){
			move = VideoGameRightPos;
		}
	}
	this.Processor.QueueInput(int64(move));
}

func (this *IntCodeVideoGame) Play() error  {
	for{
		this.MoveJoystick();
		this.ReadOutput();
		if(this.CountBlocks() <= 0){
			Log.Info("Game complete!")
			break;
		}
		if(this.GameOver){
			Log.Info("GAME OVER");
			this.RenderGameBoard();
			break;
		}
		this.TickCount++;
		this.RenderGameBoard();
		//if(this.TickCount >= 10){
		//	Log.Info("Early exit")
		//	break;
		//}
	}
	return nil;
}

func (this *IntCodeVideoGame) RenderGameBoard() {
	Log.Info("\nStep %d\n%s\nScore:%d\n", this.TickCount, this.Print(), this.CurrentScore);
}

func (this *IntCodeVideoGame) ReadOutput() error {
	for{
		if(this.IsComplete()){
			return nil;
		}
		coordX, err, halted := this.Processor.ReadNextOutput();
		if(err != nil){
			return err;
		}
		if(halted){
			return nil;
		}
		if(this.Processor.PendingInput){
			Log.Info("Breaking as we are pending input")
			return nil;
		}
		coordY, err, halted := this.Processor.ReadNextOutput();
		if(err != nil){
			return err;
		}
		if(halted){
			return errors.New("Unexpected halt while reading yCoord");
		}
		val, err, halted := this.Processor.ReadNextOutput();
		if(err != nil){
			return err;
		}
		if(halted){
			return errors.New("Unexpected halt while reading val");
		}
		if(coordX == VideoGameScoreSignifierX && coordY == VideoGameScoreSignifierY){
			Log.Info("Score update %d", val);
			if(int(val) < this.CurrentScore){
				this.GameOver = true;
			}
			this.CurrentScore = int(val);
		} else{

			blockType := int(val);
			if(blockType == VideoGamePaddleTile){
				this.LastPaddlePos = &IntVec2{};
				this.LastPaddlePos.X = int(coordX);
				this.LastPaddlePos.Y = int(coordY);
			} else if(blockType == VideoGameBallTile){
				this.LastBallPos = &IntVec2{};
				this.LastBallPos.X = int(coordX);
				this.LastBallPos.Y = int(coordY);
			}

			this.Grid.SetValue(int(coordX), int(coordY), blockType);
		}

	}
}

func (this *IntCodeVideoGame) ParseMap() error {
	for{
		if(this.IsComplete()){
			return nil;
		}
		coordX, err, halted := this.Processor.ReadNextOutput();
		if(err != nil){
			return err;
		}
		if(halted){
			return nil;
		}
		coordY, err, halted := this.Processor.ReadNextOutput();
		if(err != nil){
			return err;
		}
		if(halted){
			return errors.New("Unexpected halt while reading yCoord");
		}
		val, err, halted := this.Processor.ReadNextOutput();
		if(err != nil){
			return err;
		}
		if(halted){
			return errors.New("Unexpected halt while reading val");
		}
		this.Grid.SetValue(int(coordX), int(coordY), int(val));
		this.TileCount++;
	}
}



func (this *IntCodeVideoGame) Print() string {
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
				if(val > 0){
					buff += this.PrintGameObject(this.Grid.GetValue(i, j));
				} else{
					buff += " ";
				}
			}
		}
		buff += "\n";
	}
	//Log.Info("Furthest point is %d,%d", furthestX, furthestY);

	return buff;
}

func (this *IntCodeVideoGame) PrintGameObject(gameObject int) string {
	switch(gameObject){
		case VideoGameWallTile:
			return "|";
		case VideoGameBlockTile:
			return "█";
		case VideoGameBallTile:
			return "Ø";
		case VideoGamePaddleTile:
			return "-";
	}
	return "";
}


func (this *IntCodeVideoGame) CountBlocks() int {
	xMin := this.Grid.MinRow();
	xMax := this.Grid.MaxRow();

	yMin := this.Grid.MinCol();
	yMax := this.Grid.MaxCol();

	sum := 0;
	for j := yMin; j<= yMax; j++{
		for i := xMin; i<= xMax; i++{
			if(this.Grid.GetValue(i, j) == VideoGameBlockTile){
				sum++;
			}
		}
	}
	return sum;
}


func (this *IntCodeVideoGame) IsComplete() bool{
	return this.Processor.HasHalted;
}