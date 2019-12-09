package main

import (
	"errors"
	"io/ioutil"
	"strconv"
	"strings"
)

type IntcodeMachineV2 struct {
	Registers []int64;
	InitialState []int64;
	InstructionPointer int64;
	InstructionsExecuted int;
	InputQueue []int64;
	LastOutput int;
	LastOutputValue int64;
	HasHalted bool;
	PauseOnOutput bool;
}





func (this *IntcodeMachineV2) SetInputValue(val int64) {
	Log.Info("Setting input value to %d", val);
	this.InputQueue = append(this.InputQueue, val);
}

func (this *IntcodeMachineV2) QueueInput(val int64) {
	Log.Info("Queuing input value %d", val);
	this.InputQueue = append(this.InputQueue, val);
}


func (this *IntcodeMachineV2) Load(fileName string) error {
	Log.Info("Loading intcode v2 machine from %s", fileName)
	this.Registers = make([]int64, 0);
	this.InputQueue = make([]int64, 0);

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
	this.LastOutputValue = 0;

	this.InitialState = make([]int64, len(this.Registers));
	copy(this.InitialState, this.Registers);
	this.InputQueue = make([]int64, 0);


	Log.Info("Finished parsing machine initial state - %d registers", len(this.Registers));

	return nil;
}

func (this *IntcodeMachineV2) Reset() {
	//Log.Info("Performing reset");
	this.InstructionPointer = 0;
	this.LastOutputValue = 0;
	this.InputQueue = make([]int64, 0);
	this.LastOutput = 0;
	this.HasHalted = false;
	copy(this.Registers, this.InitialState);
}

func (this *IntcodeMachineV2) Execute() error {
	for {
		if(this.InstructionPointer < 0 || this.InstructionPointer > int64(len(this.Registers))){
			return errors.New("Instruction pointer went outside bounds " + strconv.FormatInt(this.InstructionPointer, 10));
		}

		instruction, err := this.ParseOperation();
		if(err != nil){
			return err;
		}

		Log.Info("[EXEC - " + strconv.Itoa(this.InstructionsExecuted) + "] " + instruction.Describe());
		//Log.Fatal("Early exit " + instruction.Describe());

		switch(instruction.Operation){
			case IntCodeOpCodeAdd:
				err := this.ExecuteAdd(instruction);
				if(err != nil){
					return err;
				}
				break;
			case IntCodeOpCodeMul:
				err := this.ExecuteMul(instruction);
				if(err != nil){
					return err;
				}
				break;
			case IntCodeOpCodeInput:
				err := this.ExecuteInput(instruction);
				if(err != nil){
					return err;
				}
				break;
			case IntCodeOpCodeOutput:
				this.LastOutput = this.InstructionsExecuted;
				err := this.ExecuteOutput(instruction);
				if(err != nil){
					return err;
				}
				if(this.PauseOnOutput){
					return nil;
				}
				break;
			case IntCodeOpCodeJumpIfTrue:
				err := this.ExecuteJumpIfTrue(instruction);
				if(err != nil){
					return err;
				}
				break;
			case IntCodeOpCodeJumpIfFalse:
				err := this.ExecuteJumpIfFalse(instruction);
				if(err != nil){
					return err;
				}
				break;
			case IntCodeOpCodeLessThan:
				err := this.ExecuteLessThan(instruction);
				if(err != nil){
					return err;
				}
				break;
			case IntCodeOpCodeEquals:
				err := this.ExecuteEquals(instruction);
				if(err != nil){
					return err;
				}
				break;
			case IntCodeOpCodeHalt:
				//Log.Info("Program halting after %d instructions executed - last output executed was %d", this.InstructionsExecuted, this.LastOutput);
				this.HasHalted = true;
				return nil;
				break;
			default:
				return errors.New("Unimplemented opcode " + instruction.Describe());
				break;
		}

		this.InstructionsExecuted++;
	}
}

func (this *IntcodeMachineV2) ParseOperation() (*IntcodeInstruction, error){
	res := &IntcodeInstruction{};
	rawVal := this.Registers[this.InstructionPointer];
	Log.Info("Parsing instruction from %d", rawVal);
	op1 := nthDigit64(rawVal, 0);
	op2 := nthDigit64(rawVal, 1);
	res.Operation = op1 + (op2 * 10);
	res.FistParameterMode = nthDigit64(rawVal, 2);
	res.SecondParameterMode = nthDigit64(rawVal, 3);
	res.ThirdParameterMode = nthDigit64(rawVal, 4);
	res.DeriveLength();

	return res, nil;
}

