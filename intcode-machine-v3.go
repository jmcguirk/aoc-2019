package main

import (
	"errors"
	"io/ioutil"
	"math/big"
	"strconv"
	"strings"
)

type IntcodeMachineV3 struct {
	Registers []*big.Int;
	InitialState []*big.Int;
	InstructionPointer int64;
	InstructionsExecuted int;
	InputQueue []int64;
	LastOutput int;
	LastOutputValue *big.Int;
	HasHalted bool;
	PauseOnOutput bool;
	RelativeBase int64;
	PauseOnInput bool;
	PendingInput bool;
	HasDefaultInput bool;
	DefaultInputValue int64;
	PauseOnDefaultInput bool;
}






func (this *IntcodeMachineV3) QueueInput(val int64) {
	//Log.Info("Queuing input value %d", val);
	this.InputQueue = append(this.InputQueue, val);
}

func (this *IntcodeMachineV3) SetDefaultInput(val int64) {
	this.DefaultInputValue = val;
	this.HasDefaultInput = true;
}


func (this *IntcodeMachineV3) Load(fileName string) error {
	//Log.Info("Loading intcode v3 machine from %s", fileName)
	this.Registers = make([]*big.Int, 0);
	this.InputQueue = make([]int64, 0);
	this.RelativeBase = 0;
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
			this.Registers = append(this.Registers, big.NewInt(parsed));
		}
	}

	// Reset the function pointer to start at the beginning
	this.InstructionPointer = 0;
	this.LastOutputValue = nil;

	this.InitialState = make([]*big.Int, len(this.Registers));
	copy(this.InitialState, this.Registers);
	this.InputQueue = make([]int64, 0);


	//Log.Info("Finished parsing machine initial state - %d registers", len(this.Registers));

	return nil;
}

func (this *IntcodeMachineV3) ReadNextOutput() (int64, error, bool) {
	err := this.Execute();
	if(err != nil){
		return -1, err, false;
	}
	if(this.HasHalted){
		return -1, nil, true;
	}
	if(this.LastOutputValue == nil){
		return -1, nil, false;
	}
	return this.LastOutputValue.Int64(), nil, false;
}


func (this *IntcodeMachineV3) Reset() {
	//Log.Info("Performing reset");
	this.InstructionPointer = 0;
	this.LastOutputValue = nil;
	this.InputQueue = make([]int64, 0);
	this.LastOutput = 0;
	this.HasHalted = false;
	this.HasHalted = false;
	this.RelativeBase = 0;
	copy(this.Registers, this.InitialState);
}

func (this *IntcodeMachineV3) Execute() error {
	for {
		if(this.InstructionPointer < 0 || this.InstructionPointer > int64(len(this.Registers))){
			return errors.New("Instruction pointer went outside bounds " + strconv.FormatInt(this.InstructionPointer, 10));
		}

		instruction, err := this.ParseOperation();
		if(err != nil){
			return err;
		}


		if(this.InstructionsExecuted % 10000 == 0){
			//Log.Info("[EXEC - " + strconv.Itoa(this.InstructionsExecuted) + "] " + instruction.Describe());
		}
		//Log.Info("[EXEC - " + strconv.Itoa(this.InstructionsExecuted) + "] " + instruction.Describe());
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
			if(this.PauseOnInput && len(this.InputQueue) <= 0){
				//Log.Info("Pausing for input");
				this.PendingInput = true;
				return nil;
			}
			breakPostProcess := false;
			if(this.PauseOnDefaultInput && this.HasDefaultInput && len(this.InputQueue) <= 0){
				breakPostProcess = true;
			}
			this.PendingInput = false;
			err := this.ExecuteInput(instruction);
			if(err != nil){
				return err;
			}
			if(breakPostProcess){
				return nil;
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
		case IntCodeOpCodeAdjustRelativeOffset:
			err := this.ExecuteRelativeOffset(instruction);
			if(err != nil){
				return err;
			}
			break;
		default:
			return errors.New("Unimplemented opcode " + instruction.Describe());
			break;
		}

		this.InstructionsExecuted++;
	}
}

func (this *IntcodeMachineV3) ExpandRegisters(destPos int64) {
	for {
		if(int64(len(this.Registers)) > destPos){
			break;
		}
		//Log.Info("expanding registers to %d", len(this.Registers) * 2);
		this.Registers = append(this.Registers, make([]*big.Int, len(this.Registers))...);
	}

}

