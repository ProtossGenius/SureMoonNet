package smn_analysis_go

type AnalysisItf interface {
	Read(info UnitInfo)
}

type OnNodeRead func(unitInfo UnitInfo, sn *StateNode, sm *StateMachine) error

type UnitInfo struct {
	UnitType string
	UnitVal  string
}

type BlockInfo interface {
	Type() string
}

type StateNode struct {
	sm     *StateMachine
	onread OnNodeRead
	Block  BlockInfo
}

func (this *StateNode) Read(info UnitInfo) {
	this.onread(info, this, this.sm)
}

type StateMachine struct {
	ChanSize     int
	infoChan     chan BlockInfo
	nowStateNode *StateNode
}

func (this *StateMachine) Read(info UnitInfo) {
	this.nowStateNode.Read(info)
}

func (this *StateMachine) Init() *StateMachine {
	if this.ChanSize <= 0 {
		this.ChanSize = 10000
	}
	this.infoChan = make(chan BlockInfo, this.ChanSize)
	return this
}

func (this *StateMachine) stateChange(node *StateNode) {
	go func() { this.infoChan <- node.Block }() // should never block.
	this.nowStateNode = node
}

func (this *StateMachine) GetInfoChan() <-chan BlockInfo {
	return this.infoChan
}
