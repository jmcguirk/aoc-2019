package main

import (
	"errors"
)

const RepairStateVisited = 1;
const RepairStateWall = -1;
const RepairStateOxygenSystem = 2;
const RepairStateO2 = 3;

const MoveResultWall = 0;
const MoveResultSuccess = 1;
const MoveResultOxygen = 2;

type RobotRepair struct {
	Processor *IntcodeMachineV3;
	Grid *IntegerGrid2D;
	Position *IntVec2;
	Id string;
	InstructionFileName string;
	StepCount int;
	WorkDone int;
	SystemLocation *IntVec2;
}

type UnexploredRepairNode struct {
	NewPosition *IntVec2;
	StartPosition *IntVec2;
	Direction int;
}

func (this *RobotRepair) Init(instructionFile string, initialPosition *IntVec2, grid *IntegerGrid2D) error {
	this.Processor = &IntcodeMachineV3{};
	this.Processor.PauseOnOutput = true;
	this.Processor.PauseOnInput = true;
	this.InstructionFileName = instructionFile;
	this.Position = &IntVec2{};
	this.Position.X = initialPosition.X;
	this.Position.Y = initialPosition.Y;
	this.Grid = grid;

	this.StepCount = 0;
	this.Grid.SetValue(initialPosition.X, initialPosition.Y, RepairStateVisited);

	err := this.Processor.Load(this.InstructionFileName)

	if err != nil {
		return err;
	}
	return nil;

}

func (this *RobotRepair) Explore() error {
	Log.Info("Repair robot starting exploration!");
	frontier := make([]*UnexploredRepairNode, 0);
	for{
		candidates := this.GetAllDirections(this.Position);
		for _, node := range candidates {
			if(!this.IsVisited(node)){
				frontier = append([]*UnexploredRepairNode{node}, frontier...);
			}
		}
		if(len(frontier) <= 0){
			break;
		}
		next := frontier[0];
		frontier = frontier[1:];
		if(!this.IsVisited(next)){
			if(!next.StartPosition.Eq(this.Position)){
				err := this.PathToPosition(next.StartPosition);
				if(err != nil){
					return err;
				}
			}
			err := this.ExploreNode(next);
			if(err != nil){
				return err;
			}
		}
	}
	return nil;
}

func (this *RobotRepair) ExploreNode(node *UnexploredRepairNode) error{
	//Log.Info("Exploring node res %s" , node.ToString());
	moveRes, err := this.Move(node.Direction);
	if(err != nil){
		return err;
	}
	if(moveRes == MoveResultWall){
		this.Grid.SetValue(node.NewPosition.X, node.NewPosition.Y, RepairStateWall);
	} else if (moveRes == MoveResultOxygen){
		this.Grid.SetValue(node.NewPosition.X, node.NewPosition.Y, RepairStateOxygenSystem);
		this.SystemLocation = &IntVec2{};
		this.SystemLocation = node.NewPosition.Clone();
	} else{
		this.Grid.SetValue(node.NewPosition.X, node.NewPosition.Y, RepairStateVisited);
	}
	return nil;
}

func (this *RobotRepair) Move(direction int) (int, error){
	this.StepCount++;
	this.Processor.QueueInput(int64(direction));
	res, err, halted := this.Processor.ReadNextOutput();
	if(err != nil){
		return -1, err;
	}
	if(halted){
		return -1, errors.New("Unexpected halt");
	}
	//Log.Info("Got res %d" , res);
	if(res == MoveResultSuccess || res == MoveResultOxygen){
		this.AdvancePosition(direction);
	}
	return int(res), nil;
}

func (this *RobotRepair) AdvancePosition(direction int){
	newPos := this.Position.Clone();
	switch(direction){
		case DirectionNorth:
			newPos.Y--;
			break;
		case DirectionSouth:
			newPos.Y++;
			break;
		case DirectionEast:
			newPos.X++;
			break;
		case DirectionWest:
			newPos.X--;
			break;
	}
	//Log.Info("advancing from " + this.Position.ToString() + " to " + newPos.ToString());
	this.Position = newPos;
}