func (this *IntcodeMachineV3) ParseOperation() (*IntcodeInstruction, error){
	res := &IntcodeInstruction{};
	rawVal := this.Registers[this.InstructionPointer];
	//Log.Info("Parsing instruction from %d", rawVal.Int64());
	op1 := nthDigit(rawVal, 0);
	op2 := nthDigit(rawVal, 1);
	res.Operation = op1 + (op2 * 10);
	res.FistParameterMode = nthDigit(rawVal, 2);
	res.SecondParameterMode = nthDigit(rawVal, 3);
	res.ThirdParameterMode = nthDigit(rawVal, 4);
	res.DeriveLength();

	return res, nil;
}

func (this *IntcodeMachineV3) ExecuteOutput(instruction *IntcodeInstruction) error {
	term1Register := this.InstructionPointer+1;
	if(term1Register >= int64(len(this.Registers))){
		return errors.New("Add instruction went outside bounds " + strconv.FormatInt(this.InstructionPointer, 10));
	}

	//readPos := term1Register;
	term1 := this.GetValueAtRegister(term1Register);
	if(instruction.FistParameterMode == ParameterModePosition){
		if( term1.Int64() < 0){
			return errors.New("Output read position 1 " + strconv.FormatInt(term1.Int64(), 10));
		}
		term1 = this.GetValueAtRegister(term1.Int64());
		//readPos = term1.Int64();
	} else if(instruction.FistParameterMode == ParameterModeRelative){
		offset := this.RelativeBase + term1.Int64();
		if(offset < 0){
			return errors.New("Output read position 1 " + strconv.FormatInt(offset, 10));
		}
		term1 = this.GetValueAtRegister(offset);
		//readPos = offset;
	}

	//Log.Info("[OUTPUT] - " + strconv.FormatInt(term1.Int64(), 10) + " from read pos " + strconv.FormatInt(readPos, 10));
	this.LastOutputValue = term1;
	this.InstructionPointer += int64(instruction.OperationLength);
	return nil;
}

func (this *IntcodeMachineV3) ExecuteAdd(instruction *IntcodeInstruction) error {

	if(instruction.ThirdParameterMode == ParameterModeImmediate){
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
		if(term1.Int64() < 0){
			return errors.New("Add had bad read position 1 " + strconv.FormatInt(term1.Int64(), 10));
		}
		term1 = this.GetValueAtRegister(term1.Int64());
	} else if(instruction.FistParameterMode == ParameterModeRelative){
		offset := this.RelativeBase + term1.Int64();
		if(offset < 0){
			return errors.New("Add read position 1 " + strconv.FormatInt(offset, 10));
		}
		term1 = this.GetValueAtRegister(offset);
	}

	term2 := this.GetValueAtRegister(term2Register);
	if(instruction.SecondParameterMode == ParameterModePosition){
		if(term2.Int64() < 0){
			return errors.New("Add had bad read position 2 " + strconv.FormatInt(term2.Int64(), 10));
		}
		term2 = this.GetValueAtRegister(term2.Int64());
	} else if(instruction.SecondParameterMode == ParameterModeRelative){
		offset := this.RelativeBase + term2.Int64();
		if(offset < 0){
			return errors.New("Output read position 1 " + strconv.FormatInt(offset, 10));
		}
		term2 = this.GetValueAtRegister(offset);
	}



	destPosition := this.GetValueAtRegister(destRegister).Int64();
	if(instruction.ThirdParameterMode == ParameterModeRelative){
		destPosition = this.RelativeBase + destPosition;
	}

	if(destPosition < 0){
		return errors.New("Add had bad dest position " + strconv.FormatInt(destPosition, 10));
	}
	if(destPosition >= int64(len(this.Registers))){
		this.ExpandRegisters(destPosition);
	}

	//if(destPosition >= int64(len(this.Registers)) || destPosition < 0){
		//return errors.New("Add had bad dest position " + strconv.FormatInt(destPosition, 10));
	//}

	bigI := new(big.Int);
	bigI.Set(term1);

	sum := bigI.Add(bigI, term2);
	this.SetValueAtRegister(destPosition, sum);

	//Log.Info("executed add %d stored at %d", sum.Int64(), destPosition);

	this.InstructionPointer += int64(instruction.OperationLength);
	return nil;
}