func (this *IntcodeMachineV2) ExecuteOutput(instruction *IntcodeInstruction) error {
	term1Register := this.InstructionPointer+1;
	if(term1Register >= int64(len(this.Registers))){
		return errors.New("Add instruction went outside bounds " + strconv.FormatInt(this.InstructionPointer, 10));
	}

	term1 := this.GetValueAtRegister(term1Register);
	if(instruction.FistParameterMode == ParameterModePosition){
		if(term1 >= int64(len(this.Registers)) || term1 < 0){
			return errors.New("Add had bad read position 1 " + strconv.FormatInt(term1, 10));
		}
		term1 = this.GetValueAtRegister(term1);
	}

	//Log.Info("[OUTPUT] - " + strconv.FormatInt(term1, 10));
	this.LastOutputValue = term1;
	this.InstructionPointer += int64(instruction.OperationLength);
	return nil;
}

func (this *IntcodeMachineV2) ExecuteAdd(instruction *IntcodeInstruction) error {

	if(instruction.ThirdParameterMode != ParameterModePosition){
		return errors.New("Add instructions can't write to immediate mode registers");
	}

	term1Register := this.InstructionPointer+1;
	term2Register := this.InstructionPointer+2;
	destRegister := this.InstructionPointer+3;
	if(term1Register >= int64(len(this.Registers)) || term2Register >= int64(len(this.Registers)) || destRegister >= int64(len(this.Registers))){
		return errors.New("Add instruction went outside bounds " + strconv.FormatInt(this.InstructionPointer, 10));
	}

	term1 := this.GetValueAtRegister(term1Register);
	if(instruction.FistParameterMode == ParameterModePosition){
		if(term1 >= int64(len(this.Registers)) || term1 < 0){
			return errors.New("Add had bad read position 1 " + strconv.FormatInt(term1, 10));
		}
		term1 = this.GetValueAtRegister(term1);
	}

	term2 := this.GetValueAtRegister(term2Register);
	if(instruction.SecondParameterMode == ParameterModePosition){
		if(term2 >= int64(len(this.Registers)) || term2 < 0){
			return errors.New("Add had bad read position 2 " + strconv.FormatInt(term2, 10));
		}
		term2 = this.GetValueAtRegister(term2);
	}



	destPosition := this.GetValueAtRegister(destRegister);

	if(destPosition >= int64(len(this.Registers)) || destPosition < 0){
		return errors.New("Add had bad dest position " + strconv.FormatInt(destPosition, 10));
	}
	sum := term1 + term2;
	this.SetValueAtRegister(destPosition, sum);
	Log.Info("executed add %d stored at %d", sum, destPosition);
	this.InstructionPointer += int64(instruction.OperationLength);
	return nil;
}

func (this *IntcodeMachineV2) ExecuteLessThan(instruction *IntcodeInstruction) error {

	if(instruction.ThirdParameterMode != ParameterModePosition){
		return errors.New("LE instructions can't write to immediate mode registers");
	}

	term1Register := this.InstructionPointer+1;
	term2Register := this.InstructionPointer+2;
	destRegister := this.InstructionPointer+3;
	if(term1Register >= int64(len(this.Registers)) || term2Register >= int64(len(this.Registers)) || destRegister >= int64(len(this.Registers))){
		return errors.New("LE instruction went outside bounds " + strconv.FormatInt(this.InstructionPointer, 10));
	}

	term1 := this.GetValueAtRegister(term1Register);
	if(instruction.FistParameterMode == ParameterModePosition){
		if(term1 >= int64(len(this.Registers)) || term1 < 0){
			return errors.New("LE had bad read position 1 " + strconv.FormatInt(term1, 10));
		}
		term1 = this.GetValueAtRegister(term1);
	}

	term2 := this.GetValueAtRegister(term2Register);
	if(instruction.SecondParameterMode == ParameterModePosition){
		if(term2 >= int64(len(this.Registers)) || term2 < 0){
			return errors.New("LE had bad read position 2 " + strconv.FormatInt(term2, 10));
		}
		term2 = this.GetValueAtRegister(term2);
	}



	destPosition := this.GetValueAtRegister(destRegister);

	if(destPosition >= int64(len(this.Registers)) || destPosition < 0){
		return errors.New("LE had bad dest position " + strconv.FormatInt(destPosition, 10));
	}
	if(term1 < term2){
		this.SetValueAtRegister(destPosition, 1);
	} else{
		this.SetValueAtRegister(destPosition, 0);
	}


	this.InstructionPointer += int64(instruction.OperationLength);
	return nil;
}

