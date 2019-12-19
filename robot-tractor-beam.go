package main

type RobotTractorBeam struct {
	Processor *IntcodeMachineV3;
	Grid *IntegerGrid2D;
	InstructionFileName string;
}



func (this *RobotTractorBeam) Init(instructionFile string, grid *IntegerGrid2D) error {
	this.Processor = &IntcodeMachineV3{};
	this.Processor.PauseOnOutput = true;
	this.InstructionFileName = instructionFile;
	this.Grid = grid;

	err := this.Processor.Load(this.InstructionFileName)

	if err != nil {
		return err;
	}
	return nil;

}

func (this *RobotTractorBeam) ParseGrid(maxX int, maxY int) (int, error) {
	Log.Info("Starting grid parse");
	total := 0;
	for j := 0; j < maxY; j++ {
		for i := 0; i < maxX; i++ {
			this.Processor.QueueInput(int64(i));
			this.Processor.QueueInput(int64(j));
			res, err, _ := this.Processor.ReadNextOutput();
			if(err != nil){
				return -1, err;
			}
			this.Grid.SetValue(i, j, int(res));
			if(res > 0){
				total++;
			}
			this.Processor.Reset();
		}
		Log.Info("Rendered row %d", j);
	}
	return total, nil;
}