func (this *IntcodeMachineV3) ExecuteLessThan(instruction *IntcodeInstruction) error {

	if(instruction.ThirdParameterMode == ParameterModeImmediate){
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
		if(term1.Int64() < 0){
			return errors.New("LE had bad read position 1 " + strconv.FormatInt(term1.Int64(), 10));
		}
		term1 = this.GetValueAtRegister(term1.Int64());
	} else if(instruction.FistParameterMode == ParameterModeRelative){
		offset := this.RelativeBase + term1.Int64();
		if(offset < 0){
			return errors.New("LE read position 1 " + strconv.FormatInt(offset, 10));
		}
		term1 = this.GetValueAtRegister(offset);
	}

	term2 := this.GetValueAtRegister(term2Register);
	if(instruction.SecondParameterMode == ParameterModePosition){
		if(term2.Int64() < 0){
			return errors.New("LE had bad read position 2 " + strconv.FormatInt(term2.Int64(), 10));
		}
		term2 = this.GetValueAtRegister(term2.Int64());
	} else if(instruction.SecondParameterMode == ParameterModeRelative){
		offset := this.RelativeBase + term2.Int64();
		if(offset < 0){
			return errors.New("LE read position 2 " + strconv.FormatInt(offset, 10));
		}
		term2 = this.GetValueAtRegister(offset);
	}



	destPosition := this.GetValueAtRegister(destRegister).Int64();
	if(instruction.ThirdParameterMode == ParameterModeRelative){
		destPosition = this.RelativeBase + destPosition;
	}

	if(destPosition < 0){
		return errors.New("LE had bad dest position " + strconv.FormatInt(destPosition, 10));
	}
	if(destPosition >= int64(len(this.Registers))){
		this.ExpandRegisters(destPosition);
	}

	if(term1.Int64() < term2.Int64()){
		this.SetValueAtRegister(destPosition, big.NewInt(1));
	} else{
		this.SetValueAtRegister(destPosition,  big.NewInt(0));
	}


	this.InstructionPointer += int64(instruction.OperationLength);
	return nil;
}

func (this *IntcodeMachineV3) ExecuteEquals(instruction *IntcodeInstruction) error {

	if(instruction.ThirdParameterMode == ParameterModeImmediate){
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
		if(term1.Int64() < 0){
			return errors.New("EQ had bad read position 1 " + strconv.FormatInt(term1.Int64(), 10));
		}
		term1 = this.GetValueAtRegister(term1.Int64());
	} else if(instruction.FistParameterMode == ParameterModeRelative){
		offset := this.RelativeBase + term1.Int64();
		if(offset < 0){
			return errors.New("EQ read position 1 " + strconv.FormatInt(offset, 10));
		}
		term1 = this.GetValueAtRegister(offset);
	}


	term2 := this.GetValueAtRegister(term2Register);
	if(instruction.SecondParameterMode == ParameterModePosition){
		if(term2.Int64() < 0){
			return errors.New("EQ had bad read position 2 " + strconv.FormatInt(term2.Int64(), 10));
		}
		term2 = this.GetValueAtRegister(term2.Int64());
	} else if(instruction.SecondParameterMode == ParameterModeRelative){
		offset := this.RelativeBase + term2.Int64();
		if(offset < 0){
			return errors.New("EQ read position 2 " + strconv.FormatInt(offset, 10));
		}
		term2 = this.GetValueAtRegister(offset);
	}


	destPosition := this.GetValueAtRegister(destRegister).Int64();
	if(instruction.ThirdParameterMode == ParameterModeRelative){
		destPosition = this.RelativeBase + destPosition;
	}

	if(destPosition < 0){
		return errors.New("EQ had bad dest position " + strconv.FormatInt(destPosition, 10));
	}
	if(destPosition >= int64(len(this.Registers))){
		this.ExpandRegisters(destPosition);
	}
	if(term1.Int64() == term2.Int64()){
		this.SetValueAtRegister(destPosition, big.NewInt(1));
	} else{
		this.SetValueAtRegister(destPosition,  big.NewInt(0));
	}


	this.InstructionPointer += int64(instruction.OperationLength);
	return nil;
}