func (this *IntcodeMachineV2) ExecuteEquals(instruction *IntcodeInstruction) error {

	if(instruction.ThirdParameterMode != ParameterModePosition){
		return errors.New("EQ instructions can't write to immediate mode registers");
	}

	term1Register := this.InstructionPointer+1;
	term2Register := this.InstructionPointer+2;
	destRegister := this.InstructionPointer+3;
	if(term1Register >= int64(len(this.Registers)) || term2Register >= int64(len(this.Registers)) || destRegister >= int64(len(this.Registers))){
		return errors.New("EQ instruction went outside bounds " + strconv.FormatInt(this.InstructionPointer, 10));
	}

	term1 := this.GetValueAtRegister(term1Register);
	if(instruction.FistParameterMode == ParameterModePosition){
		if(term1 >= int64(len(this.Registers)) || term1 < 0){
			return errors.New("EQ had bad read position 1 " + strconv.FormatInt(term1, 10));
		}
		term1 = this.GetValueAtRegister(term1);
	}

	term2 := this.GetValueAtRegister(term2Register);
	if(instruction.SecondParameterMode == ParameterModePosition){
		if(term2 >= int64(len(this.Registers)) || term2 < 0){
			return errors.New("EQ had bad read position 2 " + strconv.FormatInt(term2, 10));
		}
		term2 = this.GetValueAtRegister(term2);
	}



	destPosition := this.GetValueAtRegister(destRegister);

	if(destPosition >= int64(len(this.Registers)) || destPosition < 0){
		return errors.New("EQ had bad dest position " + strconv.FormatInt(destPosition, 10));
	}
	if(term1 == term2){
		this.SetValueAtRegister(destPosition, 1);
	} else{
		this.SetValueAtRegister(destPosition, 0);
	}


	this.InstructionPointer += int64(instruction.OperationLength);
	return nil;
}

func (this *IntcodeMachineV2) ExecuteJumpIfTrue(instruction *IntcodeInstruction) error {
	term1Register := this.InstructionPointer+1;
	term2Register := this.InstructionPointer+2;
	if(term1Register >= int64(len(this.Registers)) || term2Register >= int64(len(this.Registers))){
		return errors.New("JIT instruction went outside bounds " + strconv.FormatInt(this.InstructionPointer, 10));
	}

	term1 := this.GetValueAtRegister(term1Register);
	if(instruction.FistParameterMode == ParameterModePosition){
		if(term1 >= int64(len(this.Registers)) || term1 < 0){
			return errors.New("JIT had bad read position 1 " + strconv.FormatInt(term1, 10));
		}
		term1 = this.GetValueAtRegister(term1);
	}

	term2 := this.GetValueAtRegister(term2Register);
	if(instruction.SecondParameterMode == ParameterModePosition){
		if(term2 >= int64(len(this.Registers)) || term2 < 0){
			return errors.New("JIT had bad read position 2 " + strconv.FormatInt(term2, 10));
		}
		term2 = this.GetValueAtRegister(term2);
	}

	if(term1 != 0){
		this.InstructionPointer = term2;
	} else{
		this.InstructionPointer += int64(instruction.OperationLength);
	}

	return nil;
}

