package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type PlanetarySystem struct {
	Bodies []*PlanetaryBody;
	LastBodyId int;
	FileName string;
	TimeIndex int;
	HasRestedX bool;
	XPeriod 	int;
	HasRestedY bool;
	YPeriod 	int;
	HasRestedZ bool;
	ZPeriod		int;
}

type PlanetaryBody struct {
	Position *IntVec3;
	Velocity *IntVec3;
	OriginalPosition *IntVec3;
	Id			int;
}



func (this *PlanetaryBody) TotalEnergy() int {
	return this.PotentialEnergy() * this.KineticEnergy()
}

func (this *PlanetaryBody) PotentialEnergy() int {
	return int(math.Abs(float64(this.Position.X)) + math.Abs(float64(this.Position.Y)) + math.Abs(float64(this.Position.Z)));
}

func (this *PlanetaryBody) PositionSignature() int {
	return this.Position.X * 1000000 + this.Position.Y * 1000 + this.Position.Z;
}

func (this *PlanetaryBody) Signature() string {
	return this.Print();
}



func (this *PlanetaryBody) KineticEnergy() int {
	return int(math.Abs(float64(this.Velocity.X)) + math.Abs(float64(this.Velocity.Y)) + math.Abs(float64(this.Velocity.Z)));
}

func (this *PlanetaryBody) ApplyVelocity() {
	this.Position.X += this.Velocity.X;
	this.Position.Y += this.Velocity.Y;
	this.Position.Z += this.Velocity.Z;
}

func (this *PlanetaryBody) Print() string {
	return fmt.Sprintf("pos=<x=%d,y=%d,z=%d>, vel=<x=%d, y=%d, z=%d>", this.Position.X, this.Position.Y, this.Position.Z, this.Velocity.X, this.Velocity.Y, this.Velocity.Z);
}


func (this *PlanetaryBody) IsAtRest() bool {
	return this.Velocity.X == 0 && this.Velocity.Y == 0 && this.Velocity.Z == 0;
}


func (this *PlanetaryBody) IsAtRestX() bool {
	return this.Velocity.X == 0;
}

func (this *PlanetaryBody) IsAtRestY() bool {
	return this.Velocity.Y == 0;
}

func (this *PlanetaryBody) IsAtRestZ() bool {
	return this.Velocity.Z == 0;
}



func (this *PlanetaryBody) IsAtOriginalPosition() bool {
	return this.Position.X == this.OriginalPosition.X && this.Position.Y == this.OriginalPosition.Y && this.Position.Z == this.OriginalPosition.Z;
}


func (this *PlanetarySystem) TotalEnergy() int {
	sum := 0;
	for _, body := range this.Bodies{
		sum += body.TotalEnergy();
	}
	return sum;
}

func (this *PlanetarySystem) Print() string{
	buff := fmt.Sprintf("\nAfter %d steps, %d Total Energy:\n", this.TimeIndex, this.TotalEnergy());
	for _, body := range this.Bodies{
		buff += body.Print() + "\n";
	}
	return buff;
}


func (this *PlanetarySystem) Step() {
	for i, body := range this.Bodies{
		for j, neighbor := range this.Bodies{
			if(i == j){
				continue;
			}
			if(body.Position.X != neighbor.Position.X){
				if(body.Position.X > neighbor.Position.X){
					body.Velocity.X--;
				} else{
					body.Velocity.X++;
				}
			}
			if(body.Position.Y != neighbor.Position.Y){
				if(body.Position.Y > neighbor.Position.Y){
					body.Velocity.Y--;
				} else{
					body.Velocity.Y++;
				}
			}
			if(body.Position.Z != neighbor.Position.Z){
				if(body.Position.Z > neighbor.Position.Z){
					body.Velocity.Z--;
				} else{
					body.Velocity.Z++;
				}
			}
		}
	}
	for _, body := range this.Bodies{
		body.ApplyVelocity();

	}

	this.TimeIndex++;

	if(!this.HasRestedX){
		if(this.IsAtRestOnX()){
			this.HasRestedX = true;
			Log.Info("Rested on X on step %d", this.TimeIndex);
			this.XPeriod = this.TimeIndex;
		}
	}
	if(!this.HasRestedY){
		if(this.IsAtRestOnY()){
			this.HasRestedY = true;
			Log.Info("Rested on Y on step %d", this.TimeIndex);
			this.YPeriod = this.TimeIndex;
		}
	}
	if(!this.HasRestedZ){
		if(this.IsAtRestOnZ()){
			this.HasRestedZ = true;
			Log.Info("Rested on Z on step %d", this.TimeIndex);
			this.ZPeriod = this.TimeIndex;
		}
	}
}



