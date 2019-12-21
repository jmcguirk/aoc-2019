package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const SpringCodeOpAnd = "AND";
const SpringCodeOpOr = "OR";
const SpringCodeOpNot = "NOT";

type SpringCodeInstruction interface {
	ToString() string;
	Parse(line string, lineNum int) error;
	SerializeToInput(processor *SpringCodeProcessor);
	GetLineNum() int;
}


type SpringCodeAND struct{
	Register1 int;
	Register2 int;
	LineNum int;
}

func (this *SpringCodeAND) Parse(lineRaw string, lineNum int) error {
	line := strings.TrimSpace(lineRaw);
	lineParts := strings.Split(line, " ");
	if(len(lineParts) < 3){
		return errors.New("Unexpected length " + strconv.Itoa(len(lineParts)));
	}
	if(len(lineParts[1]) != 1){
		return errors.New("Unexpected arg 1 length " + strconv.Itoa(len(lineParts[1])));
	}
	if(len(lineParts[2]) != 1){
		return errors.New("Unexpected arg 2 length " + strconv.Itoa(len(lineParts[2])));
	}
	this.LineNum = lineNum;
	this.Register1 = int(lineParts[1][0]);
	this.Register2 = int(lineParts[2][0]);
	return nil;
}


func (this *SpringCodeAND) ToString() string {
	return fmt.Sprintf("%d - %s %c %c", this.LineNum, SpringCodeOpAnd, this.Register1, this.Register2);
}

func (this *SpringCodeAND) GetLineNum() int {
	return this.LineNum;
}

func (this *SpringCodeAND) SerializeToInput(processor *SpringCodeProcessor) {
	processor.EnqueueRawInput(AsciiA);
	processor.EnqueueRawInput(AsciiN);
	processor.EnqueueRawInput(AsciiD);
	processor.EnqueueRawInput(AsciiSpace);
	processor.EnqueueRawInput(this.Register1);
	processor.EnqueueRawInput(AsciiSpace);
	processor.EnqueueRawInput(this.Register2);
	processor.EnqueueRawInput(AsciiNewLine);
}







type SpringCodeOR struct{
	Register1 int;
	Register2 int;
	LineNum int;
}

func (this *SpringCodeOR) ToString() string {
	return fmt.Sprintf("%d - %s %c %c", this.LineNum, SpringCodeOpOr, this.Register1, this.Register2);
}

func (this *SpringCodeOR) GetLineNum() int {
	return this.LineNum;
}

func (this *SpringCodeOR) SerializeToInput(processor *SpringCodeProcessor) {
	processor.EnqueueRawInput(AsciiO);
	processor.EnqueueRawInput(AsciiR);
	processor.EnqueueRawInput(AsciiSpace);
	processor.EnqueueRawInput(this.Register1);
	processor.EnqueueRawInput(AsciiSpace);
	processor.EnqueueRawInput(this.Register2);
	processor.EnqueueRawInput(AsciiNewLine);
}

func (this *SpringCodeOR) Parse(lineRaw string, lineNum int) error {
	line := strings.TrimSpace(lineRaw);
	lineParts := strings.Split(line, " ");
	if(len(lineParts) < 3){
		return errors.New("Unexpected length " + strconv.Itoa(len(lineParts)));
	}
	if(len(lineParts[1]) != 1){
		return errors.New("Unexpected arg 1 length " + strconv.Itoa(len(lineParts[1])));
	}
	if(len(lineParts[2]) != 1){
		return errors.New("Unexpected arg 2 length " + strconv.Itoa(len(lineParts[2])));
	}
	this.LineNum = lineNum;
	this.Register1 = int(lineParts[1][0]);
	this.Register2 = int(lineParts[2][0]);
	return nil;
}




type SpringCodeNOT struct{
	Register1 int;
	Register2 int;
	LineNum int;
}

func (this *SpringCodeNOT) Parse(lineRaw string, lineNum int) error {
	line := strings.TrimSpace(lineRaw);
	lineParts := strings.Split(line, " ");
	if(len(lineParts) < 3){
		return errors.New("Unexpected length " + strconv.Itoa(len(lineParts)));
	}
	if(len(lineParts[1]) != 1){
		return errors.New("Unexpected arg 1 length " + strconv.Itoa(len(lineParts[1])));
	}
	if(len(lineParts[2]) != 1){
		return errors.New("Unexpected arg 2 length " + strconv.Itoa(len(lineParts[2])));
	}
	this.LineNum = lineNum;
	this.Register1 = int(lineParts[1][0]);
	this.Register2 = int(lineParts[2][0]);
	return nil;
}

func (this *SpringCodeNOT) GetLineNum() int {
	return this.LineNum;
}


func (this *SpringCodeNOT) ToString() string {
	return fmt.Sprintf("%d - %s %c %c", this.LineNum, SpringCodeOpNot, this.Register1, this.Register2);
}


func (this *SpringCodeNOT) SerializeToInput(processor *SpringCodeProcessor) {
	processor.EnqueueRawInput(AsciiN);
	processor.EnqueueRawInput(AsciiO);
	processor.EnqueueRawInput(AsciiT);
	processor.EnqueueRawInput(AsciiSpace);
	processor.EnqueueRawInput(this.Register1);
	processor.EnqueueRawInput(AsciiSpace);
	processor.EnqueueRawInput(this.Register2);
	processor.EnqueueRawInput(AsciiNewLine);
}