func (this *IntcodeMachineV2) ExecuteJumpIfFalse(instruction *IntcodeInstruction) error {
	term1Register := this.InstructionPointer+1;
	term2Register := this.InstructionPointer+2;
	if(term1Register >= int64(len(this.Registers)) || term2Register >= int64(len(this.Registers))){
		return errors.New("JIF instruction went outside bounds " + strconv.FormatInt(this.InstructionPointer, 10));
	}

	term1 := this.GetValueAtRegister(term1Register);
	if(instruction.FistParameterMode == ParameterModePosition){
		if(term1 >= int64(len(this.Registers)) || term1 < 0){
			return errors.New("JIF had bad read position 1 " + strconv.FormatInt(term1, 10));
		}
		term1 = this.GetValueAtRegister(term1);
	}

	term2 := this.GetValueAtRegister(term2Register);
	if(instruction.SecondParameterMode == ParameterModePosition){
		if(term2 >= int64(len(this.Registers)) || term2 < 0){
			return errors.New("JIF had bad read position 2 " + strconv.FormatInt(term2, 10));
		}
		term2 = this.GetValueAtRegister(term2);
	}

	if(term1 == 0){
		this.InstructionPointer = term2;
	} else{
		this.InstructionPointer += int64(instruction.OperationLength);
	}

	return nil;
}

func (this *IntcodeMachineV2) ExecuteInput(instruction *IntcodeInstruction) error {

	if(instruction.FistParameterMode != ParameterModePosition){
		return errors.New("Input instructions can't write to immediate mode registers");
	}

	term1Register := this.InstructionPointer+1;
	if(term1Register >= int64(len(this.Registers))){
		return errors.New("Input instruction went outside bounds " + strconv.FormatInt(this.InstructionPointer, 10));
	}
	destPosition := this.GetValueAtRegister(term1Register);

	if(len(this.InputQueue) <= 0){
		return errors.New("Input instruction had no pending input");
	}

	inputVal := this.InputQueue[0];
	this.InputQueue = this.InputQueue[1:];

	Log.Info("[INPUT] - Proccessed Input " + strconv.FormatInt(inputVal, 10) + " - stored in " + strconv.FormatInt(destPosition, 10));
	this.SetValueAtRegister(destPosition, inputVal);
	this.InstructionPointer += int64(instruction.OperationLength);
	return nil;
}

func (this *IntcodeMachineV2) ExecuteMul(instruction *IntcodeInstruction) error {
	if(instruction.ThirdParameterMode != ParameterModePosition){
		return errors.New("Mul instructions can't write to immediate mode registers");
	}

	term1Register := this.InstructionPointer+1;
	term2Register := this.InstructionPointer+2;
	destRegister := this.InstructionPointer+3;
	if(term1Register >= int64(len(this.Registers)) || term2Register >= int64(len(this.Registers)) || destRegister >= int64(len(this.Registers))){
		return errors.New("Mul instruction went outside bounds " + strconv.FormatInt(this.InstructionPointer, 10));
	}

	term1 := this.GetValueAtRegister(term1Register);
	if(instruction.FistParameterMode == ParameterModePosition){
		if(term1 >= int64(len(this.Registers)) || term1 < 0){
			return errors.New("Mul had bad read position 1 " + strconv.FormatInt(term1, 10));
		}
		term1 = this.GetValueAtRegister(term1);
	}

	term2 := this.GetValueAtRegister(term2Register);
	if(instruction.SecondParameterMode == ParameterModePosition){
		if(term2 >= int64(len(this.Registers)) || term2 < 0){
			return errors.New("Mul had bad read position 2 " + strconv.FormatInt(term2, 10));
		}
		term2 = this.GetValueAtRegister(term2);
	}



	destPosition := this.GetValueAtRegister(destRegister);

	if(destPosition >= int64(len(this.Registers)) || destPosition < 0){
		return errors.New("Mul had bad dest position " + strconv.FormatInt(destPosition, 10));
	}
	sum := term1 * term2;
	this.SetValueAtRegister(destPosition, sum);
	this.InstructionPointer += int64(instruction.OperationLength);
	return nil;
}

func (this *IntcodeMachineV2) GetValueAtRegister(register int64) int64 {
	return this.Registers[register];
}

func (this *IntcodeMachineV2) SetValueAtRegister(register int64, val int64)  {
	this.Registers[register] = val;
}

func (this *IntcodeMachineV2) PrintContents()  {
	buff := "";
	for i, val := range this.Registers{
		if(i > 0){
			buff += ", ";
		}
		buff += strconv.FormatInt(val, 10);
	}
	Log.Info("Machine Contents - " + buff)
}