func (this *IntcodeMachineV3) ExecuteJumpIfTrue(instruction *IntcodeInstruction) error {
	term1Register := this.InstructionPointer+1;
	term2Register := this.InstructionPointer+2;
	if(term1Register >= int64(len(this.Registers)) || term2Register >= int64(len(this.Registers))){
		return errors.New("JIT instruction went outside bounds " + strconv.FormatInt(this.InstructionPointer, 10));
	}

	term1 := this.GetValueAtRegister(term1Register);
	if(instruction.FistParameterMode == ParameterModePosition){
		if(term1.Int64() < 0){
			return errors.New("JIT had bad read position 1 " + strconv.FormatInt(term1.Int64(), 10));
		}
		term1 = this.GetValueAtRegister(term1.Int64());
	} else if(instruction.FistParameterMode == ParameterModeRelative){
		offset := this.RelativeBase + term1.Int64();
		if(offset < 0){
			return errors.New("JIT read position 1 " + strconv.FormatInt(offset, 10));
		}
		term1 = this.GetValueAtRegister(offset);
	}

	term2 := this.GetValueAtRegister(term2Register);
	if(instruction.SecondParameterMode == ParameterModePosition){
		if(term2.Int64() < 0){
			return errors.New("JIT had bad read position 2 " + strconv.FormatInt(term2.Int64(), 10));
		}
		term2 = this.GetValueAtRegister(term2.Int64());
	} else if(instruction.SecondParameterMode == ParameterModeRelative){
		offset := this.RelativeBase + term2.Int64();
		if(offset < 0){
			return errors.New("JIT read position 1 " + strconv.FormatInt(offset, 10));
		}
		term2 = this.GetValueAtRegister(offset);
	}

	if(term1.Int64() != 0){
		this.InstructionPointer = term2.Int64();
	} else{
		this.InstructionPointer += int64(instruction.OperationLength);
	}

	return nil;
}

func (this *IntcodeMachineV3) ExecuteJumpIfFalse(instruction *IntcodeInstruction) error {
	term1Register := this.InstructionPointer+1;
	term2Register := this.InstructionPointer+2;
	if(term1Register >= int64(len(this.Registers)) || term2Register >= int64(len(this.Registers))){
		return errors.New("JIF instruction went outside bounds " + strconv.FormatInt(this.InstructionPointer, 10));
	}

	term1 := this.GetValueAtRegister(term1Register);
	if(instruction.FistParameterMode == ParameterModePosition){
		if(term1.Int64() < 0){
			return errors.New("JIF had bad read position 1 " + strconv.FormatInt(term1.Int64(), 10));
		}
		term1 = this.GetValueAtRegister(term1.Int64());
	} else if(instruction.FistParameterMode == ParameterModeRelative){
		offset := this.RelativeBase + term1.Int64();
		if(offset < 0){
			return errors.New("JIF read position 1 " + strconv.FormatInt(offset, 10));
		}
		term1 = this.GetValueAtRegister(offset);
	}

	term2 := this.GetValueAtRegister(term2Register);
	if(instruction.SecondParameterMode == ParameterModePosition){
		if(term2.Int64() < 0){
			return errors.New("JIF had bad read position 2 " + strconv.FormatInt(term2.Int64(), 10));
		}
		term2 = this.GetValueAtRegister(term2.Int64());
	} else if(instruction.SecondParameterMode == ParameterModeRelative){
		offset := this.RelativeBase + term2.Int64();
		if(offset < 0){
			return errors.New("JIF read position 1 " + strconv.FormatInt(offset, 10));
		}
		term2 = this.GetValueAtRegister(offset);
	}

	if(term1.Int64() == 0){
		this.InstructionPointer = term2.Int64();
	} else{
		this.InstructionPointer += int64(instruction.OperationLength);
	}

	return nil;
}

func (this *IntcodeMachineV3) ExecuteRelativeOffset(instruction *IntcodeInstruction) error {

	term1Register := this.InstructionPointer+1;
	if(term1Register >= int64(len(this.Registers))){
		return errors.New("Add instruction went outside bounds " + strconv.FormatInt(this.InstructionPointer, 10));
	}

	term1 := this.GetValueAtRegister(term1Register);
	if(instruction.FistParameterMode == ParameterModePosition){
		if(term1.Int64() < 0){
			return errors.New("Add had bad read position 1 " + strconv.FormatInt(term1.Int64(), 10));
		}
		term1 = this.GetValueAtRegister(term1.Int64());
	} else if(instruction.FistParameterMode == ParameterModeRelative){
		offset := this.RelativeBase + term1.Int64();
		if(offset < 0){
			return errors.New("Add read position 1 " + strconv.FormatInt(offset, 10));
		}
		term1 = this.GetValueAtRegister(offset);
	}
	this.RelativeBase += term1.Int64();
	//Log.Info("adjusted relative base to %d - new value is %d", term1.Int64(), this.RelativeBase);

	this.InstructionPointer += int64(instruction.OperationLength);
	return nil;
}



