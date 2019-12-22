package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type Deck struct {
	Cards []int;
	ScratchCopy []int;
	Table []int;
	CardOfInterest int;
	CardOfInterestIndex int;
	Size int;
}

const DeckOperationDealOpCode = "deal";
const DeckOperationCutOpCode = "cut";
const DeckOperationResetOpCode = "reset";
const DeckOperationDealIntoStackLiteral = "deal into new stack";

func (this *Deck) Init(size int) {
	arr := make([]int, size);
	this.Cards = arr;
	arr = make([]int, size);
	this.ScratchCopy = arr;
	arr = make([]int, size);
	this.Table = arr;

	this.Size = size;
	this.Reset();
}

func (this *Deck) InitSlim(size int) {
	this.Size = size;
}


func (this *Deck) Reset() {
	for i := 0; i < this.Size; i++{
		this.Cards[i] = i;
	}
}

func (this *Deck) ResetSlim() {
	this.CardOfInterestIndex = this.CardOfInterest+1;
}

func (this *Deck) ResetScratch() {
	for i := 0; i < this.Size; i++{
		this.ScratchCopy[i] = 0;
	}
}

func (this *Deck) DealIntoNewStack() {
	this.Reverse();
}

func (this *Deck) DealIntoNewSlackSlim() {
	this.CardOfInterestIndex = this.Size - 1  - this.CardOfInterestIndex;
}

func (this *Deck) Reverse() {
	for i := 0; i < this.Size; i++{
		this.ScratchCopy[this.Size - i - 1] = this.Cards[i];
	}
	swp := this.Cards;
	this.Cards = this.ScratchCopy;
	this.ScratchCopy = swp;
}



func (this *Deck) ParseShuffleInstructionSet(fileName string) ([]DeckOperation, error) {
	file, err := ioutil.ReadFile(fileName);
	if err != nil {
		return nil, err;
	}
	res := make([]DeckOperation, 0);
	fileContents := strings.TrimSpace(string(file));
	parts := strings.Split(fileContents, "\n");
	lineNum := 1;
	for _, val := range parts {
		trimmed := strings.TrimSpace(val);
		if(trimmed != ""){
			lineParts := strings.Split(trimmed, " ");

			opCode := lineParts[0];
			var instruction DeckOperation;
				switch(opCode){
					case DeckOperationDealOpCode:
						if(trimmed == DeckOperationDealIntoStackLiteral){
							instruction = &DeckOperationDealToNewStack{};
						} else{
							instruction = &DeckOperationDeal{};
						}
						break;
					case DeckOperationCutOpCode:
						instruction = &DeckOperationCut{};
						break;
					case DeckOperationResetOpCode:
						instruction = &DeckOperationReset{};
						break;
					case "//": // Comment
						continue;
				default:
					return nil, errors.New("Unknown deck operation " + opCode);
			}
			if(instruction != nil){
				err = instruction.Parse(trimmed, lineNum);
				if(err != nil){
					return nil, err;
				}
				res = append(res, instruction);
				lineNum++;
			}
		}
	}
	return res, nil;
}

func (this *Deck) DescribeOperations(ops []DeckOperation) string {
	var str strings.Builder;
	for _, op := range ops{
		str.WriteString(op.ToString());
		str.WriteString("\n");
	}
	return str.String();
}

func (this *Deck) Apply(ops []DeckOperation) {
	for _, op := range ops{
		op.Apply(this);
		this.ResetScratch();
	}
}

func (this *Deck) SetCardOfInterest(card int){
	this.CardOfInterest = card;
	this.CardOfInterestIndex = this.CardOfInterest;
}

func (this *Deck) ApplySlim(ops []DeckOperation) {

	for _, op := range ops{
		op.ApplySlim(this);
	}
}

func (this *Deck) IndexOf(cardNum int) int {
	for i, v := range this.Cards{
		if(v == cardNum){
			return i;
		}
	}
	return -1;
}

func (this *Deck) CutSlim(n int) {
	if(n == 0){
		return;
	}
	if(n > 0){
		// This cut occurs before us. Shift left
		if(this.CardOfInterestIndex >= n){
			this.CardOfInterestIndex -= n;
		} else{
			// We are part of this cut. Our new index is offset from the end
			this.CardOfInterestIndex = (this.Size - n) + this.CardOfInterestIndex;
		}
	} else{
		n = n * -1;
		pivot := this.Size - n;
		if(this.CardOfInterestIndex < pivot){ // We are not in this cut, move to the right
			this.CardOfInterestIndex += n;
		} else{ // We are in the cut. Move to the front, our new index is our distance from the pivot
			this.CardOfInterestIndex = this.CardOfInterestIndex - pivot;
		}
	}
}




