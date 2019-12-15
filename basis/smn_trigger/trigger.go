package smn_trigger

import "sync"

//the signal what's the trigger accept
type TriSignal interface{}

type TriOnSignal interface {
	/*
		Error is passed to the error trigger for processing, no value will be returned here, because the trigger does not
	promise to process any return results
	 */
	OnSignal(sig TriSignal)
}

type Trigger struct {
	regList []TriOnSignal
	lock    sync.Mutex
}

func NewTrigger() *Trigger {
	return &Trigger{regList:make([]TriOnSignal, 0, 15)}
}

func (this *Trigger)OnRegist(onSignal TriOnSignal) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.regList = append(this.regList, onSignal)
}

func (this *Trigger)OnDeal(signal TriSignal)  {
	this.lock.Lock()
	defer this.lock.Unlock()
	for _, dlr := range this.regList{
		dlr.OnSignal(signal)
	}
}
