package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)


type TunnelSearchSystemMultiAgent struct {
	CanonicalGrid *IntegerGrid2D;
	InstructionFileName string;
	Partitions []*TunnelPartition
	Doors map[int]*TunnelSearchDoor;
	AllKeys map[int]*TunnelSearchKey;
	BestPath int;
}

type TunnelPartition struct {
	StartPos *IntVec2;
	Doors map[int]*TunnelSearchDoor;
	AllKeys map[int]*TunnelSearchKey;
	PartitionId int;
	StartingPaths []*TunnelSearchKeyPairPath;
}



type TunnelAgentState struct {
	Partition *TunnelPartition;
	CurrKey int;
}

func (this *TunnelAgentState) Clone() *TunnelAgentState {
	res := &TunnelAgentState{};
	res.Partition = this.Partition;
	res.CurrKey = this.CurrKey;
	return res;
}

type TunnelGlobalState struct {
	Agents []*TunnelAgentState;
	HeldKeys []int;
	Depth int;
	TotalCost int;
}


func (this *TunnelGlobalState) Clone() *TunnelGlobalState {
	res := &TunnelGlobalState{};
	res.TotalCost = this.TotalCost;
	res.HeldKeys = make([]int, 0);
	res.Depth = this.Depth;
	res.HeldKeys = append(res.HeldKeys, this.HeldKeys...);
	res.Agents = make([]*TunnelAgentState, len(this.Agents));
	for i, a := range this.Agents{
		res.Agents[i] = a.Clone();
	}
	return res;
}

func (this *TunnelSearchSystemMultiAgent) Init(gridFile string) error {

	this.CanonicalGrid = &IntegerGrid2D{};
	this.CanonicalGrid.Init();

	this.AllKeys = make(map[int]*TunnelSearchKey);
	this.Doors = make(map[int]*TunnelSearchDoor);

	this.Partitions = make([]*TunnelPartition, 0);

	file, err := os.Open(gridFile);
	if err != nil {
		Log.FatalError(err);
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)


	y := 0;
	x := 0;



	for scanner.Scan() {
		lineRaw := strings.TrimSpace(scanner.Text());
		if (lineRaw != "") {
			for _, c := range lineRaw {
				val := int(c);
				this.CanonicalGrid.SetValue(x, y, val);
				if (val == TunnelStartCode) {
					partition := &TunnelPartition{};
					partition.PartitionId = len(this.Partitions);
					partition.StartPos = &IntVec2{};
					partition.StartPos.X = x;
					partition.StartPos.Y = y;
					this.Partitions = append(this.Partitions, partition);
					//this.StartPos = &IntVec2{};
					//this.StartPos.X = x;
					//this.StartPos.Y = y;
				} else if (val >= TunnelDoorStart && val <= TunnelDoorEnd) {

					door := &TunnelSearchDoor{}
					door.X = x;
					door.Y = y;
					door.KeyCode = val;
					this.Doors[val] = door;
				} else if (val >= TunnelKeyStart && val <= TunnelKeyEnd) {
					key := &TunnelSearchKey{};
					key.X = x;
					key.Y = y;
					key.StartPos = &IntVec2{};
					key.StartPos.X = x;
					key.StartPos.Y = y;
					key.KeyCode = val;
					this.AllKeys[val] = key;
					//this.KeyCount++;
				}
				x++;
			}
		}
		x = 0;
		y++;
	}



	Log.Info("Parse finished - %d keys, %d doors  in %d partitions", len(this.AllKeys), len(this.Doors), len(this.Partitions));

	for _, partition := range this.Partitions{
		partition.AllKeys = make(map[int]*TunnelSearchKey, 0);
		partition.Doors = make(map[int]*TunnelSearchDoor, 0);
		partition.StartingPaths = make([]*TunnelSearchKeyPairPath, 0);
		reachablePoints := this.CanonicalGrid.Reachable(partition.StartPos, TunnelWallCode);
		for _, v2 := range reachablePoints{
			val := this.CanonicalGrid.GetValue(v2.X, v2.Y);
			k, exists := this.AllKeys[val];
			if(exists){
				partition.AllKeys[k.KeyCode] = k;
			}
		}
	}

	for _, partition := range this.Partitions{

		//Log.Info("Partition %d has %d keys", i, len(partition.AllKeys));

		blacklist := make([]int, 0);
		blacklist = append(blacklist, TunnelWallCode);


		for _, key := range partition.AllKeys{
			path := this.CanonicalGrid.ShortestPathWithBlacklist(partition.StartPos, key.StartPos, blacklist);
			if(path != nil){

				p := &TunnelSearchKeyPairPath{};
				p.Key = key;
				p.BestPath = path;
				for _, node := range path {
					val := this.CanonicalGrid.GetValue(node.X, node.Y);
					_, exists := this.Doors[val];
					if(exists){
						p.KeysRequired = append(p.KeysRequired, val + TunnelKeyCaseOffset);
					}
				}
				partition.StartingPaths = append(partition.StartingPaths, p);
				//Log.Info("origin to %c requires %d keys", key.KeyCode, len(p.KeysRequired));
			}
			//Log.Info("Partition %d had %d initially reachable keys", i, len(partition.StartingPaths));
		}




		blacklist = make([]int, 0);
		blacklist = append(blacklist, TunnelWallCode);

		for _, key := range partition.AllKeys{
			key.Paths = make([]*TunnelSearchKeyPairPath, 0);
			for _, key2 := range partition.AllKeys{
				if(key2.KeyCode == key.KeyCode){
					continue;
				}
				path := this.CanonicalGrid.ShortestPathWithBlacklist(key.StartPos, key2.StartPos, blacklist);
				if(path != nil){
					pair := &TunnelSearchKeyPairPath{};
					pair.BestPath = path;
					pair.KeysRequired = make([]int, 0);
					for _, node := range path {
						val := this.CanonicalGrid.GetValue(node.X, node.Y);
						_, exists := this.Doors[val];
						if(exists){
							pair.KeysRequired = append(pair.KeysRequired, val + TunnelKeyCaseOffset);
						}
					}
					pair.Key = key2;
					key.Paths = append(key.Paths, pair);
					//Log.Info("%c to %c requires %d keys", key.KeyCode, key2.KeyCode, len(pair.KeysRequired));

				}
			}
			sort.SliceStable(key.Paths, func(i, j int) bool {
				return len(key.Paths[i].BestPath) < len(key.Paths[j].BestPath);
			});

		}


	}

	agentOrder := make([]int, 4);
	agentOrder[0] = 0;
	agentOrder[1] = 1;
	agentOrder[2] = 2;
	agentOrder[3] = 3;

	permutations := make([][]int, 0);


	PermInt(agentOrder, func(a []int) {
		cpy := make([]int, len(a));
		copy(cpy, a);
		permutations = append(permutations, cpy);
	})

	state := &TunnelGlobalState{};
	state.Agents = make([]*TunnelAgentState, 4);
	for i, partition := range this.Partitions{
		p := &TunnelAgentState{};
		p.Partition = partition;
		state.Agents[i] = p;
	}

	state.Depth = 0;
	state.HeldKeys = make([]int, 0);

	this.BestPath = 2104;
	this.FindShortestPaths(state, permutations);



	return nil;
}


