package main

const IntCodeOpCodeAdd = 1;
const IntCodeOpCodeMul = 2;
const IntCodeOpCodeHalt = 99;
const IntCodeOpCodeInput = 3;
const IntCodeOpCodeOutput = 4;
const IntCodeOpCodeJumpIfTrue = 5;
const IntCodeOpCodeJumpIfFalse = 6;
const IntCodeOpCodeLessThan = 7;
const IntCodeOpCodeEquals = 8;
const IntCodeOpCodeAdjustRelativeOffset = 9;


const ParameterModePosition = 0;
const ParameterModeImmediate = 1;
const ParameterModeRelative = 2;
