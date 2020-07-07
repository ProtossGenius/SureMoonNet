package smn_analysis

import (
	"fmt"
	"strings"
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

const (
	//ResultEnd product by end.
	ResultEnd = -1
	//ResultPFromDft product from default node.
	ResultPFromDft = -2
	//ResultError end with error.
	ResultError = -3
)

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

//ProductError end with error.
type ProductError struct {
	Error string
}

//ProductType result's type. usally should >= 0.
func (*ProductError) ProductType() int {
	return -3
}

//StateNode StateMachine's node.
type StateNode struct {
	sm     *StateMachine
	Result ProductItf
	reader StateNodeReader
	Datas  map[string]interface{}
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

//SendProduct send product to StateMachine.
func (sn *StateNode) SendProduct(result ProductItf) {
	if result == nil {
		return
	}
	sn.sm.resultChan <- result
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

func (sm *StateMachine) cleanNodes() {
	sm.nowStateNode.GetProduct()
	if sm.nowStateNode.Result.ProductType() == -1 {
		return
	}

	sm.resultChan <- sm.nowStateNode.Result
}

//End insert a end product.
func (sm *StateMachine) End() {
	sm.cleanNodes()
	sm.resultChan <- &ProductEnd{}
}

//Err insert a err product.
func (sm *StateMachine) Err(err string) {
	sm.cleanNodes()
	sm.resultChan <- &ProductError{err}
}

//ErrEnd insert a error product.
func (sm *StateMachine) ErrEnd(err string) {
	sm.cleanNodes()
	sm.resultChan <- &ProductError{err}
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
	sm        *StateMachine
	SNodeList []*StateNode
	LiveMap   map[*StateNode]byte
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
func (dsn *DftStateNodeReader) Register(node StateNodeReader) *DftStateNodeReader {
	dsn.SNodeList = append(dsn.SNodeList, (&StateNode{}).Init(dsn.sm, node))
	return dsn
}

//Clean .
func (dsn *DftStateNodeReader) Clean() {
	for _, node := range dsn.SNodeList {
		node.CleanReader()
		dsn.LiveMap[node] = 0
	}
}

//PreRead DftStateNodeReader don't do PreRead, because it should deal all registed reader.
func (dsn *DftStateNodeReader) PreRead(stateNode *StateNode, input InputItf) (isEnd bool, err error) {
	return false, nil
}

func (dsn *DftStateNodeReader) Read(stateNode *StateNode, input InputItf) (isEnd bool, err error) {
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

/*StateNodeListReader .
 */
type StateNodeListReader struct {
	list   []StateNodeReader
	ptr    int
	result ProductItf
}

//NewStateNodeListReader .
func NewStateNodeListReader(readers ...StateNodeReader) StateNodeReader {
	return &StateNodeListReader{list: readers, ptr: 0}
}

//Name reader's name.
func (s *StateNodeListReader) Name() string {
	return "StateNodeListReader"
}

//Current get current StateNodeReader .
func (s *StateNodeListReader) Current() StateNodeReader {
	return s.list[s.ptr]
}

//Ptr get pointer num.
func (s *StateNodeListReader) Ptr() int {
	return s.ptr
}

//Size get size.
func (s *StateNodeListReader) Size() int {
	return len(s.list)
}

//PreRead only see if should stop read.
func (s *StateNodeListReader) PreRead(stateNode *StateNode, input InputItf) (isEnd bool, err error) {
	current := s.Current()
	lend, lerr := current.PreRead(stateNode, input)
	if lerr != nil {
		return true, lerr
	}
	if lend {
		s.ptr++
		if s.ptr == s.Size() {
			stateNode.Result = current.GetProduct()
			s.result = stateNode.Result
			return true, nil
		}
	}
	return false, nil
}

//Read real read. even isEnd == true the input be readed.
func (s *StateNodeListReader) Read(stateNode *StateNode, input InputItf) (isEnd bool, err error) {
	current := s.Current()
	lend, lerr := current.Read(stateNode, input)
	if lerr != nil {
		return true, lerr
	}
	if lend {
		s.ptr++
		stateNode.Result = current.GetProduct()
		if s.ptr == s.Size() {
			s.result = stateNode.Result
			return true, nil
		}
	}
	return false, nil
}

//GetProduct return result.
func (s *StateNodeListReader) GetProduct() ProductItf {
	return s.result
}

//Clean let the Reader like new.  it will be call before first Read.
func (s *StateNodeListReader) Clean() {
	for _, reader := range s.list {
		reader.Clean()
	}
}

//StateNodeSelectReader .
type StateNodeSelectReader struct {
	ReaderList []StateNodeReader
	LiveMap    map[StateNodeReader]bool
	Result     ProductItf
}

//NewStateNodeSelectReader .
func NewStateNodeSelectReader(list ...StateNodeReader) *StateNodeSelectReader {
	return &StateNodeSelectReader{
		ReaderList: list,
		LiveMap:    make(map[StateNodeReader]bool, len(list)),
	}
}

//Name reader's name.
func (s *StateNodeSelectReader) Name() string {
	return "StateNodeSelectReader"
}

//PreRead only see if should stop read.
func (s *StateNodeSelectReader) PreRead(stateNode *StateNode, input InputItf) (isEnd bool, err error) {
	errList := []string{}
	for cr := range s.LiveMap {
		lend, lerr := cr.PreRead(stateNode, input)
		if lerr != nil {
			errList = append(errList, lerr.Error())
			delete(s.LiveMap, cr)
			continue
		}

		if lend {
			s.Result = cr.GetProduct()
			return true, nil
		}
	}
	if len(s.LiveMap) == 0 {
		return true, fmt.Errorf("Error in StateNodeSelectReader, error list : \n%s", strings.Join(errList, "\n"))
	}
	return false, nil
}

//Read real read. even isEnd == true the input be readed.
func (s *StateNodeSelectReader) Read(stateNode *StateNode, input InputItf) (isEnd bool, err error) {
	errList := []string{}
	for cr := range s.LiveMap {
		lend, lerr := cr.Read(stateNode, input)
		if lerr != nil {
			errList = append(errList, lerr.Error())
			delete(s.LiveMap, cr)
			continue
		}

		if lend {
			s.Result = cr.GetProduct()
			return true, nil
		}
	}
	if len(s.LiveMap) == 0 {
		return true, fmt.Errorf("Error in StateNodeSelectReader, error list : \n%s", strings.Join(errList, "\n"))
	}
	return false, nil
}

//GetProduct return result.
func (s *StateNodeSelectReader) GetProduct() ProductItf {
	return s.Result
}

//Clean let the Reader like new.  it will be call before first Read.
func (s *StateNodeSelectReader) Clean() {
	for _, reader := range s.ReaderList {
		reader.Clean()
		s.LiveMap[reader] = true
	}

	s.Result = nil
}
