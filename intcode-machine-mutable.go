package main

import (
	"errors"
	"io/ioutil"
	"strconv"
	"strings"
)

type IntcodeMachineMutable struct {
	Registers []int64;
	InitialState []int64;
	InstructionPointer int64;
	InstructionsExecuted int;
}

const InstructionPointerJumpSize = 4;



func (this *IntcodeMachineMutable) Load(fileName string) error {
	Log.Info("Loading intcode machine from %s", fileName)
	this.Registers = make([]int64, 0);

	file, err := ioutil.ReadFile(fileName);
	if err != nil {
		return err;
	}
	fileContents := strings.TrimSpace(string(file));
	parts := strings.Split(fileContents, ",");
	for _, val := range parts {
		trimmed := strings.TrimSpace(val);
		if(trimmed != ""){
			parsed, err := strconv.ParseInt(val, 10, 64);
			if(err != nil){
				return err;
			}
			this.Registers = append(this.Registers, parsed);
		}
	}

	// Reset the function pointer to start at the beginning
	this.InstructionPointer = 0;


	this.InitialState = make([]int64, len(this.Registers));
	copy(this.InitialState, this.Registers);


	Log.Info("Finished parsing machine initial state - %d registers", len(this.Registers));

	return nil;
}

func (this *IntcodeMachineMutable) Reset() {
	//Log.Info("Performing reset");
	this.InstructionPointer = 0;
	copy(this.Registers, this.InitialState);
}

func (this *IntcodeMachineMutable) Execute() error {
	for {
		if(this.InstructionPointer < 0 || this.InstructionPointer > int64(len(this.Registers))){
			return errors.New("Instruction pointer went outside bounds " + strconv.FormatInt(this.InstructionPointer, 10));
		}
		instruction := this.Registers[this.InstructionPointer];
		switch(instruction){
			case IntCodeOpCodeAdd:
				err := this.ExecuteAdd();
				if(err != nil){
					return err;
				}
				break;
			case IntCodeOpCodeMul:
				err := this.ExecuteMul();
				if(err != nil){
					return err;
				}
				break;
			case IntCodeOpCodeHalt:
				//Log.Info("Program halting after %d instructions executed", this.InstructionsExecuted);
				return nil;
				break;
		}
		this.InstructionPointer += InstructionPointerJumpSize;
		this.InstructionsExecuted++;
	}
}

func (this *IntcodeMachineMutable) ExecuteAdd() error {
	term1Register := this.InstructionPointer+1;
	term2Register := this.InstructionPointer+2;
	destRegister := this.InstructionPointer+3;
	if(term1Register >= int64(len(this.Registers)) || term2Register >= int64(len(this.Registers)) || destRegister >= int64(len(this.Registers))){
		return errors.New("Add instruction went outside bounds " + strconv.FormatInt(this.InstructionPointer, 10));
	}

	term1Position := this.GetValueAtRegister(term1Register);
	term2Position := this.GetValueAtRegister(term2Register);
	destPosition := this.GetValueAtRegister(destRegister);

	if(term1Position >= int64(len(this.Registers)) || term1Position < 0){
		return errors.New("Add had bad read position 1 " + strconv.FormatInt(term1Position, 10));
	}

	if(term2Position >= int64(len(this.Registers)) || term2Position < 0){
		return errors.New("Add had bad read position 2 " + strconv.FormatInt(term2Position, 10));
	}

	if(destPosition >= int64(len(this.Registers)) || destPosition < 0){
		return errors.New("Add had bad dest position " + strconv.FormatInt(destPosition, 10));
	}

	term1 := this.GetValueAtRegister(term1Position);
	term2 := this.GetValueAtRegister(term2Position);

	sum := term1 + term2;

	this.SetValueAtRegister(destPosition, sum);
	return nil;
}

func (this *IntcodeMachineMutable) ExecuteMul() error {
	term1Register := this.InstructionPointer+1;
	term2Register := this.InstructionPointer+2;
	destRegister := this.InstructionPointer+3;
	if(term1Register >= int64(len(this.Registers)) || term2Register >= int64(len(this.Registers)) || destRegister >= int64(len(this.Registers))){
		return errors.New("Mul instruction went outside bounds " + strconv.FormatInt(this.InstructionPointer, 10));
	}

	term1Position := this.GetValueAtRegister(term1Register);
	term2Position := this.GetValueAtRegister(term2Register);
	destPosition := this.GetValueAtRegister(destRegister);

	if(term1Position >= int64(len(this.Registers)) || term1Position < 0){
		return errors.New("Mul had bad read position 1 " + strconv.FormatInt(term1Position, 10));
	}

	if(term2Position >= int64(len(this.Registers)) || term2Position < 0){
		return errors.New("Mul had bad read position 2 " + strconv.FormatInt(term2Position, 10));
	}

	if(destPosition >= int64(len(this.Registers)) || destPosition < 0){
		return errors.New("Mul had bad dest position " + strconv.FormatInt(destPosition, 10));
	}

	term1 := this.GetValueAtRegister(term1Position);
	term2 := this.GetValueAtRegister(term2Position);

	sum := term1 * term2;

	this.SetValueAtRegister(destPosition, sum);
	return nil;
}

func (this *IntcodeMachineMutable) GetValueAtRegister(register int64) int64 {
	return this.Registers[register];
}

func (this *IntcodeMachineMutable) SetValueAtRegister(register int64, val int64)  {
	this.Registers[register] = val;
}

func (this *IntcodeMachineMutable) PrintContents()  {
	buff := "";
	for i, val := range this.Registers{
		if(i > 0){
			buff += ", ";
		}
		buff += strconv.FormatInt(val, 10);
	}
	Log.Info("Machine Contents - " + buff)
}
