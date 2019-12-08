package main

import (
	"github.com/nfnt/resize"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"os"
	"strconv"
)

type IntImage struct {
	FileName string;
	Width int;
	Height int;
	LastLayerId int;
	Layers []*IntImageLayer;
}

func (this *IntImage) Parse(width int, height int, fileName string) error {
	this.Width = width;
	this.Height = height;
	this.FileName = fileName;
	this.LastLayerId = 1;
	Log.Info("Parsing int image machine from %s", fileName)
	this.Layers = make([]*IntImageLayer, 0);

	fileContents, err := ioutil.ReadFile(fileName);
	if(err != nil){
		return err;
	}

	layerSize := this.Width * this.Height;
	parsedInLayer := 0;
	var layer *IntImageLayer;

	for i, strV := range fileContents{

		if(layer == nil || parsedInLayer >= layerSize){
			parsedInLayer = 0;
			layer = &IntImageLayer{};
			layer.Init(this, this.LastLayerId, i);
			this.Layers = append(this.Layers, layer);
			this.LastLayerId++;
		}
		parsedInLayer++;

		v, err := strconv.ParseInt(string(strV), 10, 64);
		if(err != nil){
			return err;
		}

		layer.AddValue(int(v));
	}
	return nil;
}

type IntImageLayer struct {
	LayerId int;
	Hist map[int]int;
	StartIndex int;
	Data [][]int;
	ContainingImage *IntImage;
	ValueCount int;
}

func (this *IntImageLayer) Init(image *IntImage, layerId int, startIndex int) {
	this.ContainingImage = image;
	this.LayerId = layerId;
	this.Hist = make(map[int]int);
	this.Data = make([][]int, image.Height);
	for i := 0; i < image.Height; i++ {
		this.Data[i] = make([]int, image.Width);
	}

	this.StartIndex = startIndex;
}

func (this *IntImageLayer) AddValue(value int) {


	row := this.ValueCount / this.ContainingImage.Width;
	col := this.ValueCount - (row * this.ContainingImage.Width);

	//Log.Info("Adding %d - %d,%d - Count %d", value, row, col, this.ValueCount);

	this.Data[row][col] = value;

	_, exists := this.Hist[value];
	if(!exists){
		this.Hist[value] = 0;
	}
	this.Hist[value]++;

	this.ValueCount++;
}

func (this *IntImageLayer) GetHistValue(value int) int {

	val, exists := this.Hist[value];
	if(!exists){
		return 0;
	}
	return val;
}


func (this *IntImage) Flatten() [][]int{
	transparentPixel := 2;
	data := make([][]int, this.Height);
	for i := 0; i < this.Height; i++ {
		data[i] = make([]int, this.Width);
		for j, _ := range data[i]{
			data[i][j] = transparentPixel;
		}
	}
	for _, layer := range this.Layers {
		for i, row := range layer.Data {
			for j, v := range row{

				if(data[i][j] == transparentPixel){
					data[i][j] = v;
				}


			}
		}
	}
	return data;
}

func (this *IntImage) FlattenAndDraw() string {


	//whitePixel := 1;
	//blackPixel := 0;

	data := this.Flatten();


	buff := "\n";
	for rowIndex, _ := range data {
		row := data[rowIndex];
		rowBuff := "";
		for _, j := range row{
			if(j == 0){
				rowBuff += " ";
			} else{
				rowBuff += strconv.Itoa(j);
			}

		}
		buff += rowBuff + "\n";
	}

	return buff;
}

func (this *IntImage) FlattenAndRenderToFile(fileName string, targetWidth int) {


	//whitePixel := 1;
	//blackPixel := 0;

	data := this.Flatten();

	padding := 2;

	width := this.Width + (padding*2);
	height := this.Height +  (padding*2);

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	for x := 0; x < width; x++ {
		for y := 0; y < width; y++ {
			img.Set(x, y, color.White);
		}
	}

	// Set color for each pixel.
	for x := 0; x < this.Width; x++ {
		for y := 0; y < this.Height; y++ {
			if(data[y][x] != 0){
				img.Set(x+padding, y+padding, color.Black)
			}
		}
	}

	// Encode
	m := resize.Resize(uint(targetWidth), 0, img, resize.MitchellNetravali)
	//as PNG.
	f, _ := os.Create(fileName)
	png.Encode(f, m)
	Log.Info("Rendered int image %s to %s ", this.FileName, fileName);
}



func (this *IntImageLayer) Draw() string {

	buff := "";
	for rowIndex, _ := range this.Data {
		row := this.Data[rowIndex];
		rowBuff := "";
		for _, j := range row{
			rowBuff += strconv.Itoa(j);
		}
		buff += rowBuff + "\n";
	}

	return buff;
}

func (this *IntImage) Draw() string {

	buff := "\n";

	for _, layer := range this.Layers {
		buff += "Layer " + strconv.Itoa(layer.LayerId);
		buff += "\n";
		buff += layer.Draw();
		buff += "\n";
	}

	return buff;
}