func (this *Deck) Cut(n int) {
	if(n == 0){
		return;
	}
	// Back slice
	if(n > 0){
		pivot := this.Size - n;
		for i := 0; i < n; i++{
			this.ScratchCopy[pivot+i] = this.Cards[i];
		}
		//Front Slice
		for i := 0; i < this.Size - n; i++{
			this.ScratchCopy[i] = this.Cards[i+n];
		}
	} else if(n < 0){
		n = n * -1;
		pivot := this.Size - n;

		// Front half
		for i := 0; i < n; i++{
			this.ScratchCopy[i] = this.Cards[pivot + i];
		}
		// Back half
		for i := n; i < this.Size; i++{
			this.ScratchCopy[i] = this.Cards[i-n];
		}

	}
	swp := this.Cards;
	this.Cards = this.ScratchCopy;
	this.ScratchCopy = swp;

}

func (this *Deck) Deal(n int) {
	index := 0;
	for _, i := range this.Cards {
		this.ScratchCopy[index] = i;
		index += n;
		index = index % this.Size;
	}
	swp := this.Cards;
	this.Cards = this.ScratchCopy;
	this.ScratchCopy = swp;
}

func (this *Deck) DealSlim(n int) {
	this.CardOfInterestIndex = (this.CardOfInterestIndex * n) % this.Size;
}

func (this *Deck) Print() string {
	var str strings.Builder;
	for i := 0; i < this.Size; i++{
		if(i > 0){
			str.WriteString(",");
		}
		str.WriteString(strconv.Itoa(this.Cards[i]));
	}
	return str.String();
}

type DeckOperation interface {
	ToString() string;
	Parse(line string, lineNum int) error;
	Apply(deck *Deck);
	ApplySlim(deck *Deck);
}

type DeckOperationDealToNewStack struct {
	LineNum int;
}

func (this *DeckOperationDealToNewStack) ToString() string{
	return fmt.Sprintf("%d - deal into new stack", this.LineNum);
}

func (this *DeckOperationDealToNewStack) Parse(line string, lineNum int) error{
	this.LineNum = lineNum;
	return nil;
}

func (this *DeckOperationDealToNewStack) Apply(deck *Deck){
	deck.DealIntoNewStack();
}

func (this *DeckOperationDealToNewStack) ApplySlim(deck *Deck){
	deck.DealIntoNewSlackSlim();
}


type DeckOperationCut struct {
	LineNum int;
	CutAmount int;
}

func (this *DeckOperationCut) ToString() string{
	return fmt.Sprintf("%d - cut %d", this.LineNum, this.CutAmount);
}

func (this *DeckOperationCut) Parse(line string, lineNum int) error{
	this.LineNum = lineNum;
	parts := strings.Split(line, " ");
	if(len(parts) < 2){
		return errors.New("Failed to parse cut, not enough parts");
	}
	parsedInt, err := strconv.ParseInt(parts[1], 10, 64);
	if(err != nil){
		return err;
	}
	this.CutAmount = int(parsedInt);
	return nil;
}

func (this *DeckOperationCut) Apply(deck *Deck){
	deck.Cut(this.CutAmount);
}

func (this *DeckOperationCut) ApplySlim(deck *Deck){
	deck.CutSlim(this.CutAmount);
}




type DeckOperationDeal struct {
	LineNum int;
	DealAmount int;
}

func (this *DeckOperationDeal) ToString() string{
	return fmt.Sprintf("%d - deal with increment %d", this.LineNum, this.DealAmount);
}

func (this *DeckOperationDeal) Parse(line string, lineNum int) error{
	this.LineNum = lineNum;
	parts := strings.Split(line, " ");
	if(len(parts) < 4){
		return errors.New("Failed to parse deal, not enough parts");
	}
	parsedInt, err := strconv.ParseInt(parts[3], 10, 64);
	if(err != nil){
		return err;
	}
	this.DealAmount = int(parsedInt);
	return nil;
}

func (this *DeckOperationDeal) Apply(deck *Deck){
	deck.Deal(this.DealAmount);
}

func (this *DeckOperationDeal) ApplySlim(deck *Deck){
	deck.DealSlim(this.DealAmount);
}


type DeckOperationReset struct {
	LineNum int;
}

func (this *DeckOperationReset) ToString() string{
	return fmt.Sprintf("%d - reset", this.LineNum);
}

func (this *DeckOperationReset) Parse(line string, lineNum int) error{
	return nil;
}

func (this *DeckOperationReset) Apply(deck *Deck){
	deck.Reset();
}
func (this *DeckOperationReset) ApplySlim(deck *Deck){
	deck.ResetSlim();
}
