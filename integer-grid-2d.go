package main

import (
	"errors"
	"strconv"
	"strings"
)

type IntegerGrid2D struct {
	Data map[int]*map[int]int;
}

func (this *IntegerGrid2D) Init() {
	this.Data = make(map[int]*map[int]int);
}

func (this *IntegerGrid2D) Visit(x int, y int) int {
	_, exists := this.Data[x];
	if(!exists){
		newMap := make(map[int]int);
		this.Data[x] = &newMap;
	}
	rowData := *this.Data[x];
	_, exists = rowData[y];
	if(!exists){
		rowData[y] = 0;
	}
	rowData[y]++;
	return rowData[y];
}

func (this *IntegerGrid2D) IsVisited(x int, y int) bool {
	return this.GetValue(x, y) > 0;
}

func (this *IntegerGrid2D) SetValue(x int, y int, value int) {
	_, exists := this.Data[x];
	if(!exists){
		newMap := make(map[int]int);
		this.Data[x] = &newMap;
	}
	rowData := *this.Data[x];
	_, exists = rowData[y];
	if(!exists){
		rowData[y] = 0;
	}
	rowData[y] = value;
}


func (this *IntegerGrid2D) GetValue(x int, y int) int {
	_, exists := this.Data[x];
	if(!exists){
		return 0;
	}
	rowData := *this.Data[x];
	_, exists = rowData[y];
	if(!exists){
		return 0;
	}
	return rowData[y];
}

type TwistyLine struct {
	Id  			int;
	LineSegments	[]*TwistyLineSegment;
	VisitedGrid		*IntegerGrid2D;
}



const LineDirectionUp = 0;
const LineDirectionDown = 1;
const LineDirectionLeft = 2;
const LineDirectionRight = 3;

func (this *TwistyLine) Parse(line string) error {
	this.LineSegments = make([]*TwistyLineSegment, 0);
	this.VisitedGrid = &IntegerGrid2D{};
	this.VisitedGrid.Init();
	parts := strings.Split(line, ",");
	for _, part := range parts{
		trimmed := strings.TrimSpace(part);
		if(trimmed != ""){
			direction := string(trimmed[0]);
			mag := trimmed[1:];
			seg := &TwistyLineSegment{};

			magVal, err := strconv.ParseInt(mag, 10, 64);
			if(err != nil){
				return err;
			}
			seg.Magnitude = int(magVal);
			switch(direction){
				case "U":
					seg.Direction = LineDirectionUp;
					break;
				case "D":
					seg.Direction = LineDirectionDown;
					break;
				case "L":
					seg.Direction = LineDirectionLeft;
					break;
				case "R":
					seg.Direction = LineDirectionRight;
					break;
				default:
					 return errors.New("Unknown line direction " + direction);
			}


			this.LineSegments = append(this.LineSegments, seg);
		}
	}
	return nil;
}



type TwistyLineSegment struct {
	Magnitude int;
	Direction int;
}

func (this *TwistyLine) StepsToIntersection(intersection *IntVec2) int {
	if(!this.VisitedGrid.IsVisited(intersection.X, intersection.Y)){
		return -1;
	}
	return this.VisitedGrid.GetValue(intersection.X, intersection.Y);
}


func (this *TwistyLine) Apply(grid *IntegerGrid2D) []*IntVec2  {
	res := make([]*IntVec2, 0);
	x := 0;
	y := 0;
	step := 0;
	for _, segment := range this.LineSegments {
		for i := 0; i < segment.Magnitude; i++ {
			switch segment.Direction {
				case LineDirectionUp:
					y++;
					break;
				case LineDirectionDown:
					y--;
					break;
				case LineDirectionLeft:
					x--;
					break;
				case LineDirectionRight:
					x++;
					break;
			}
			step++;
			if(x == 0 && y == 0){ // Don't bother marking the origin
				continue;
			}
			if(this.VisitedGrid.IsVisited(x, y)){
				continue;
			}
			this.VisitedGrid.SetValue(x, y, step);
			intersects := grid.Visit(x, y);
			if(intersects > 1){
				point := &IntVec2{};
				point.X = x;
				point.Y = y;
				res = append(res, point);
			}
		}
	}
	return res;
}