func (this *TunnelSearchSystemMultiAgent) GetRelevantPaths(agent *TunnelAgentState) []*TunnelSearchKeyPairPath {
	if(agent.CurrKey > 0){
		return agent.Partition.AllKeys[agent.CurrKey].Paths;
	}
	return agent.Partition.StartingPaths;
}

func (this *TunnelSearchSystemMultiAgent) FindShortestPaths(state *TunnelGlobalState, permutations [][]int) {

	if(len(state.HeldKeys) == len(this.AllKeys)){
		if(this.BestPath > state.TotalCost){
			Log.Info("New best path found %d", state.TotalCost)
			state.TotalCost = this.BestPath;
		}
		return;
	}
	if(state.Depth == 0){
		operationDone := make(chan bool)
		for _, order := range permutations {
			go this.FindShortestPathsWithOrder(state, permutations, order);
		}
		<-operationDone
	} else{
		for _, order := range permutations {
			this.FindShortestPathsWithOrder(state, permutations, order);
		}
	}

}

func (this *TunnelSearchSystemMultiAgent) FindShortestPathsWithOrder(state *TunnelGlobalState, permutations [][]int, order[]int) {

	if(len(state.HeldKeys) == len(this.AllKeys)){
		if(this.BestPath > state.TotalCost){
			Log.Info("New best path found %d", state.TotalCost)
			state.TotalCost = this.BestPath;
		}
		return;
	}
	for _, i := range order{


		agent := state.Agents[i];
		possiblePaths := this.GetRelevantPaths(agent);

		for _, path := range possiblePaths {
			subState := state.Clone();
			subState.Depth++;
			agent = subState.Agents[i];
			//Log.Info("Considering path to %c for agent %d with path %s", path.Key.KeyCode, i, this.PrintAscii(state.HeldKeys));
			//path := reachableKeys[key];
			newCost := subState.TotalCost + len(path.BestPath);
			if (newCost >= this.BestPath) {
				continue;
			}
			dupe := false;
			for _,v:= range subState.HeldKeys{
				if(v == path.Key.KeyCode){
					dupe = true;
					break;
				}
			}
			if(dupe){
				continue;
			}
			locked := false;
			for _, k:= range path.KeysRequired{
				hasKey := false;
				for _,v:= range subState.HeldKeys{
					if(v == k){
						hasKey = true;
						break;
					}
				}
				if(!hasKey){
					locked = true;
					break;
				}
			}
			if(locked){
				continue;
			}
			subState.TotalCost = newCost;

			subState.HeldKeys = append(subState.HeldKeys, path.Key.KeyCode);
			agent.CurrKey = path.Key.KeyCode;
			//Log.Info("Considering %d %s after acquiring %c", newCost, this.PrintAscii(subState.HeldKeys), path.Key.KeyCode);
			if (len(subState.HeldKeys) == len(this.AllKeys)) {
				//Log.Info("Found path %d - %s", newCost, this.PrintAscii(subState.HeldKeys));
				if (newCost < this.BestPath) {
					Log.Info("Found best path %d - %s", newCost, this.PrintAscii(subState.HeldKeys));
					this.BestPath = newCost;
					return;
				}
				continue;
			}
			this.FindShortestPaths(subState, permutations);
		}

	}
}

func (this *TunnelSearchSystemMultiAgent) PrintAscii(res[] int) string {
	buff := "";
	for i, v := range res{
		if(i > 0){
			buff += ",";
		}
		buff += fmt.Sprintf("%c", v);
	}

	return buff;
}