func (this *RobotRepair) PathToPosition(pos *IntVec2) error{
	//Log.Info("Requesting path");
	path := this.Grid.ShortestPath(this.Position, pos, RepairStateWall);
	//Log.Info("Got path");
	if(path == nil || len(path) == 0){
		return errors.New("No path found");
	}
	for _, p := range path {
		if(p.X == this.Position.X && p.Y < this.Position.Y){
			_, err := this.Move(DirectionNorth);
			if(err != nil){
				return err;
			}
		} else if(p.X == this.Position.X && p.Y > this.Position.Y){
			_, err := this.Move(DirectionSouth);
			if(err != nil){
				return err;
			}
		} else if(p.X > this.Position.X && p.Y == this.Position.Y){
			_, err := this.Move(DirectionEast);
			if(err != nil){
				return err;
			}
		} else if(p.X < this.Position.X && p.Y == this.Position.Y){
			_, err := this.Move(DirectionWest);
			if(err != nil){
				return err;
			}
		} else{
			return errors.New("Broken path");
		}
	}
	return nil;
}

func (this *UnexploredRepairNode) ToString() string{
	buff := "";
	buff += this.StartPosition.ToString();
	buff += " - ";
	switch(this.Direction){
		case DirectionNorth:
			buff += "N";
			break;
		case DirectionSouth:
			buff += "S";
			break;
		case DirectionEast:
			buff += "E";
			break;
		case DirectionWest:
			buff += "W";
			break;
	}
	return buff;
}



func (this *RobotRepair) IsVisited(node *UnexploredRepairNode) bool{
	return this.Grid.HasValue(node.NewPosition.X, node.NewPosition.Y);
}

func (this *RobotRepair) GetAllDirections(vec2 *IntVec2) []*UnexploredRepairNode {
	res := make([]*UnexploredRepairNode, 4);

	north := &UnexploredRepairNode{};
	north.NewPosition = &IntVec2{};
	north.NewPosition.X = vec2.X;
	north.NewPosition.Y = vec2.Y-1;
	north.Direction = DirectionNorth;
	north.StartPosition = vec2.Clone();
	res[0] = north;

	south := &UnexploredRepairNode{};
	south.NewPosition = &IntVec2{};
	south.NewPosition.X = vec2.X;
	south.NewPosition.Y = vec2.Y+1;
	south.Direction = DirectionSouth;
	south.StartPosition = vec2.Clone();
	res[1] = south;

	west := &UnexploredRepairNode{};
	west.NewPosition = &IntVec2{};
	west.NewPosition.X = vec2.X-1;
	west.NewPosition.Y = vec2.Y;
	west.Direction = DirectionWest;
	west.StartPosition = vec2.Clone();
	res[2] = west;

	east := &UnexploredRepairNode{};
	east.NewPosition = &IntVec2{};
	east.NewPosition.X = vec2.X+1;
	east.NewPosition.Y = vec2.Y;
	east.Direction = DirectionEast;
	east.StartPosition = vec2.Clone();
	res[3] = east;

	return res;
}



func (this *RobotRepair) PrintState(){
	Log.Info("[Step %d] Repair Robot is at %s", this.StepCount, this.Position.ToString());
	Log.Info("\n%s\n", this.PrintGrid());
}

func (this *RobotRepair) PrintGrid() string{
	xMin := this.Grid.MinRow();
	xMax := this.Grid.MaxRow();

	yMin := this.Grid.MinCol();
	yMax := this.Grid.MaxCol();

	buff := "";
	for j := yMin; j<= yMax; j++{
		for i := xMin; i<= xMax; i++{
			if(!this.Grid.HasValue(i, j)){
				buff += " ";
			} else{
				if(i == this.Position.X && j == this.Position.X){
					buff += "D";
				}
				val := this.Grid.GetValue(i, j);
				if(val == RepairStateWall){
					buff += "#"
				} else if(val == RepairStateOxygenSystem) {
					buff += "S";
				} else if(val == RepairStateO2) {
					buff += "O";
				} else{
					buff += ".";
				}
			}
		}
		buff += "\n";
	}
	return buff;
}

func (this *RobotRepair) SimulateOxygenFill() int{


	steps := 0;

	frontier := make([]*IntVec2, 0);
	frontier = append(frontier, this.SystemLocation);

	for{
		if(len(frontier) <= 0){
			break;
		}
		newFrontier := make([]*IntVec2, 0);
		for _, node := range frontier{
			this.Grid.SetValue(node.X, node.Y, RepairStateO2);
			edges := this.Grid.GenerateEdges(node);
			for _, edge := range edges{
				if(this.Grid.GetValue(edge.X, edge.Y) == RepairStateVisited){
					newFrontier = append(newFrontier, edge);
				}
			}
		}
		frontier = newFrontier;
		steps++;
	}
	return steps - 1;
}

func (this *RobotRepair) IsComplete() bool{
	return this.Processor.HasHalted;
}