func (this *IntcodeMachineV3) ExecuteInput(instruction *IntcodeInstruction) error {

	if(instruction.FistParameterMode == ParameterModeImmediate){
		return errors.New("Input instructions can't write to immediate mode registers");
	}

	term1Register := this.InstructionPointer+1;
	if(term1Register >= int64(len(this.Registers))){
		return errors.New("Input instruction went outside bounds " + strconv.FormatInt(this.InstructionPointer, 10));
	}
	destPosition := this.GetValueAtRegister(term1Register).Int64();
	if(instruction.FistParameterMode == ParameterModeRelative){
		destPosition = this.RelativeBase + destPosition;
	}


	if(destPosition < 0){
		return errors.New("Input had bad dest position " + strconv.FormatInt(destPosition, 10));
	}
	if(destPosition >= int64(len(this.Registers))){
		this.ExpandRegisters(destPosition);
	}


	var inputVal int64;
	if(len(this.InputQueue) <= 0 && !this.HasDefaultInput){
		return errors.New("Input instruction had no pending input");
	}


	if(len(this.InputQueue) > 0){
		inputVal = this.InputQueue[0];
		this.InputQueue = this.InputQueue[1:];
	} else{
		inputVal = this.DefaultInputValue;
	}


	//Log.Info("[INPUT] - Proccessed Input " + strconv.FormatInt(inputVal, 10) + " - stored in " + strconv.FormatInt(destPosition, 10));
	this.SetValueAtRegister(destPosition, big.NewInt(inputVal));
	this.InstructionPointer += int64(instruction.OperationLength);
	return nil;
}

func (this *IntcodeMachineV3) ExecuteMul(instruction *IntcodeInstruction) error {
	if(instruction.ThirdParameterMode == ParameterModeImmediate){
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
		if(term1.Int64() < 0){
			return errors.New("Mul had bad read position 1 " + strconv.FormatInt(term1.Int64(), 10));
		}
		term1 = this.GetValueAtRegister(term1.Int64());
	} else if(instruction.FistParameterMode == ParameterModeRelative){
		offset := this.RelativeBase + term1.Int64();
		if(offset < 0){
			return errors.New("MUL read position 1 " + strconv.FormatInt(offset, 10));
		}
		term1 = this.GetValueAtRegister(offset);
	}

	term2 := this.GetValueAtRegister(term2Register);
	if(instruction.SecondParameterMode == ParameterModePosition){
		if(term2.Int64() < 0){
			return errors.New("Mul had bad read position 2 " + strconv.FormatInt(term2.Int64(), 10));
		}
		term2 = this.GetValueAtRegister(term2.Int64());
	} else if(instruction.SecondParameterMode == ParameterModeRelative){
		offset := this.RelativeBase + term2.Int64();
		if(offset < 0){
			return errors.New("Mul read position 1 " + strconv.FormatInt(offset, 10));
		}
		term2 = this.GetValueAtRegister(offset);
	}



	destPosition := this.GetValueAtRegister(destRegister).Int64();
	if(instruction.ThirdParameterMode == ParameterModeRelative){
		destPosition = this.RelativeBase + destPosition;
	}
	if(destPosition < 0){
		return errors.New("Mul had bad dest position " + strconv.FormatInt(destPosition, 10));
	}
	if(destPosition >= int64(len(this.Registers))){
		this.ExpandRegisters(destPosition);
	}
	bigI := new(big.Int);
	bigI.Set(term1);
	sum := bigI.Mul(bigI, term2);
	this.SetValueAtRegister(destPosition, sum);
	this.InstructionPointer += int64(instruction.OperationLength);
	return nil;
}

func (this *IntcodeMachineV3) GetValueAtRegister(register int64) *big.Int {
	if(register >= int64(len(this.Registers))){
		return big.NewInt(0);
	}
	cpy := big.NewInt(0);
	existing := this.Registers[register];
	if(existing != nil){
		cpy.Set(existing);
	}
	return cpy;
}

func (this *IntcodeMachineV3) SetValueAtRegister(register int64, val *big.Int)  {
	cpy := big.NewInt(0);
	cpy.Set(val);
	this.Registers[register] = cpy;
}

func (this *IntcodeMachineV3) PrintContents()  {
	buff := "";
	for i, val := range this.Registers{
		if(i > 0){
			buff += ", ";
		}
		buff += strconv.FormatInt(val.Int64(), 10);
	}
	Log.Info("Machine Contents - " + buff)
}
