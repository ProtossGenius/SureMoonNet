package smn_analysis

import (
	"fmt"
)

//OnNodeRead node read. maybe no use now?
type OnNodeRead func(stateNode *StateNode, input InputItf) (isEnd bool, err error)

//StateNodeReader state node reader.
type StateNodeReader interface {
	//Name reader's name.
	Name() string
	//PreRead only see if should stop read.
	PreRead(stateNode *StateNode, input InputItf) (isEnd bool, err error)
	//Read real read. even isEnd == true the input be readed.
	Read(stateNode *StateNode, input InputItf) (isEnd bool, err error)
	//GetProduct return result.
	GetProduct() ProductItf
	//Clean let the Reader like new.  it will be call before first Read.
	Clean()
}

//InputItf StateMachine's input.
type InputItf interface {
	//Copy should give copy to prevent user change it.
	Copy() InputItf
}

//ProductItf StateMachine's result.
type ProductItf interface {
	//ProductType result's type. usally should >= 0.
	ProductType() int
}

//ProductEnd last result.
type ProductEnd struct{}

//ProductType result's type. end's id = -1.
func (*ProductEnd) ProductType() int {
	return -1
}

//ProductDftNode get product type from default.
type ProductDftNode struct {
	Reason string
}

//ProductType .
func (p *ProductDftNode) ProductType() int {
	return -2
}

//StateNode StateMachine's node.
type StateNode struct {
	sm     *StateMachine
	Result ProductItf
	reader StateNodeReader
}

//Init .
func (sn *StateNode) Init(sm *StateMachine, SNReader StateNodeReader) *StateNode {
	sn.reader = SNReader
	sn.sm = sm
	return sn
}

//PreRead !!!Warning!!!
//should not deal the input in this function.
//only check if end.
func (sn *StateNode) PreRead(input InputItf) (isEnd bool, err error) {
	return sn.reader.PreRead(sn, input)
}

func (sn *StateNode) Read(input InputItf) (isEnd bool, err error) {
	return sn.reader.Read(sn, input)
}

//CleanReader run clean for reader.
func (sn *StateNode) CleanReader() {
	sn.Result = nil
	sn.reader.Clean()
}

//GetProduct get result.
func (sn *StateNode) GetProduct() {
	sn.Result = sn.reader.GetProduct()
}

//ChangeStateNode .
func (sn *StateNode) ChangeStateNode(nextNode *StateNode) {
	sn.sm.changeStateNode(nextNode)
}

//StateMachine state machine to formulate a state-tree and get result by input.
type StateMachine struct {
	ChanSize     int
	resultChan   chan ProductItf
	nowStateNode *StateNode
	//when a StateNode end, will let nowStateNOde = DfeStateNode.
	//StateNode's PreRead should always return (isEnd=false, err=nil)
	DftStateNode *StateNode
	//isPreadEnd pre-read's result.
	isPreadEnd bool
	//isReadend read's result.
	isReadEnd bool
}

func (sm *StateMachine) Read(input InputItf) error {
	var err error
	if sm.nowStateNode == nil {
		sm.useDefault()
	}
	sm.isPreadEnd, err = sm.nowStateNode.PreRead(input)

	if iserr(err) {
		return err
	}

	if sm.isPreadEnd {
		sm.useDefault() // now Reader is DftStateNodeReader, and it don't have PreRead so can call Read direct.
	}

	sm.isReadEnd, err = sm.nowStateNode.Read(input)
	if iserr(err) {
		return err
	}
	if sm.isReadEnd {
		sm.useDefault()
	}
	return nil
}

//IsPreadEnd return the last read status.
func (sm *StateMachine) IsPreadEnd() bool {
	return sm.isPreadEnd
}

//IsReadEnd is read end.
func (sm *StateMachine) IsReadEnd() bool {
	return sm.isReadEnd
}

//IsEndHappened if have end then return true.
func (sm *StateMachine) IsEndHappened() bool {
	return sm.IsPreadEnd() || sm.IsReadEnd()
}

//Init base on new(StateMachine).Init() write style.
func (sm *StateMachine) Init() *StateMachine {
	if sm.ChanSize <= 0 {
		sm.ChanSize = 10000
	}
	sm.resultChan = make(chan ProductItf, sm.ChanSize)
	return sm
}

