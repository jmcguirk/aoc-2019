package main

type DirectedGraph struct {
	LastNodeId 		int;
	LastEdgeId 		int;
	Nodes  map[int]*Node;
	Edges  map[int]*Edge;
	LabelToNode map[string]*Node;
	LabelToEdge map[string]*Edge;
}

func (this *DirectedGraph) Init() {
	this.LastNodeId = 1;
	this.Nodes = make(map[int]*Node);
	this.Edges = make(map[int]*Edge);
	this.LabelToNode = make(map[string]*Node);
}

func (this *DirectedGraph) AllNodes() []*Node {
	res := make([]*Node, 0);
	for _, node := range this.Nodes{
		res = append(res, node);
	}
	return res;
}



func (this *DirectedGraph) GetOrCreateNode(label string)*Node {
	res, exists := this.LabelToNode[label];
	if(exists){
		return res;
	}
	res = &Node{};
	res.Init(this);
	res.Label = label;
	res.Id = this.LastNodeId;
	this.LastNodeId++;
	this.LabelToNode[res.Label] = res;
	this.Nodes[res.Id] = res;
	return res;
}

func (this *DirectedGraph) CreateEdge(from *Node, to *Node)*Edge {
	res := &Edge{};
	res.From = from;
	res.To = to;
	from.Edges = append(from.Edges, res);
	res.Id = this.LastEdgeId;
	this.LastEdgeId++;
	this.Edges[res.Id] = res;
	return res;
}