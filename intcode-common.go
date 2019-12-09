package main

type IntcodeInstruction struct{
	Operation int;
	FistParameterMode int;
	SecondParameterMode int;
	ThirdParameterMode int;
	OperationLength int;
}

func (this *IntcodeInstruction) DeriveLength() {
	switch(this.Operation){
	case IntCodeOpCodeAdd:
		this.OperationLength = 4;
		break;
	case IntCodeOpCodeMul:
		this.OperationLength = 4;
		break;
	case IntCodeOpCodeHalt:
		this.OperationLength = 1;
		break;
	case IntCodeOpCodeInput:
		this.OperationLength = 2;
		break;
	case IntCodeOpCodeOutput:
		this.OperationLength = 2;
		break;
	case IntCodeOpCodeJumpIfTrue:
		this.OperationLength = 3;
		break;
	case IntCodeOpCodeJumpIfFalse:
		this.OperationLength = 3;
		break;
	case IntCodeOpCodeLessThan:
		this.OperationLength = 4;
		break;
	case IntCodeOpCodeEquals:
		this.OperationLength = 4;
	case IntCodeOpCodeAdjustRelativeOffset:
		this.OperationLength = 2;
		break;
	}
}


func (this *IntcodeInstruction) Describe() string {
	buff := "";
	switch(this.Operation){
	case IntCodeOpCodeAdd:
		buff += "ADD";
		break;
	case IntCodeOpCodeMul:
		buff += "MUL";
		break;
	case IntCodeOpCodeHalt:
		buff += "HALT";
		break;
	case IntCodeOpCodeInput:
		buff += "INPUT";
		break;
	case IntCodeOpCodeOutput:
		buff += "OUTPUT";
		break;
	case IntCodeOpCodeJumpIfTrue:
		buff += "JIT";
		break;
	case IntCodeOpCodeJumpIfFalse:
		buff += "JIF";
		break;
	case IntCodeOpCodeLessThan:
		buff += "LT";
		break;
	case IntCodeOpCodeEquals:
		buff += "EQ";
		break;
	case IntCodeOpCodeAdjustRelativeOffset:
		buff += "REL";
		break;
	}
	//buff += "(";
	//buff += strconv.Itoa(this.OperationLength);
	//buff += ")";
	switch this.FistParameterMode {
	case ParameterModePosition:
		buff += " POS";
		break;
	case ParameterModeImmediate:
		buff += " IMM";
		break;
	case ParameterModeRelative:
		buff += " REL";
		break;
	}
	if(this.OperationLength > 2){
		switch this.SecondParameterMode {
		case ParameterModePosition:
			buff += " POS";
			break;
		case ParameterModeImmediate:
			buff += " IMM";
			break;
		case ParameterModeRelative:
			buff += " REL";
			break;
		}
	}
	if(this.OperationLength > 3){
		switch this.ThirdParameterMode {
		case ParameterModePosition:
			buff += " POS";
			break;
		case ParameterModeImmediate:
			buff += " IMM";
			break;
		case ParameterModeRelative:
			buff += " REL";
			break;
		}
	}
	return buff;
}