func (sm *StateMachine) changeStateNode(node *StateNode) {
	if sm.nowStateNode != nil && sm.nowStateNode != sm.DftStateNode {
		beforeNode := sm.nowStateNode
		beforeNode.GetProduct()
		if beforeNode.Result != nil {
			sm.resultChan <- beforeNode.Result
		}
	}
	sm.nowStateNode = node
}

//End input end-input.
func (sm *StateMachine) End() {
	sm.nowStateNode.GetProduct()
	sm.resultChan <- sm.nowStateNode.Result
	if sm.nowStateNode.Result.ProductType() == -1 {
		return
	}
	sm.resultChan <- &ProductEnd{}
}

func (sm *StateMachine) useDefault() {
	sm.changeStateNode(sm.DftStateNode)
	sm.DftStateNode.CleanReader()
}

//GetResultChan .
func (sm *StateMachine) GetResultChan() <-chan ProductItf {
	return sm.resultChan
}

//DftStateNodeReader choice node.
type DftStateNodeReader struct {
	StateNodeReader
	sm        *StateMachine
	SNodeList []*StateNode
	LiveMap   map[*StateNode]byte
	first     bool //is first call after clean
}

//Name reader's name.
func (dsn *DftStateNodeReader) Name() string {
	return "DftStateNodeReader"
}

//GetProduct .
func (dsn *DftStateNodeReader) GetProduct() ProductItf {
	if len(dsn.LiveMap) == 0 {
		return &ProductEnd{}
	}
	res := &ProductDftNode{Reason: "Err when GetProduct ProductDftNode, MutiNode Lived they're  :"}
	for key := range dsn.LiveMap {
		res.Reason += ", " + key.reader.Name()
	}

	return res
}

//NewDftStateNodeReader .
func NewDftStateNodeReader(machine *StateMachine) *DftStateNodeReader {
	dft := &DftStateNodeReader{sm: machine, SNodeList: make([]*StateNode, 0), LiveMap: make(map[*StateNode]byte)}
	machine.DftStateNode = (&StateNode{}).Init(machine, dft)
	return dft
}

//Register .
func (dsn *DftStateNodeReader) Register(node StateNodeReader) {
	dsn.SNodeList = append(dsn.SNodeList, (&StateNode{}).Init(dsn.sm, node))
}

//Clean .
func (dsn *DftStateNodeReader) Clean() {
	for k := range dsn.LiveMap {
		delete(dsn.LiveMap, k)
	}
	for _, node := range dsn.SNodeList {
		node.CleanReader()
	}
	dsn.first = true
}

//PreRead DftStateNodeReader don't do PreRead, because it should deal all registed reader.
func (dsn *DftStateNodeReader) PreRead(stateNode *StateNode, input InputItf) (isEnd bool, err error) {
	return false, nil
}

func (dsn *DftStateNodeReader) Read(stateNode *StateNode, input InputItf) (isEnd bool, err error) {
	if dsn.first {
		for _, val := range dsn.SNodeList {
			dsn.LiveMap[val] = 0
		}
		dsn.first = false
	}
	liveCnt := 0
	endCnt := 0
	errStr := ""
	var nextNode *StateNode
	for k := range dsn.LiveMap {
		kend, kerr := k.PreRead(input.Copy())
		if iserr(kerr) {
			errStr += kerr.Error()
		}
		if kend || iserr(kerr) {
			delete(dsn.LiveMap, k)
			continue
		}
		kend, kerr = k.Read(input.Copy())
		if iserr(kerr) {
			delete(dsn.LiveMap, k)
			continue
		}
		liveCnt++
		if liveCnt == 1 {
			nextNode = k
		}
		if kend {
			endCnt++
			nextNode = k
		}
	}
	if liveCnt == 0 {
		return true, fmt.Errorf(ErrNoMatchStateNode, errStr)
	}
	if endCnt > 1 {
		return true, fmt.Errorf(ErrTooMuchMatchStateNode)
	}
	if len(dsn.LiveMap) == 1 {
		stateNode.ChangeStateNode(nextNode)
		if endCnt == 1 {
			return true, nil
		}
	}
	if endCnt != 0 && len(dsn.LiveMap) != 1 { // todo: as success return?
		return true, fmt.Errorf(ErrTooMuchMatchStateNodeWhenHasEnd)
	}
	return false, nil
}
