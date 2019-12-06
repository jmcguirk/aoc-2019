package main;

import (
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"time"
)

func FormatDuration(duration time.Duration) string{
	return FormatDurationMS(int64(duration.Seconds() * 1000));
}

func FormatDurationMS(durationMS int64) string{
	if(durationMS < 1000){
		return strconv.FormatInt(durationMS, 10) + "ms";
	}
	secs := float64(durationMS) / 1000;
	if(secs < 300){
		return fmt.Sprintf("%.2fs", secs);
	}
	mins := secs / 60;
	if(mins < 60){
		return fmt.Sprintf("%.2fm", mins);
	}
	hr := mins / 60;
	return fmt.Sprintf("%.2fh", hr);
}


type IntVec2 struct{
	X 		int;
	Y		int;
}


func (this *IntVec2) ManhattanDistance(other *IntVec2) int{
	xComp := this.X - other.X;
	if(xComp < 0){
		xComp *= -1;
	}
	yComp := this.Y - other.Y;
	if(yComp < 0){
		yComp *= -1;
	}
	return xComp + yComp;
}

func nthDigit(bigI *big.Int, n int64) int {
	var quotient big.Int
	quotient.Exp(big.NewInt(10), big.NewInt(n), nil)

	bigI.Div(bigI, &quotient)

	var result big.Int
	result.Mod(bigI, big.NewInt(10))

	return int(result.Int64());
}

func nthDigit64(val int64, n int64) int {
	var quotient big.Int
	quotient.Exp(big.NewInt(10), big.NewInt(n), nil)

	bigI := big.NewInt(val);
	bigI.Div(bigI, &quotient)

	var result big.Int
	result.Mod(bigI, big.NewInt(10))

	return int(result.Int64());
}


func IsGTEOrEqual(registersA []int, registersB []int) bool {
	for i, v := range registersA{
		if(v > registersB[i]){
			return true;
		}
		if(v < registersB[i]){
			return false;
		}
	}
	return true;
}

func ReverseSlice(s interface{}) {
	size := reflect.ValueOf(s).Len()
	swap := reflect.Swapper(s)
	for i, j := 0, size-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
}

func IsGTE(registersA []int, registersB []int) bool {
	for i, v := range registersA{
		if(v > registersB[i]){
			return true;
		}
		if(v < registersB[i]){
			return false;
		}
	}
	return false;
}

func IsEQ(registersA []int, registersB []int) bool {
	for i, v := range registersA{
		if(v != registersB[i]){
			return false;
		}
	}
	return true;
}