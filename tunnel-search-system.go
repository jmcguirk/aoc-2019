package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
)

const TunnelWallCode = int('#');
const TunnelEmptyCode = int('.');
const TunnelStartCode = int('@');
const TunnelDoorStart = int('A');
const TunnelDoorEnd = int('Z');
const TunnelKeyStart = int('a');
const TunnelKeyEnd = int('z');
const TunnelKeyCaseOffset = int('a') - int('A');

type TunnelSearchSystem struct {
	CanonicalGrid *IntegerGrid2D;
	InstructionFileName string;
	StartPos *IntVec2;

	Doors map[int]*TunnelSearchDoor;
	AllKeys map[int]*TunnelSearchKey;
	KeyCount int;
	BestPath int;
}

type TunnelSearchNode struct {
	X	int;
	Y 	int;
	KeysHeld []int;
	Id string;
	NeighborStates []*TunnelSearchNode;
	KeyUnlocked int;
}


func (this *TunnelSearchNode) CanTraversDoor(door int) bool {
	for _, k := range this.KeysHeld{
		if(k - TunnelKeyCaseOffset == door){
			return true;
		}
	}
	return false;
}

func (this *TunnelSearchNode) GenerateId() string{
	this.Id = fmt.Sprintf("%d,%d,%c,%s", this.X,this.Y, this.KeyUnlocked,this.KeySignature());
	return this.Id;
}

func (this *TunnelSearchNode) KeySignature() string{
	buff := "";
	for _, c :=  range this.KeysHeld{
		buff += fmt.Sprintf("%c", c);
	}
	return buff;
}

func (this *TunnelSearchNode) HasKey(key int) bool {
	for _, k := range this.KeysHeld{
		if(k == key){
			return true;
		}
	}
	return false;
}


type TunnelSearchKey struct {
	X	int;
	Y 	int;
	StartPos *IntVec2;
	UnlocksDoor bool;
	KeyCode int;
	Paths []*TunnelSearchKeyPairPath;
}

type TunnelSearchKeyPairPath struct {
	Key *TunnelSearchKey;
	KeysRequired []int;
	BestPath []*IntVec2;
}


type TunnelSearchDoor struct {
	X	int;
	Y 	int;
	KeyCode int;
}

func (this *TunnelSearchSystem) Init(gridFile string) error {

	this.CanonicalGrid = &IntegerGrid2D{};
	this.CanonicalGrid.Init();

	this.AllKeys = make(map[int]*TunnelSearchKey);
	this.Doors = make(map[int]*TunnelSearchDoor);

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
		if(lineRaw != ""){
			for _, c := range lineRaw{
				val := int(c);
				this.CanonicalGrid.SetValue(x, y, val);
				if(val == TunnelStartCode){
					this.StartPos = &IntVec2{};
					this.StartPos.X = x;
					this.StartPos.Y = y;
				} else if(val >= TunnelDoorStart && val <= TunnelDoorEnd){

					door := &TunnelSearchDoor{}
					door.X = x;
					door.Y = y;
					door.KeyCode = val;
					this.Doors[val] = door;
				} else if(val >= TunnelKeyStart && val <= TunnelKeyEnd){
					key := &TunnelSearchKey{};
					key.X = x;
					key.Y = y;
					key.StartPos = &IntVec2{};
					key.StartPos.X = x;
					key.StartPos.Y = y;
					key.KeyCode = val;
					this.AllKeys[val] = key;
					this.KeyCount++;
				}
				x++;
			}
		}
		x = 0;
		y ++;
	}

	Log.Info("Parse finished - %d keys, %d doors - start position is %d, %d", len(this.AllKeys), len(this.Doors), this.StartPos.X, this.StartPos.Y);

	for _, k := range this.AllKeys{
		_, exists := this.Doors[k.KeyCode-TunnelKeyCaseOffset];
		k.UnlocksDoor = exists;
	}


	return nil;
}


func (this *TunnelSearchSystem) GenerateId(x int, y int, keyFound int, keys[] int) string{
	return fmt.Sprintf("%d,%d,%c,%s", x, y, keyFound, keys);
}

