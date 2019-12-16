
package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type FrequencySystem struct {
	Digits []int;
	CurrentDigits []int;
	BasePattern []int;
	UnpackedPatterns [][]int;
	FileName string;
	PhaseCount int;
	MessageOffset int;
	Log bool;
}

func (this *FrequencySystem) Parse(fileName string) error{
	this.Digits = make([]int, 0);
	this.BasePattern = make([]int, 4);
	this.BasePattern[0] = 0;
	this.BasePattern[1] = 1;
	this.BasePattern[2] = 0;
	this.BasePattern[3] = -1;
	fileContents, err := ioutil.ReadFile(fileName);
	if(err != nil){
		return err;
	}
	digitsRaw := strings.TrimSpace(string(fileContents));
	offsetBuff := "";
	for i, char := range digitsRaw{
		digit, err := strconv.ParseInt(string(char), 10, 64);
		if(err != nil){
			return err;
		}
		if(i < 7){
			offsetBuff += string(char);
		}
		this.Digits = append(this.Digits, int(digit));
		this.CurrentDigits = append(this.CurrentDigits, int(digit));
	}
	offsetVal, err := strconv.ParseInt(offsetBuff, 10, 64);
	if(err != nil){
		return err;
	}
	this.MessageOffset = int(offsetVal);
	Log.Info("message offset at %d", this.MessageOffset);


	/*
	this.UnpackedPatterns = make([][]int, len(this.CurrentDigits));
	for i, _ := range this.UnpackedPatterns{
		cpyCount := i + 1;
		unpacked := make([]int, cpyCount * len(this.BasePattern));
		k := 0;
		for j := range this.BasePattern{
			for cpy := 0; cpy < cpyCount; cpy++{
				unpacked[k] = this.BasePattern[j];
				k++;
			}
		}
		this.UnpackedPatterns[i] = unpacked;
		//Log.Info("Unpacked %d %d ", i, len(unpacked));
	}*/

	Log.Info("Completed parsing frequency system, %d digits loaded", len(this.Digits));
	return nil;
}

func (this *FrequencySystem) GetPatternValue(digitIndex int, iterationCount int) int{
	/*
	lookup := this.UnpackedPatterns[iterationCount];
	return lookup[(digitIndex+ 1) % len(lookup)];
	*/
	return this.BasePattern[((digitIndex+1)/(iterationCount+1))%len(this.BasePattern)];
}

func (this *FrequencySystem) Step(count int) error{

	debugBuff := "";
	for n := 0; n < count; n++{
		newDigits := make([]int, len(this.CurrentDigits));

		for i, _ := range newDigits{
			sum := 0;
			if(this.Log){
				debugBuff = "";
			}

			for j, v := range this.CurrentDigits{
				if(this.Log){
					if(debugBuff != ""){
						debugBuff += "  +  ";
					}
				}
				val := this.GetPatternValue(j, i);
				sum += v * val;
				if(this.Log){
					debugBuff += fmt.Sprintf("%d*%d", v, val);
				}

			}
			if(sum < 0){
				sum = sum * -1;
			}
			if(sum >= 10){
				sum = sum % 10;
			}
			newDigits[i] = sum;
			if(this.Log) {
				//debugBuff += fmt.Sprintf(" = %d", sum);
				//Log.Info(debugBuff);
			}
		}

		this.CurrentDigits = newDigits;
		this.PhaseCount++;
		if(this.Log){
			Log.Info(this.ToString());
		}

	}

	return nil;
}


func (this *FrequencySystem) StepMulti(multiple int, count int) error{

	copied := make([]int, 0);
	for i:= 0; i < multiple; i++{
		copied = append(copied, this.CurrentDigits...);
	}

	// Copy just the upper half - we don't care about simulating the signal before this offset
	newDigits := make([]int, 0);
	for i:= this.MessageOffset; i < len(copied); i++{
		newDigits = append(newDigits, (copied[i]));
	}
	Log.Info("New array is %d digits - %d" , len(newDigits), len(this.CurrentDigits) * multiple);


	for n := 0; n < count; n++{
		sum := 0;

		for i := (len(newDigits) - 1); i >= 0; i-- {
			sum += newDigits[i];
			newDigits[i] = sum % 10;
		}

		this.PhaseCount++;
	}


	this.CurrentDigits = newDigits;


	return nil;
}

func (this *FrequencySystem) ToString() string{
	buff := "";

	digits := "";
	for _, digit := range this.CurrentDigits{
		digits += strconv.Itoa(digit);
	}

	buff += fmt.Sprintf("After %d phase: %s\n", this.PhaseCount, digits);
	return buff;
}

func (this *FrequencySystem) ToShortString() string{
	buff := "";

	digits := "";
	for i := 0; i < 8; i++{
		digits += strconv.Itoa(this.CurrentDigits[i]);;
	}

	buff += fmt.Sprintf("After %d phase: %s\n", this.PhaseCount, digits);
	return buff;
}