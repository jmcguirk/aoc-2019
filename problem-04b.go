package main

import (
	"math/big"
	"os"
	"strconv"
)

type Problem4B struct {

	Registers []int;
	MaxRegisters []int;
	MinRegisters []int;
	Matches int;
}

func (this *Problem4B) Solve() {
	Log.Info("Problem 4B solver beginning!")


	this.Registers = make([]int, 6);
	this.MinRegisters = make([]int, 6);
	this.MaxRegisters = make([]int, 6);


	min := big.NewInt(138307);
	max := big.NewInt(654504);

	this.LoadIntoBuff(this.Registers, min);
	this.LoadIntoBuff(this.MinRegisters, min);
	this.LoadIntoBuff(this.MaxRegisters, max);


	this.LoopOdometer(0, false);
}


func (this *Problem4B) LogAndExit() {
	Log.Info("Completed odometer loop - found %d matches ", this.Matches);
	os.Exit(0);
}



func (this *Problem4B) LoopOdometer(index int, parentHasDupe bool) {

	min := this.Registers[index];
	if(index > 0){
		min = this.Registers[index-1];

	}
	max := 9;


	for i := min; i <= max; i++ {

		weHaveDupe := false;
		if(!parentHasDupe){
			if(index > 0){
				if(i == this.Registers[index - 1]){
					weHaveDupe = true;
				}
			}
		}

		this.Registers[index] = i;

		if(index < len(this.Registers) - 1){
			this.LoopOdometer(index+1, parentHasDupe || weHaveDupe);
		} else{
			if(IsGTE(this.Registers, this.MaxRegisters)){
				this.LogAndExit();
			}
			if(parentHasDupe || weHaveDupe){
				if(IsGTE(this.Registers, this.MinRegisters) && this.ContainsPreciseSequence(this.Registers)){
					this.Matches++;
					this.LogBuff(this.Registers);
				}
			}

		}
	}
}

func (this *Problem4B) ContainsPreciseSequence(registers []int) bool {

	seqCount := 0;
	lastVal := -1;
	for _, v := range registers{
		if(lastVal < -1){
			lastVal = v;
			seqCount = 1;
		} else {
			if(lastVal == v){
				seqCount++;
			} else{
				if(seqCount == 2){
					return true;
				}
				seqCount = 1;
				lastVal = v;
			}
		}
	}
	if(seqCount == 2){
		return true;
	}
	return false;
}

func (this *Problem4B) LoadIntoBuff(registers []int, bigInt *big.Int) {

	Log.Info("Loading %d " , bigInt.Int64());
	for i, _ := range registers{
		cpy := big.NewInt(0);
		cpy.Set(bigInt)
		registers[i] = nthDigit(cpy,int64(len(registers) - i - 1));
	}
}

func (this *Problem4B) LogRegisters() {
	this.LogBuff(this.Registers);
}

func (this *Problem4B) IsEQStr(val string) bool {
	buff := "";
	for _, v := range this.Registers{
		buff += strconv.Itoa(v);
	}
	if(buff == val) {
		return true;
	}
	return false;
}





func (this *Problem4B) LogBuff(registers[]int) {

	buff := "";
	for _, v := range registers{
		buff += strconv.Itoa(v);
	}
	Log.Info(buff);
}
