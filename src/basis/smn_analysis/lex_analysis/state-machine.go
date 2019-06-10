package lex_analysis

import (
	"fmt"
)

type OnNodeRead func(stateNode *StateNode, input InputItf) (isEnd bool, err error)

type StateNodeReader interface {
	PreRead(stateNode *StateNode, input InputItf) (isEnd bool, err error)
	Read(stateNode *StateNode, input InputItf) (isEnd bool, err error)
	Clean()
}

type InputItf interface {
}

type ProductItf interface {
}

type StateNode struct {
	sm     *StateMachine
	Result ProductItf
	reader StateNodeReader
}

func (this *StateNode) Init(sm *StateMachine, SNReader StateNodeReader) *StateNode {
	this.reader = SNReader
	this.sm = sm
	return this
}

// !!!Warning!!!
//should not deal the input in this function.
//only check if end.
func (this *StateNode) PreRead(input InputItf) (isEnd bool, err error) {
	return this.reader.PreRead(this, input)
}

func (this *StateNode) Read(input InputItf) (isEnd bool, err error) {
	return this.reader.PreRead(this, input)
}

func (this *StateNode) CleanReader() {
	this.reader.Clean()
}

func (this *StateNode) ChangeStateNode(nextNode *StateNode) {
	this.sm.changeStateNode(nextNode)
}

type StateMachine struct {
	ChanSize     int
	resultChan   chan ProductItf
	nowStateNode *StateNode
	//when a StateNode end, will let nowStateNOde = DfeStateNode.
	//StateNode's PreRead should always return (isEnd=false, err=nil)
	DftStateNode *StateNode
}

func (this *StateMachine) Read(input InputItf) error {
	isEnd, err := this.nowStateNode.PreRead(input)
	if iserr(err) {
		return err
	}
	if isEnd {
		this.nowStateNode = this.DftStateNode
	}
	isEnd, err = this.nowStateNode.Read(input)
	if iserr(err) {
		return err
	}
	if isEnd {
		this.backDefault(this.DftStateNode)
	}
	return nil
}

func (this *StateMachine) Init() *StateMachine {
	if this.ChanSize <= 0 {
		this.ChanSize = 10000
	}
	this.resultChan = make(chan ProductItf, this.ChanSize)
	return this
}

func (this *StateMachine) changeStateNode(node *StateNode) {
	beforeNode := this.nowStateNode
	go func() {
		if beforeNode.Result != nil {
			this.resultChan <- beforeNode.Result
		}
	}() // should never block.
	this.nowStateNode = node
}

func (this *StateMachine) backDefault(node *StateNode) {
	this.DftStateNode.CleanReader()
	this.nowStateNode = this.DftStateNode
}

func (this *StateMachine) GetResultChan() <-chan ProductItf {
	return this.resultChan
}

type DftStateNodeReader struct {
	StateNodeReader
	SNodeList []*StateNode
	LiveMap   map[*StateNode]byte
	first     bool //is first call after clean
}

func NewDftStateNodeReader() *DftStateNodeReader {
	return &DftStateNodeReader{SNodeList: make([]*StateNode, 0), LiveMap: make(map[*StateNode]byte)}
}

func (this *DftStateNodeReader) Clean() {
	for k := range this.LiveMap {
		k.CleanReader()
		delete(this.LiveMap, k)
	}
	this.first = true
}

func (this *DftStateNodeReader) PreRead(stateNode *StateNode, input InputItf) (isEnd bool, err error) {
	return false, nil
}

func (this *DftStateNodeReader) Read(stateNode *StateNode, input InputItf) (isEnd bool, err error) {
	if this.first {
		for _, val := range this.SNodeList {
			this.LiveMap[val] = 0
		}
	}
	liveCnt := 0
	endCnt := 0
	errStr := ""
	var nextNode *StateNode
	for k := range this.LiveMap {
		kend, kerr := k.PreRead(input)
		if iserr(kerr) {
			errStr += kerr.Error()
		}
		if kend || iserr(kerr) {
			delete(this.LiveMap, k)
			continue
		}
		kend, kerr = k.Read(input)
		if iserr(kerr) {
			delete(this.LiveMap, k)
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
	if len(this.LiveMap) == 1 {
		stateNode.ChangeStateNode(nextNode)
		if endCnt == 1 {
			return true, nil
		}
	}
	if endCnt != 0 && len(this.LiveMap) != 1 { // todo: as success return?
		return true, fmt.Errorf(ErrTooMuchMatchStateNodeWhenHasEnd)
	}
	return false, nil
}
