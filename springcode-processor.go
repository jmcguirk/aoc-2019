package main

import (
	"errors"
	"io/ioutil"
	"strconv"
	"strings"
)

type SpringCodeProcessor struct {
	Processor *IntcodeMachineV3;
	InstructionFileName string;
	MaxInstructionQueue int;
	CurrentProgram string;
	BeginExecutionKeyword string;
	Program []SpringCodeInstruction;
}

func (this *SpringCodeProcessor) Init(instructionFile string, maxInstructionQueue int) error {
	this.InstructionFileName = instructionFile;
	this.Program = make([]SpringCodeInstruction, 0);
	this.Processor = &IntcodeMachineV3{};
	this.Processor.PauseOnOutput = true;
	this.MaxInstructionQueue = maxInstructionQueue;
	err := this.Processor.Load(this.InstructionFileName)

	if err != nil {
		return err;
	}
	return nil;

}

func (this *SpringCodeProcessor) SetBeginExecutionKeyword(word string){
	this.BeginExecutionKeyword = word;
}

func (this *SpringCodeProcessor) ExecuteLoadedProgram() (error, []int) {
	for _, instruction := range this.Program {
		instruction.SerializeToInput(this);
	}
	outputBuff := make([]int, 0);
	this.FinalizeProgram();
	for{
		res, err, hasHalted := this.Processor.ReadNextOutput();
		if(err != nil){
			return err, nil;
		}
		if(hasHalted){
			break;
		}
		outputBuff = append(outputBuff, int(res));
	}
	return nil, outputBuff;
}

func (this *SpringCodeProcessor) FinalizeProgram() () {
	if(this.BeginExecutionKeyword != ""){
		for _, v := range this.BeginExecutionKeyword{
			this.EnqueueRawInput(int(v));
		}
		this.EnqueueRawInput(AsciiNewLine);
	}
}

func (this *SpringCodeProcessor) LoadProgram(instructionFile string) error {

	Log.Info("Loading springcode program from " + instructionFile);
	this.Program = make([]SpringCodeInstruction, 0);

	file, err := ioutil.ReadFile(instructionFile);
	if err != nil {
		return err;
	}
	fileContents := strings.TrimSpace(string(file));
	parts := strings.Split(fileContents, "\n");
	lineNum := 1;
	for _, val := range parts {
		trimmed := strings.TrimSpace(val);
		if(trimmed != ""){
			lineParts := strings.Split(trimmed, " ");

			opCode := lineParts[0];
			var instruction SpringCodeInstruction;
			switch(opCode){
				case SpringCodeOpAnd:
					instruction = &SpringCodeAND{};
					break;
				case SpringCodeOpNot:
					instruction = &SpringCodeNOT{};
					break;
				case SpringCodeOpOr:
					instruction = &SpringCodeOR{};
					break;
				case "//": // Comment
					continue;
				default:
					return errors.New("Unknown opcode " + opCode + " at line num " + strconv.Itoa(lineNum));
			}
			if(instruction != nil){
				err = instruction.Parse(trimmed, lineNum);
				if(err != nil){
					return err;
				}
				this.Program = append(this.Program, instruction);
				lineNum++;
			}
		}
	}
	//Log.Info("Successfully parsed springcode program from " + instructionFile);
	this.CurrentProgram = instructionFile;
	this.Processor.Reset();
	return nil;
}

func (this *SpringCodeProcessor) Describe() string  {
	var str strings.Builder;
	str.WriteString("\n Spring Code Processor - Program Loaded From " + this.CurrentProgram + "\n");
	for _, instruction := range this.Program{
		str.WriteString(instruction.ToString() + "\n");
	}
	return str.String();
}

func (this *SpringCodeProcessor) AddInstruction(instruction *SpringCodeInstruction) error {
	if(len(this.Program) >= this.MaxInstructionQueue){
		return errors.New("At program queue length max " + strconv.Itoa(this.MaxInstructionQueue));
	}

	return nil;
}

func (this *SpringCodeProcessor) EnqueueRawInput(char int) {
	this.Processor.QueueInput(int64(char));
}


func (this *SpringCodeProcessor) Reset()  {
	this.Processor.Reset();

}
