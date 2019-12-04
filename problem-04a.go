package main

import (
	"math/big"
	"os"
	"strconv"
)

type Problem4A struct {

	Registers []int;
	MaxRegisters []int;
	MinRegisters []int;
	Matches int;
}

func (this *Problem4A) Solve() {
	Log.Info("Problem 4A solver beginning!")


	this.Registers = make([]int, 6);
	this.MaxRegisters = make([]int, 6);


	min := big.NewInt(138307);
	max := big.NewInt(654504);

	this.LoadIntoBuff(this.Registers, min);
	this.LoadIntoBuff(this.MaxRegisters, max);


	this.LoopOdometer(0, true, false);
}


func (this *Problem4A) LogAndExit() {
	Log.Info("Completed odometer loop - found %s matches ", this.Matches);
	os.Exit(0);
}

func (this *Problem4A) LoopOdometer(index int, atMax bool, hasDupe bool) {

	min := this.Registers[index];
	if(index > 0){
		min = this.Registers[index-1];

	}
	max := 9;


	for i := min; i <= max; i++ {
		if(!hasDupe){
			if(index > 0){
				if(i == this.Registers[index - 1]){
					hasDupe = true;
				}
			}
		}
		if(atMax){
			if(i < this.MaxRegisters[i]){
				atMax = false;
			}
		}


		this.Registers[index] = i;

		if(index < len(this.Registers) - 1){
			this.LoopOdometer(index+1, atMax, hasDupe);
		} else{
			if(hasDupe){
				this.Matches++;
				this.LogBuff(this.Registers);
			}
			if(atMax){
				this.LogAndExit();
			}

		}
	}
}


func (this *Problem4A) LoadIntoBuff(registers []int, bigInt *big.Int) {

	Log.Info("Loading %d " , bigInt.Int64());
	for i, _ := range registers{
		cpy := big.NewInt(0);
		cpy.Set(bigInt)
		registers[i] = nthDigit(cpy,int64(len(registers) - i - 1));
	}
}

func (this *Problem4A) LogRegisters() {
	this.LogBuff(this.Registers);
}


func (this *Problem4A) LogBuff(registers[]int) {

	buff := "";
	for _, v := range registers{
		buff += strconv.Itoa(v);
	}
	Log.Info(buff);
}
