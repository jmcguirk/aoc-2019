package main

import (
	"fmt"
	"strings"
	"unicode"
)

type RobotSpringDroid struct {
	Processor *SpringCodeProcessor;
	Grid *IntegerGrid2D;
}




func (this *RobotSpringDroid) Init(instructionFile string) error {
	this.Processor = &SpringCodeProcessor{};
	err := this.Processor.Init(instructionFile, 15);
	if err != nil {
		return err;
	}
	this.Grid = &IntegerGrid2D{};
	this.Grid.Init();

	return nil;

}

func (this *RobotSpringDroid) LoadProgram(instructionFile string) error {
	return this.Processor.LoadProgram(instructionFile);
}

func (this *RobotSpringDroid) Execute(run bool) error {
	if(run){
		this.Processor.SetBeginExecutionKeyword("RUN");
	} else{
		this.Processor.SetBeginExecutionKeyword("WALK");
	}
	err, out := this.Processor.ExecuteLoadedProgram();
	if(err != nil){
		return err;
	}
	Log.Info("Processed loaded program - output buff was %d long", len(out));
	if(len(out) > 1){
		Log.Info("Output: \n%s", this.RenderOutput(out));
	}

	return nil;
}

func (this *RobotSpringDroid) RenderOutput(output []int) string {
	var str strings.Builder;
	for _, v := range output{
		if(v < unicode.MaxASCII){
			str.WriteByte(byte(v));
		} else{
			str.WriteString(fmt.Sprintf("%d", v));
		}
	}
	return str.String();
}

func (this *RobotSpringDroid) DescribeLoadedProgram() string {
	return this.Processor.Describe();
}
