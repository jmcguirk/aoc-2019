package main;

import (
	"fmt"
	"math"
	"math/big"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const slopeEpsilon = 0.00001;
const distEpsilon = 0.00001;

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

func (this *IntVec2) Slope(other *IntVec2) float32{
	if(other.X == this.X){
		return math.MaxFloat32;
	}
	return float32((other.Y - this.Y)) / float32((other.X - this.X));
}

func (this *IntVec2) Distance(other *IntVec2) float32{
	distX := (this.X - other.X);
	distY := (this.Y - other.Y);
	return float32(math.Sqrt(float64((distX*distX) + (distY*distY))));
}

func (this *IntVec2) Angle (other *IntVec2) float32{
	return float32(math.Atan2(float64(other.Y - this.Y), float64(other.X - this.X)));
	//atan2(y2 - y1, x2 - x1) * 180 / PI;
}

func (this *IntVec2) GetVisiblePoints(points []*IntVec2) []*IntVec2{
	res := make([]*IntVec2, 0);
	candidate := this;
	for _, neighbor := range points {
		if(neighbor == candidate){
			continue;
		}
		isOccluded := false;
		slopeN := candidate.Slope(neighbor);
		distN := candidate.Distance(neighbor);
		for _, occluder := range points {
			if(occluder == neighbor || occluder == candidate){
				continue;
			}
			slopeO := candidate.Slope(occluder);
			if(math.Abs(float64(slopeN - slopeO)) <= slopeEpsilon){

				if(math.Abs(float64((candidate.Distance(occluder) + neighbor.Distance(occluder)) - distN)) <= distEpsilon){
					isOccluded = true;
				}
			}
		}
		if(!isOccluded){
			res = append(res, neighbor);
		}
	}
	return res;
}

func Filter(target *IntVec2, points []*IntVec2) []*IntVec2{
	res := make([]*IntVec2, 0);
	for _, candidate := range points {
		if(candidate.X != target.X || candidate.Y != target.Y){
			res = append(res, candidate);
		}
	}
	return res;
}



func nthDigit(input *big.Int, n int64) int {
	var quotient big.Int
	quotient.Exp(big.NewInt(10), big.NewInt(n), nil)

	bigI := new(big.Int);
	bigI.Set(input);

	bigI.Div(bigI, &quotient)

	var result big.Int
	result.Mod(bigI, big.NewInt(10))

	return int(result.Int64());
}

// Perm calls f with each permutation of a.
func Perm(a []int64, f func([]int64)) {
	perm(a, f, 0)
}

// Permute the values at index i to len(a)-1.
func perm(a []int64, f func([]int64), i int) {
	if i > len(a) {
		f(a)
		return
	}
	perm(a, f, i+1)
	for j := i + 1; j < len(a); j++ {
		a[i], a[j] = a[j], a[i]
		perm(a, f, i+1)
		a[i], a[j] = a[j], a[i]
	}
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


func UpperAlphaCharacters() string {
	p := make([]byte, 26)
	for i := range p {
		p[i] = 'a' + byte(i)
	}
	return strings.ToUpper(string(p));
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