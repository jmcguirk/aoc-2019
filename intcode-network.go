package main

import "fmt"

const NATNetworkAddress = 255;

type IntcodeNetwork struct {
	Terminals []*IntcodeTerminal;
	LastAddress int;
	SimulationStep int;
	LastPacketId int;
	PendingPackets []*IntcodePacket;
	EnableNAT bool;
	LastNATPacket *IntcodePacket
	LastFlushedNATPacket *IntcodePacket
	TotalPacketsDelivered int;
	AddressOfInterest int;
}

func (this *IntcodeNetwork) Init(){
	this.PendingPackets = make([]*IntcodePacket, 0);
	this.Terminals = make([]*IntcodeTerminal, 0);
	this.LastAddress = 0;
	this.AddressOfInterest = -1;
}

func (this *IntcodeNetwork) AddTerminals(instructionFileName string, count int) error{
	for i := 0; i < count; i++{
		terminal := &IntcodeTerminal{}
		err := terminal.Init(this.LastAddress, instructionFileName, this);
		if(err != nil){
			return err;
		}
		this.Terminals = append(this.Terminals, terminal);
		this.LastAddress++;
	}
	Log.Info("Added %d terminals to network using file %s", count, instructionFileName)
	return nil;
}

func (this *IntcodeNetwork) Simulate() error{
	for {

		// Doing this in two loops - not sure if actually needed or if I should be queuing immediate?
		for _, terminal := range this.Terminals{
			p, err := terminal.ReceiveAndSend();
			if(err != nil){
				return err;
			}
			if(p != nil){

				p.PacketId = this.LastPacketId;
				p.GeneratedOnFrame = this.SimulationStep;
				this.LastPacketId++;
				if(terminal.Address == 5){
					//Log.Info("Terminal %d sending %s", terminal.Address, p.Describe());
				}
				//this.PendingPackets = append(this.PendingPackets, packet);
				if(this.EnableNAT && p.ToAddress == NATNetworkAddress){
					//Log.Info("Received NAT packet %s", p.Describe());
					this.LastNATPacket = p;
					continue;
				}
				if(p.ToAddress >= len(this.Terminals) || p.ToAddress < 0){
					continue;
				}
				deliveryAddress := this.Terminals[p.ToAddress];
				if(p.ToAddress == this.AddressOfInterest) {
					Log.Info("Found packet to address of interest %s ", p.Describe());
					return nil;
				}
				deliveryAddress.Deliver(p);
				this.TotalPacketsDelivered++;
			}
		}
		//
		//for _, p := range this.PendingPackets{

		//}
		this.PendingPackets = nil;

		if(this.EnableNAT && this.LastNATPacket != nil){
			isIdle := true;
			for _, terminal := range this.Terminals{
				if(!terminal.IsIdle()){
					//Log.Info("Terminal %d is still active %d %d", terminal.Address, len(terminal.Inbox), terminal.IdleFrames);
					isIdle = false;
					break;
				}
			}
			if(isIdle){
				//Log.Info("Network determined idle on frame %d, sending packet %s to address 0", this.SimulationStep, this.LastNATPacket.Describe());
				if(this.LastFlushedNATPacket != nil && this.LastFlushedNATPacket.Y == this.LastNATPacket.Y){
					Log.Info("Dupe value received %d", this.LastNATPacket.Y);
					return nil;
				}
				this.Terminals[0].Deliver(this.LastNATPacket);
				this.LastFlushedNATPacket = this.LastNATPacket;
			}
		}

		this.SimulationStep++;
		//Log.Info("Step %d delivered %d packets", this.SimulationStep, this.TotalPacketsDelivered);
	}
	return nil;
}



type IntcodeTerminal struct {
	Address int;
	Processor *IntcodeMachineV3;
	ContainingNetwork *IntcodeNetwork;
	InstructionFile string;
	Inbox []*IntcodePacket;
	HasHalted bool;
	IdleFrames int;
}

type IntcodePacket struct {
	FromAddress int;
	ToAddress int;
	GeneratedOnFrame int;
	PacketId int;
	X int64;
	Y int64;
}

func (this *IntcodePacket) Describe() string{
	return fmt.Sprintf("Packet %d - From: %d, To: %d, X: %d, Y: %d", this.PacketId, this.FromAddress, this.ToAddress, this.X, this.Y);
}

func (this *IntcodeTerminal) Init(address int, instructionFileName string, network *IntcodeNetwork) error{
	this.ContainingNetwork = network;
	this.Processor = &IntcodeMachineV3{};
	this.InstructionFile = instructionFileName;
	err := this.Processor.Load(instructionFileName);
	if(err != nil){
		return err;
	}
	//this.Processor.PauseOnDefaultInput = true;
	//this.Processor.SetDefaultInput(-1);
	this.Processor.QueueInput(int64(address));
	this.Processor.PauseOnOutput = true;
	this.Processor.PauseOnInput = true;
	this.Address = address;
	this.Inbox = make([]*IntcodePacket, 0);
	return nil;
}

func (this *IntcodeTerminal) Deliver(packet *IntcodePacket){
	//this.Processor.QueueInput(packet.X);
	//this.Processor.QueueInput(packet.Y);
	this.Inbox = append(this.Inbox, packet);
}

func (this *IntcodeTerminal) IsIdle() bool{
	return len(this.Inbox) <= 0;
}

func (this *IntcodeTerminal) ReceiveAndSend() (*IntcodePacket, error){
	if(len(this.Inbox) > 0) {
		for len(this.Inbox) > 0 {
			packet := this.Inbox[0];
			this.Processor.QueueInput(packet.X);
			this.Processor.QueueInput(packet.Y);
			this.Inbox = this.Inbox[1:];
		}
	}



	val, err, hasHalted := this.Processor.ReadNextOutput();
	if(this.Processor.PendingInput){
		this.Processor.QueueInput(-1);
		return nil, nil;
	}

	if(err != nil){
		return nil, err;
	}
	if(hasHalted){
		this.HasHalted = true;
		return nil, nil;
	}
	if(val == -1){
		this.IdleFrames++;
		return nil, nil;
	}
	address := val;
	val, err, hasHalted = this.Processor.ReadNextOutput();
	if(err != nil){
		return nil, err;
	}
	if(hasHalted){
		this.HasHalted = true;
		return nil, nil;
	}
	if(val == -1){
		return nil, nil;
	}
	x := val;
	val, err, hasHalted = this.Processor.ReadNextOutput();
	if(err != nil){
		return nil, err;
	}
	if(hasHalted){
		this.HasHalted = true;
		return nil, nil;
	}
	if(val == -1){
		return nil, nil;
	}
	y := val;
	packet := &IntcodePacket{};
	packet.FromAddress = this.Address;
	packet.ToAddress = int(address);
	packet.X = x;
	packet.Y = y;
	this.IdleFrames = 0;
	return packet, nil;
}