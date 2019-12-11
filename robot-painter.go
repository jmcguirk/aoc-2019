package main

type RobotPainter struct {
	Processor *IntcodeMachineV3;
	Grid *IntegerGrid2D;
	Position *IntVec2;
	Id string;
	InstructionFileName string;
	Orientation int;
	StepCount int;
	WorkDone int;
}

func (this *RobotPainter) Init(instructionFile string, id string, initialPosition *IntVec2, grid *IntegerGrid2D) error {
	this.Processor = &IntcodeMachineV3{};
	this.Processor.PauseOnOutput = true;
	this.InstructionFileName = instructionFile;
	this.Id = id;
	this.Position = &IntVec2{};
	this.Position.X = initialPosition.X;
	this.Position.Y = initialPosition.Y;
	this.Grid = grid;

	this.Orientation = OrientationNorth;
	this.StepCount = 0;

	err := this.Processor.Load(this.InstructionFileName)

	if err != nil {
		return err;
	}
	return nil;

}


func (this *RobotPainter) ReadCurrentColor(){
	this.Processor.QueueInput(int64(this.Grid.GetValue(this.Position.X, this.Position.Y)));
}

func (this *RobotPainter) MarkColor(color int){
	Log.Info("[Step %d] Painter Robot %s marking color %d at %s", this.StepCount, this.Id, color, this.Position.ToString());
	if(!this.Grid.HasValue(this.Position.X, this.Position.Y)){
		this.WorkDone++;
	}
	this.Grid.SetValue(this.Position.X, this.Position.Y, color);
}

func (this *RobotPainter) Turn(direction int){
	oldOrientation := this.Orientation;
	if(direction == 0){
		Log.Info("[Step %d] Painter Robot turning left at %s", this.StepCount, this.Position.ToString());
		this.Orientation--;
		if(this.Orientation < OrientationNorth){
			this.Orientation = OrientationWest;
		}
		Log.Info("[Step %d] Left turn complete %s to %s", this.StepCount, PrintOrientation(oldOrientation), PrintOrientation(this.Orientation));
	} else {
		Log.Info("[Step %d] Painter Robot turning right at %s", this.StepCount, this.Position.ToString());
		this.Orientation++;
		if(this.Orientation > OrientationWest){
			this.Orientation = OrientationNorth;
		}
		Log.Info("[Step %d] Right turn complete %s to %s", this.StepCount, PrintOrientation(oldOrientation), PrintOrientation(this.Orientation));
	}

}

func (this *RobotPainter) Step() error{
	this.StepCount++;
	Log.Info("[Step %d] Painter Robot %s beginning step at %s facing %s. Work Done - %d", this.StepCount, this.Id, this.Position.ToString(), PrintOrientation(this.Orientation), this.WorkDone);
	this.ReadCurrentColor();
	color, err := this.ReadOutput();
	if(err != nil){
		return err;
	}
	this.MarkColor(color);
	turn, err := this.ReadOutput();
	if(err != nil){
		return err;
	}
	this.Turn(turn);
	this.Advance();
	return nil;
}

func (this *RobotPainter) Advance(){
	newPos := &IntVec2{};
	x := this.Position.X;
	y := this.Position.Y;

	switch(this.Orientation){
		case OrientationNorth:
			y--;
			break;
		case OrientationSouth:
			y++;
			break;
		case OrientationWest:
			x--;
			break;
		case OrientationEast:
			x++;
			break;
	}


	newPos.X = x;
	newPos.Y = y;

	Log.Info("[Step %d] Painter Robot advanced from %s to %s", this.StepCount, this.Position.ToString(), newPos.ToString());
	this.Position = newPos;

}

func (this *RobotPainter) ReadOutput() (int, error){
	err := this.Processor.Execute();
	if(err != nil){
		return -1, err;
	}
	return int(this.Processor.LastOutputValue.Int64()), nil;
}

func (this *RobotPainter) PrintState(){
	Log.Info("[Step %d] Painter Robot %s, is at %s", this.StepCount, this.Id, this.Position.ToString());
}

func (this *RobotPainter) IsComplete() bool{
	return this.Processor.HasHalted;
}