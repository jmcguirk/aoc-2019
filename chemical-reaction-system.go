package main

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type ChemicalReactionSystem struct {
	Formulae map[string]*ReactionFormula;
	FileName string;
}

type ReactionFormula struct {
	ProductId string;
	Input []*ReactionProduct;
	Output *ReactionProduct;
}

type ReactionProduct struct {
	ProductId string;
	Quantity int;
}

type ReactionState struct {
	Held	map[string]int;
	Desired	map[string]int;
	TotalInputRequired int;
	InputDesired string;
}

func (this *ReactionProduct) Parse(stub string) error{
	parts := strings.Split(stub, " ");
	if(len(parts) != 2){
		return errors.New("Parse failure - wrong product parts");
	}
	qty, err := strconv.ParseInt(strings.TrimSpace(parts[0]), 10, 64);
	if(err != nil){
		return err;
	}
	this.ProductId = strings.TrimSpace(parts[1]);
	this.Quantity = int(qty);
	return nil;
}

func (this *ReactionFormula) Parse(line string) error{

	parts := strings.Split(line, "=>");

	if(len(parts) != 2){
		return errors.New("Parse failure - wrong parts");
	}

	inputRaw := strings.TrimSpace(parts[0]);
	outputRaw := strings.TrimSpace(parts[1]);

	this.Input = make([]*ReactionProduct, 0);

	inputParts := strings.Split(inputRaw, ",");
	for _, v := range inputParts{
		inputProduct := &ReactionProduct{}
		err := inputProduct.Parse(strings.TrimSpace(v));
		if(err != nil){
			return err;
		}
		this.Input = append(this.Input, inputProduct);
	}

	this.Output = &ReactionProduct{};
	err := this.Output.Parse(outputRaw);
	if(err != nil){
		return err;
	}

	this.ProductId = this.Output.ProductId;

	return nil;
}


func (this *ReactionFormula) Print() string{
	buff := "";
	for _, input := range this.Input {
		if(buff != ""){
			buff += ", ";
		}
		buff += fmt.Sprintf("%d %s", input.Quantity, input.ProductId);
	}
	buff += " => ";
	buff += fmt.Sprintf("%d %s", this.Output.Quantity, this.Output.ProductId);
	return buff;
}


func (this *ChemicalReactionSystem) PrintFormulae() string{
	buff := "\n";
	for _, formula := range this.Formulae{
		buff += formula.Print() + "\n";
	}

	return buff;
}


func (this *ChemicalReactionSystem) ApplyFormula(formula *ReactionFormula, quantityDesired int, state *ReactionState) (error){

	desiredAmount := quantityDesired;
	v, exists :=state.Held[formula.ProductId];
	if(exists){
		if(desiredAmount > v){
			state.Held[formula.ProductId] = 0;
			desiredAmount -= v;
		} else{
			state.Held[formula.ProductId] = v - desiredAmount;
			desiredAmount = 0;
		}

	}



	if(desiredAmount > 0){
		formulaTimes := int(math.Ceil(float64(desiredAmount)/float64(formula.Output.Quantity)));
		excess := (formulaTimes * formula.Output.Quantity) - desiredAmount;
		for _, i := range formula.Input {
			inputAmountNeeded := formulaTimes * i.Quantity;
			if(i.ProductId == state.InputDesired){
				//Log.Info("Spending %d %s to buy %d %s - desired - %d", inputAmountNeeded, i.ProductId, formula.Output.Quantity, formula.Output.ProductId, desiredAmount);
				state.TotalInputRequired += inputAmountNeeded;
			} else{

				_, exists :=state.Desired[i.ProductId];
				if(!exists){
					state.Desired[i.ProductId] = 0;
				}
				state.Desired[i.ProductId] += inputAmountNeeded;
				//Log.Info("Added input %d x %s to shopping list because we wanted %d x %s", inputAmountNeeded, i.ProductId, desiredAmount, formula.ProductId);
			}
		}
		if(excess > 0){
			_, exists :=state.Held[formula.Output.ProductId];
			if(!exists){
				state.Held[formula.Output.ProductId] = 0;
			}
			state.Held[formula.Output.ProductId] += excess;
		}

	}

	delete(state.Desired, formula.Output.ProductId);

	for k, d := range state.Desired{
		v, exists :=state.Held[k];
		if(exists){
			if(d > v){
				state.Held[k] = 0;
				state.Desired[k] = d - v;
			} else{
				state.Held[k] = v - d;
				state.Desired[k] = 0;
			}

		}
	}

	return nil;
}

func (this *ChemicalReactionSystem) ApplyFormulae(state *ReactionState) (error){

	if(len(state.Desired) <= 0){
		return nil;
	}

	var firstProductId string;
	var firstQuantityDesired int;
	for k, v := range state.Desired {
		firstProductId = k;
		firstQuantityDesired = v;
		break;
	}


	if(firstQuantityDesired <= 0){
		delete(state.Desired, firstProductId);
	} else{
		formula := this.Formulae[firstProductId];
		if(formula == nil){
			return errors.New("Couldn't convert " + firstProductId);
		}
		err := this.ApplyFormula(formula, firstQuantityDesired, state);
		if(err != nil){
			return err;
		}
	}


	if(len(state.Desired) > 0){
		return this.ApplyFormulae(state);
	}
	return nil;
}


func (this *ChemicalReactionSystem) GetTotalInputRequired(productId string, amountDesired int, desiredInput string) (int, error){

	_, exists := this.Formulae[productId];
	if(!exists){
		return -1, errors.New("Failed to find formula for product id")
	}

	reactionState := &ReactionState{};
	reactionState.Held = map[string]int{};
	reactionState.InputDesired = desiredInput;
	reactionState.Desired = map[string]int{};
	reactionState.Desired[productId] = amountDesired;

	err := this.ApplyFormulae(reactionState);
	if(err != nil){
		return -1, err;
	}

	return reactionState.TotalInputRequired, nil;
}


func (this *ChemicalReactionSystem) Load(fileName string) error {
	this.FileName = fileName;
	this.Formulae = make(map[string]*ReactionFormula);


	file, err := os.Open(fileName);
	if err != nil {
		Log.FatalError(err);
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)


	for scanner.Scan() {
		lineRaw := strings.TrimSpace(scanner.Text());
		if(lineRaw != ""){
			formula := &ReactionFormula{}
			err = formula.Parse(lineRaw);
			if err != nil {
				Log.FatalError(err);
			}
			this.Formulae[formula.ProductId] = formula;
		}
	}


	Log.Info("Completed parsing from %s, %d formulae", fileName, len(this.Formulae));

	//Log.Info(this.PrintFormulae());

	return nil;
}