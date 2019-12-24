package main

import (
	"fmt"
	"math"

)

type LayeredGrid struct {
	Layers map[int]*GridLayer;
	DefaultValue int;
}

type GridLayer struct {
	ContainingGrid *LayeredGrid;
	LayerId int;
	Data map[int]*map[int]int;
}

func (this *LayeredGrid) Init(defaultValue int) {
	this.Layers = make(map[int]*GridLayer);
	this.DefaultValue = defaultValue;
}

func (this *LayeredGrid) GetOrCreateLayer(layerId int) *GridLayer {
	grid, exists := this.Layers[layerId];
	if(exists){
		return grid;
	}

	if(layerId != 0){
		cpy := this.Layers[0].CloneInto(this);
		cpy.LayerId = layerId;
		for _, v := range cpy.Data{
			row := *v;
			for j, _ := range row{
				row[j] = this.DefaultValue;
			}
		}
		this.Layers[layerId] = cpy;
		return cpy;
	} else{
		grid = &GridLayer{}
		grid.Init(this, layerId);
		this.Layers[layerId] = grid;
		return grid;
	}

}

func (this *GridLayer) Init(containing *LayeredGrid, layerId int) {
	this.ContainingGrid = containing;
	this.LayerId = layerId;
	this.Data = make(map[int]*map[int]int);
}

func (this *GridLayer) CloneInto(grid *LayeredGrid) *GridLayer {
	res := &GridLayer{};
	res.Init(grid, this.LayerId);


	for k, v := range this.Data{
		cpy := make(map[int]int);
		for j, v2 := range *v{
			cpy[j] = v2;
		}
		res.Data[k] = &cpy;
	}

	return res;
}

func (this *LayeredGrid) Clone() *LayeredGrid {
	res := &LayeredGrid{};
	res.Layers = make(map[int]*GridLayer);
	for k, v := range this.Layers{
		res.Layers[k] = v.CloneInto(res);
	}

	return res;
}



func (this *GridLayer) Visit(x int, y int) int {
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

func (this *GridLayer) IsVisited(x int, y int) bool {
	return this.GetValue(x, y) > 0;
}

func (this *GridLayer) SetValue(x int, y int, value int) {
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

func (this *LayeredGrid) SetValue(x int, y int, layerId int, value int) {
	this.GetOrCreateLayer(layerId).SetValue(x, y, value);
}

func (this *LayeredGrid) GetValue(x int, y int, layerId int) int {
	return this.GetOrCreateLayer(layerId).GetValue(x, y);
}

func (this *LayeredGrid) HasValue(x int, y int, layerId int) bool {
	return this.GetOrCreateLayer(layerId).HasValue(x, y);
}

func (this *LayeredGrid) CountValue(value int) int {
	total := 0;
	for _, l := range this.Layers {
		total += l.CountValue(value);
	}
	return total;
}

func (this *GridLayer) CountValue(value int) int {
	total := 0;
	for _, v := range this.Data{
		for _, v2 := range *v{
			if(v2 == value){
				total++;
			}
		}
	}
	return total;
}

func (this *LayeredGrid) MinLayer() int {
	res := math.MaxInt64;
	for _, layer := range this.Layers{
		if(layer.LayerId < res){
			res = layer.LayerId;
		}
	}
	return res;
}

func (this *LayeredGrid) MaxLayer() int {
	res := math.MinInt64;
	for _, layer := range this.Layers{
		if(layer.LayerId > res){
			res = layer.LayerId;
		}
	}
	return res;
}


func (this *LayeredGrid) PrintAscii() string {

	layerMin := this.MinLayer();
	layerMax := this.MaxLayer();
	buff := "";
	for i := layerMin; i <= layerMax; i++{
		buff += fmt.Sprintf("\nLayer %d:\n%s\n",i, this.GetOrCreateLayer(i).PrintAscii());
	}

	//Log.Info("Furthest point is %d,%d", furthestX, furthestY);

	return buff;
}

func (this *GridLayer) PrintAscii() string {
	xMin := this.MinRow();
	xMax := this.MaxRow();

	yMin := this.MinCol();
	yMax := this.MaxCol();

	buff := "";
	for j := yMin; j<= yMax; j++{
		for i := xMin; i<= xMax; i++{
			if(!this.HasValue(i, j)){
				buff += " ";
			} else{
				val := this.GetValue(i, j);
				if(val > 0){
					buff += fmt.Sprintf("%c", val);
				} else{
					buff += " ";
				}
			}
		}
		buff += "\n";
	}
	//Log.Info("Furthest point is %d,%d", furthestX, furthestY);

	return buff;
}

func (this *GridLayer) MaxRow() int {
	res := math.MinInt32;
	for x, _ := range this.Data{
		if(x > res){
			res = x;
		}
	}
	return res;
}

func (this *GridLayer) MinRow() int {
	res := math.MaxInt32;
	for x, _ := range this.Data{
		if(x < res){
			res = x;
		}
	}
	return res;
}

func (this *GridLayer) MaxCol() int {
	res := math.MinInt32;
	for _, vals := range this.Data{
		for y, _ := range *vals{
			if(y > res){
				res = y;
			}
		}
	}
	return res;
}

func (this *GridLayer) MinCol() int {
	res := math.MaxInt32;
	for _, vals := range this.Data{
		for y, _ := range *vals{
			if(y < res){
				res = y;
			}
		}
	}
	return res;
}



func (this *GridLayer) MaxX() int {
	res := math.MinInt32;
	for x, _ := range this.Data{
		if(x > res){
			res = x;
		}
	}
	return res;
}

func (this *GridLayer) MinX() int {
	res := math.MaxInt32;
	for x, _ := range this.Data{
		if(x < res){
			res = x;
		}
	}
	return res;
}

func (this *GridLayer) MaxY() int {
	res := math.MinInt32;
	for _, vals := range this.Data{
		for y, _ := range *vals{
			if(y > res){
				res = y;
			}
		}
	}
	return res;
}

func (this *GridLayer) MinY() int {
	res := math.MaxInt32;
	for _, vals := range this.Data{
		for y, _ := range *vals{
			if(y < res){
				res = y;
			}
		}
	}
	return res;
}



func (this *GridLayer) GetValue(x int, y int) int {
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

func (this *GridLayer) HasValue(x int, y int) bool {
	_, exists := this.Data[x];
	if(!exists){
		return false;
	}
	rowData := *this.Data[x];
	_, exists = rowData[y];
	if(!exists){
		return false;
	}
	return true;
}