func (this *TunnelSearchSystem) FindShortestPath() int {

	initiallyReachableKeys := make(map[int][]*IntVec2);

	blacklist := make([]int, 0);
	blacklist = append(blacklist, TunnelWallCode);
	for _, door := range this.Doors{
		blacklist = append(blacklist, door.KeyCode);
	}

	for _, key := range this.AllKeys{
		path := this.CanonicalGrid.ShortestPathWithBlacklist(this.StartPos, key.StartPos, blacklist);
		if(path != nil){
			initiallyReachableKeys[key.KeyCode] = path;
		}
	}

	blacklist = make([]int, 0);
	blacklist = append(blacklist, TunnelWallCode);

	for _, key := range this.AllKeys{
		key.Paths = make([]*TunnelSearchKeyPairPath, 0);
		for _, key2 := range this.AllKeys{
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

	Log.Info("Considering %d reachable keys", len(initiallyReachableKeys));

	//heldKeys := make([]int, 0);
	//this.BestPath = int(math.MaxInt64);
	operationDone := make(chan bool)
	this.BestPath = 5072; // LNG
	for key, path := range initiallyReachableKeys {
		newHeld := make([]int, 0);
		newHeld = append(newHeld, key);
		cost := len(path);
		go this.FindShortestPathWithKeys(this.AllKeys[key].Paths, newHeld, cost, 0);

	}
	<-operationDone

	Log.Info("Found shortest path %d steps", this.BestPath);

	return -1;
}

func (this *TunnelSearchSystem) PrintAscii(res[] int) string {
	buff := "";
	for i, v := range res{
		if(i > 0){
			buff += ",";
		}
		buff += fmt.Sprintf("%c", v);
	}

	return buff;
}

func (this *TunnelSearchSystem) FindShortestPathWithKeys(possiblePaths []*TunnelSearchKeyPairPath, HeldKeys[] int, totalPathCost int, depth int) {

	if(depth == 0){
		Log.Info("Starting %c %d", HeldKeys[0], totalPathCost);
	}
	if(len(HeldKeys) == len(this.AllKeys)){
		if(totalPathCost < this.BestPath){
			//Log.Info("Found path %d %d vs ", totalPathCost,len(this.AllKeys), len(HeldKeys));
			this.BestPath = totalPathCost;
			return;
		}
	}

	for _, path := range possiblePaths {
		//path := reachableKeys[key];
		newCost := totalPathCost + len(path.BestPath);
		if (newCost >= this.BestPath) {
			continue;
		}
		newHeld := make([]int, 0);
		newHeld = append(newHeld, HeldKeys...);
		newHeld = append(newHeld, path.Key.KeyCode);
		//Log.Info("Considering %d %s", newCost, this.PrintAscii(newHeld));
		if (len(newHeld) == len(this.AllKeys)) {
			if (newCost < this.BestPath) {
				Log.Info("Found path %d - %s", newCost, this.PrintAscii(newHeld));
				this.BestPath = newCost;
				continue;
			}
		}

		newReachable := make([]*TunnelSearchKeyPairPath, 0);
		canonicalKey := path.Key;
		for _, nextPath := range canonicalKey.Paths {
			if (newCost+len(nextPath.BestPath) >= this.BestPath) {
				continue;
			}
			dupe := false;
			for _, v := range newHeld { // Already held, don't bother
				if (nextPath.Key.KeyCode == v) {
					dupe = true;
					break;
				}
			}
			if (dupe) {
				continue;
			}
			locked := false;
			for _, v := range nextPath.KeysRequired { // Already held, don't bother
				hasKey := false;
				for _, h := range newHeld { // Already held, don't bother
					if (h == v) {
						hasKey = true;
						break;
					}
				}
				if (!hasKey) {
					locked = true;
					break;
				}
			}
			if (!locked) {
				newReachable = append(newReachable, nextPath);
			}

		}

		this.FindShortestPathWithKeys(newReachable, newHeld, newCost, depth+1);
	}
	if(depth == 0){
		Log.Info("Completed %c", HeldKeys[0]);
	}
}

func (this *TunnelSearchSystem) GenerateChildNode(next *TunnelSearchNode, xOffset int, yOffset int) *TunnelSearchNode{
	node := &TunnelSearchNode{}
	node.X = next.X + xOffset;
	node.Y = next.Y + yOffset;
	nVal := this.CanonicalGrid.GetValue(node.X, node.Y);
	//Log.Info("Generating %d, %d - %c", node.X, node.Y, nVal);
	node.KeysHeld = make([]int, 0);
	node.KeysHeld = append(node.KeysHeld, next.KeysHeld...);
	valid := false;
	if(nVal != TunnelWallCode){
		if(nVal >= TunnelDoorStart && nVal <= TunnelDoorEnd){
			if(next.CanTraversDoor(nVal)){
				//Log.Info("walking through door!");
				valid = true;
			}
		} else if(nVal >= TunnelKeyStart && nVal <= TunnelKeyEnd) {
			valid = true;
			if(!node.HasKey(nVal)){
				//Log.Info("Picked up key! %c", nVal);
				node.KeysHeld = append(node.KeysHeld, nVal);
				node.KeyUnlocked = nVal;
			}
		} else{
			valid = true;
		}
	}
	if(!valid){
		return nil;
	}
	node.GenerateId();
	return node;
}

func (this *TunnelSearchSystem) FindShortestPathToKey(key *TunnelSearchKey) []*TunnelSearchNode{
	Log.Info("Searching from start pos %d,%d to key %c at %d,%d", this.StartPos.X, this.StartPos.Y, key.KeyCode, key.X, key.Y);


	allNodes := make(map[string]*TunnelSearchNode);

	start := &TunnelSearchNode{};
	start.KeysHeld = make([]int, 0);
	start.X = this.StartPos.X;
	start.Y = this.StartPos.Y;
	start.Id = start.GenerateId();
	allNodes[start.Id] = start;

	res := make([]*TunnelSearchNode, 0);

	visitedNodes := make(map[string]*TunnelSearchNode);
	minCostToStart := make(map[string]int);
	nearestToStart := make(map[string]*TunnelSearchNode);

	frontier := make([]*TunnelSearchNode, 0);
	frontier = append(frontier, start);
	frontierMap := make(map[string]*TunnelSearchNode);
	frontierMap[start.Id] = start;
	minCostToStart[start.Id] = 0;

	var endNode *TunnelSearchNode
	for {
		if (len(frontier) <= 0) {
			break;
		}
		sort.SliceStable(frontier, func(i, j int) bool {
			return minCostToStart[frontier[i].Id] < minCostToStart[frontier[j].Id];
		});


		next := frontier[0];
		frontier = frontier[1:];
		//Log.Info("exploring %s", next.Id);
		delete(frontierMap, next.Id);
		costToHere := minCostToStart[next.Id];

		if(next.NeighborStates == nil){
			next.NeighborStates = make([]*TunnelSearchNode, 0);

			// North
			north := this.GenerateChildNode(next, 0, -1);
			if(north != nil){
				existing, exists := allNodes[north.Id];
				if(exists){
					north = existing;
				} else{
					allNodes[north.Id] = north;
				}
				next.NeighborStates = append(next.NeighborStates, north);
			}

			// South
			south := this.GenerateChildNode(next, 0, +1);
			if(south != nil){
				existing, exists := allNodes[south.Id];
				if(exists){
					south = existing;
				} else{
					allNodes[south.Id] = south;
				}
				next.NeighborStates = append(next.NeighborStates, south);
			}

			// East
			east := this.GenerateChildNode(next, +1, 0);
			if(east != nil){
				existing, exists := allNodes[east.Id];
				if(exists){
					east = existing;
				} else{
					allNodes[east.Id] = east;
				}
				next.NeighborStates = append(next.NeighborStates, east);
			}

			// West
			west := this.GenerateChildNode(next, -1, 0);
			if(west != nil){
				existing, exists := allNodes[west.Id];
				if(exists){
					west = existing;
				} else{
					allNodes[west.Id] = west;
				}
				next.NeighborStates = append(next.NeighborStates, west);
			}

		}




		for _, edge := range next.NeighborStates{
			neighbor := edge;
			_, visited := visitedNodes[neighbor.Id];
			if(visited){
				continue;
			}

			bestToHere, bestCostExists := minCostToStart[neighbor.Id];
			if(!bestCostExists){
				bestToHere = int(math.MaxInt32);
			}

			score := costToHere + 1000*(len(this.AllKeys) - len(neighbor.KeysHeld));
			if(score < bestToHere){
				minCostToStart[neighbor.Id] = score;
				nearestToStart[neighbor.Id] = next;
				_, alreadyEnqueued := frontierMap[neighbor.Id];
				if(!alreadyEnqueued){
					//Log.Info("Enqueueing %s", neighbor.Id);
					frontierMap[neighbor.Id] = neighbor;
					frontier = append(frontier, neighbor);
				}
			}

		}
		visitedNodes[next.Id] = next;
		if(next.X == key.X && next.Y == key.Y && len(next.KeysHeld) == len(this.AllKeys)){
			//Log.Info("Path completed successfully %s", next.Id);
			endNode = next;
			break;
		}
	}

	if(endNode == nil){
		Log.Info("Path from start pos %d,%d to key %c at %d,%d was not possible", this.StartPos.X, this.StartPos.Y, key.KeyCode, key.X, key.Y);
		return nil;
	}

	_, exists := minCostToStart[endNode.Id];
	if(!exists){
		Log.Info("Path from start pos %d,%d to key %c at %d,%d was not possible", this.StartPos.X, this.StartPos.Y, key.KeyCode, key.X, key.Y);
		return nil; // No path found
	}

	nextPathStep := endNode.Id;

	for {
		next := nearestToStart[nextPathStep];
		if(next == start){
			break;
		}
		nextPathStep = next.Id;
		res = append(res, next);
	}
	res = append(res, endNode);

	//ReverseSlice(res);
	return res;
}
