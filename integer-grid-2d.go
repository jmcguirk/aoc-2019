package main

import (
	"errors"
	"github.com/nfnt/resize"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type IntegerGrid2D struct {
	Data map[int]*map[int]int;
}

func (this *IntegerGrid2D) Init() {
	this.Data = make(map[int]*map[int]int);
}

func (this *IntegerGrid2D) Clone() *IntegerGrid2D {
	res := &IntegerGrid2D{};
	res.Init();


	for k, v := range this.Data{
		cpy := make(map[int]int);
		for j, v2 := range *v{
			cpy[j] = v2;
		}
		res.Data[k] = &cpy;
	}

	return res;
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

func (this *IntegerGrid2D) PrintToFile(fileName string, targetWidth int) {
	xMin := this.MinRow();
	xMax := this.MaxRow();

	yMin := this.MinCol();
	yMax := this.MaxCol();

	padding := 3;

	baseWidth := xMax - xMin;
	baseHeight := yMax - yMin;

	width := baseWidth + (padding*2);
	height := baseHeight +  (padding*2);

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	for x := 0; x < width; x++ {
		for y := 0; y < width; y++ {
			img.Set(x, y, color.Black);
		}
	}


	for j := yMin; j<= yMax; j++{
		for i := xMin; i<= xMax; i++{
			val := this.GetValue(i, j);
			if(val > 0){
				img.Set(i+padding, j+padding, color.White);
			}
		}
	}
	// Encode
	m := resize.Resize(uint(targetWidth), 0, img, resize.MitchellNetravali)
	//as PNG.
	f, _ := os.Create(fileName)
	png.Encode(f, m)
	Log.Info("Rendered grid image %s ", fileName);

}

func (this *IntegerGrid2D) Print() string {
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
					buff += strconv.Itoa(this.GetValue(i, j));
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

func (this *IntegerGrid2D) MaxRow() int {
	res := math.MinInt32;
	for x, _ := range this.Data{
		if(x > res){
			res = x;
		}
	}
	return res;
}

func (this *IntegerGrid2D) MinRow() int {
	res := math.MaxInt32;
	for x, _ := range this.Data{
		if(x < res){
			res = x;
		}
	}
	return res;
}

func (this *IntegerGrid2D) MaxCol() int {
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

func (this *IntegerGrid2D) MinCol() int {
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



func (this *IntegerGrid2D) MaxX() int {
	res := math.MinInt32;
	for x, _ := range this.Data{
		if(x > res){
			res = x;
		}
	}
	return res;
}

func (this *IntegerGrid2D) MinX() int {
	res := math.MaxInt32;
	for x, _ := range this.Data{
		if(x < res){
			res = x;
		}
	}
	return res;
}

func (this *IntegerGrid2D) MaxY() int {
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

func (this *IntegerGrid2D) MinY() int {
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

func (this *IntegerGrid2D) HasValue(x int, y int) bool {
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

func (this *IntegerGrid2D) TileIndex(X int, Y int) int {
	return (X + TileIndexSize) + ((Y + TileIndexSize) * TileIndexOffset);
}

func (this *IntegerGrid2D) FromTileIndex(tileIndex int) (int, int) {
	y := (tileIndex/TileIndexOffset)-TileIndexSize;
	x := (tileIndex%TileIndexOffset)-TileIndexSize;
	return y, x;
}

func (this *IntegerGrid2D) GenerateEdges(from *IntVec2) []*IntVec2 {
	res := make([]*IntVec2, 0);
	north := from.Clone();
	north.Y--;
	res = append(res, north);

	south := from.Clone();
	south.Y++;
	res = append(res, south);

	east := from.Clone();
	east.X++;
	res = append(res, east);

	west := from.Clone();
	west.X--;
	res = append(res, west);

	return res;
}

func (this *IntegerGrid2D) ShortestPathWithBlacklist(from *IntVec2, to *IntVec2, blockValue []int) []*IntVec2 {

	//Log.Info("Requesting path from %s to %s", from.ToString(), to.ToString());
	res := make([]*IntVec2, 0);

	visitedNodes := &IntegerGrid2D{};
	visitedNodes.Init();
	minCostToStart := &IntegerGrid2D{};
	minCostToStart.Init();

	nearestToStart := make(map[int]int);

	frontier := make([]*IntVec2, 0);
	frontier = append(frontier, from);
	frontierMap := make(map[int]int);
	frontierMap[from.TileIndex()] = 1;

	minCostToStart.SetValue(from.X, from.Y, 1);

	for {
		if (len(frontier) <= 0) {
			break;
		}
		sort.SliceStable(frontier, func(i, j int) bool {
			vI := frontier[i];
			vJ := frontier[j];
			return minCostToStart.GetValue(vI.X, vI.Y) < minCostToStart.GetValue(vJ.X, vJ.Y);
		});

		next := frontier[0];
		frontier = frontier[1:];
		delete(frontierMap, next.TileIndex());
		costToHere := minCostToStart.GetValue(next.X, next.Y);
		edges := this.GenerateEdges(next);
		for _, edge := range edges{
			if(visitedNodes.HasValue(edge.X, edge.Y)){
				continue;
			}
			if(!this.HasValue(edge.X, edge.Y)){
				continue;
			}
			val := this.GetValue(edge.X, edge.Y);
			isBlackListed := false;
			for _, v := range blockValue{
				if(v == val){
					isBlackListed = true;
					break;
				}
			}
			if(isBlackListed){
				continue;
			}
			bestToHere := int(math.MaxInt32);

			bestCostExists := minCostToStart.HasValue(edge.X, edge.Y);
			if(bestCostExists){
				bestToHere = minCostToStart.GetValue(edge.X, edge.Y);
			}

			if(costToHere + 1 < bestToHere){
				minCostToStart.SetValue(edge.X, edge.Y, costToHere + 1);
				//Log.Info("Point %d to %d", edge.TileIndex(), next.TileIndex());
				nearestToStart[edge.TileIndex()] = next.TileIndex();
				//minCostToStart[neighbor.Id] = costToHere + edge.Weight;
				//nearestToStart[neighbor.Id] = next;
				_, alreadyEnqueued := frontierMap[edge.TileIndex()];
				if(!alreadyEnqueued){
					frontierMap[edge.TileIndex()] = 1;
					frontier = append(frontier, edge);
				}
			}

		}
		visitedNodes.SetValue(next.X, next.Y, 1);
		if(next.TileIndex() == to.TileIndex()){
			break;
		}
	}

	if(!minCostToStart.HasValue(to.X, to.Y)){
		// No path exists
		return nil;
	}

	//Log.Info("Done %d", from.TileIndex());
	nextPathStep := to.TileIndex();

	for {
		next := nearestToStart[nextPathStep];
		if(next == 0){
			log.Fatal("exit");
		}
		//Log.Info("Check %d to %d", nextPathStep, next);
		if(next == from.TileIndex()){
			break;
		}
		nextPathStep = next;
		node := &IntVec2{};
		node.FromTileIndex(next);
		res = append(res, node);
	}

	ReverseSlice(res);
	res = append(res, to);
	return res;
}

func (this *IntegerGrid2D) Reachable(from *IntVec2, blockValue int) []*IntVec2 {

	res := make([]*IntVec2, 0);
	visitedNodes := &IntegerGrid2D{};
	visitedNodes.Init();

	frontier := make([]*IntVec2, 0);
	frontier = append(frontier, from);
	frontierMap := make(map[int]int);
	frontierMap[from.TileIndex()] = 1;

	res = append(res, from)

	for {
		if (len(frontier) <= 0) {
			break;
		}
		next := frontier[0];
		frontier = frontier[1:];
		delete(frontierMap, next.TileIndex());
		edges := this.GenerateEdges(next);
		for _, edge := range edges {
			if(visitedNodes.HasValue(edge.X, edge.Y)){
				continue;
			}
			if(!this.HasValue(edge.X, edge.Y)){
				continue;
			}
			if(this.GetValue(edge.X, edge.Y) == blockValue){
				continue;
			}
			visitedNodes.SetValue(edge.X, edge.Y, 1);
			res = append(res, edge);
			_, alreadyEnqueued := frontierMap[edge.TileIndex()];
			if(!alreadyEnqueued){
				frontierMap[edge.TileIndex()] = 1;
				frontier = append(frontier, edge);
			}
		}
	}
	return res;
}

func (this *IntegerGrid2D) ShortestPath(from *IntVec2, to *IntVec2, blockValue int) []*IntVec2 {

	//Log.Info("Requesting path from %s to %s", from.ToString(), to.ToString());
	res := make([]*IntVec2, 0);

	visitedNodes := &IntegerGrid2D{};
	visitedNodes.Init();
	minCostToStart := &IntegerGrid2D{};
	minCostToStart.Init();

	nearestToStart := make(map[int]int);

	frontier := make([]*IntVec2, 0);
	frontier = append(frontier, from);
	frontierMap := make(map[int]int);
	frontierMap[from.TileIndex()] = 1;

	minCostToStart.SetValue(from.X, from.Y, 1);

	for {
		if (len(frontier) <= 0) {
			break;
		}
		sort.SliceStable(frontier, func(i, j int) bool {
			vI := frontier[i];
			vJ := frontier[j];
			return minCostToStart.GetValue(vI.X, vI.Y) < minCostToStart.GetValue(vJ.X, vJ.Y);
		});

		next := frontier[0];
		frontier = frontier[1:];
		delete(frontierMap, next.TileIndex());
		costToHere := minCostToStart.GetValue(next.X, next.Y);
		edges := this.GenerateEdges(next);
		for _, edge := range edges{
			if(visitedNodes.HasValue(edge.X, edge.Y)){
				continue;
			}
			if(!this.HasValue(edge.X, edge.Y)){
				continue;
			}
			if(this.GetValue(edge.X, edge.Y) == blockValue){
				continue;
			}
			bestToHere := int(math.MaxInt32);

			bestCostExists := minCostToStart.HasValue(edge.X, edge.Y);
			if(bestCostExists){
				bestToHere = minCostToStart.GetValue(edge.X, edge.Y);
			}

			if(costToHere + 1 < bestToHere){
				minCostToStart.SetValue(edge.X, edge.Y, costToHere + 1);
				//Log.Info("Point %d to %d", edge.TileIndex(), next.TileIndex());
				nearestToStart[edge.TileIndex()] = next.TileIndex();
				//minCostToStart[neighbor.Id] = costToHere + edge.Weight;
				//nearestToStart[neighbor.Id] = next;
				_, alreadyEnqueued := frontierMap[edge.TileIndex()];
				if(!alreadyEnqueued){
					frontierMap[edge.TileIndex()] = 1;
					frontier = append(frontier, edge);
				}
			}

		}
		visitedNodes.SetValue(next.X, next.Y, 1);
		if(next.TileIndex() == to.TileIndex()){
			break;
		}
	}

	if(!minCostToStart.HasValue(to.X, to.Y)){
		// No path exists
		return nil;
	}

	//Log.Info("Done %d", from.TileIndex());
	nextPathStep := to.TileIndex();

	for {
		next := nearestToStart[nextPathStep];
		if(next == 0){
			log.Fatal("exit");
		}
		//Log.Info("Check %d to %d", nextPathStep, next);
		if(next == from.TileIndex()){
			break;
		}
		nextPathStep = next;
		node := &IntVec2{};
		node.FromTileIndex(next);
		res = append(res, node);
	}

	ReverseSlice(res);
	res = append(res, to);
	return res;
}