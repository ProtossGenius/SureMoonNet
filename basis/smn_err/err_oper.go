package smn_err

func iserr(err error) bool {
	return err != nil
}

func noerr(err error) bool {
	return err == nil
}

type ErrDeal struct {
	onErr OnErr
}

func (this *ErrDeal) OnErr(err error) {
	this.onErr(err)
}

func NewErrDeal() *ErrDeal {
	return &ErrDeal{onErr: DftOnErr}
}

type OnErr func(err error)

func DftOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