func (this *PlanetarySystem) IsAtRest() bool {
	for _, body := range this.Bodies{
		if(!body.IsAtRest()){
			return false;
		}
	}
	return true;
}

func (this *PlanetarySystem) IsAtRestOnX() bool {
	for _, body := range this.Bodies{
		if(!body.IsAtRestX()){
			return false;
		}
	}
	return true;
}

func (this *PlanetarySystem) IsAtRestOnY() bool {
	for _, body := range this.Bodies{
		if(!body.IsAtRestY()){
			return false;
		}
	}
	return true;
}

func (this *PlanetarySystem) IsAtRestOnZ() bool {
	for _, body := range this.Bodies{
		if(!body.IsAtRestZ()){
			return false;
		}
	}
	return true;
}



func (this *PlanetarySystem) HasCycled() bool {
	for _, body := range this.Bodies{
		if(!body.IsAtRest()){
			return false;
		}
	}
	return true;
}

func (this *PlanetarySystem) HasRestedOnAllAxises() bool {
	return this.HasRestedX && this.HasRestedY && this.HasRestedZ;
}

func (this *PlanetarySystem) Load(fileName string) error {
	Log.Info("Parsing planetary system from %s", fileName)


	this.LastBodyId = 1;
	this.FileName = fileName;
	this.Bodies = make([]*PlanetaryBody, 0);


	file, err := os.Open(fileName);
	if err != nil {
		Log.FatalError(err);
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)


	for scanner.Scan() {
		lineRaw := strings.TrimSpace(scanner.Text());
		if(lineRaw != ""){
			lineRaw = strings.Replace(lineRaw, "<x=", "", -1);
			lineRaw = strings.Replace(lineRaw, "y=", "", -1);
			lineRaw = strings.Replace(lineRaw, "z=", "", -1);
			lineRaw = strings.Replace(lineRaw, ">", "", -1);
			parts := strings.Split(lineRaw,",");
			xRaw := strings.TrimSpace(parts[0]);
			yRaw := strings.TrimSpace(parts[1]);
			zRaw := strings.TrimSpace(parts[2]);
			xParsed, err := strconv.ParseInt(xRaw, 10, 64);
			if(err != nil){
				return err;
			}
			yParsed, err := strconv.ParseInt(yRaw, 10, 64);
			if(err != nil){
				return err;
			}
			zParsed, err := strconv.ParseInt(zRaw, 10, 64);
			if(err != nil){
				return err;
			}


			body := &PlanetaryBody{};
			body.Position = &IntVec3{X:int(xParsed), Y: int(yParsed), Z: int(zParsed)};
			body.OriginalPosition = &IntVec3{X:int(xParsed), Y: int(yParsed), Z: int(zParsed)};
			body.Velocity = &IntVec3{X:int(0), Y: int(0), Z: int(0)};
			body.Id = this.LastBodyId;
			this.LastBodyId++;
			this.Bodies = append(this.Bodies, body);
		}
	}

	Log.Info("Completed parsing from %s, %d planetary bodies", fileName, len(this.Bodies));

	Log.Info(this.Print());

	return nil;
}