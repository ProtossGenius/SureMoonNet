package lw_analysis

func iserr(err error) bool {
	return err != nil
}

func noerr(err error) bool {
	return err == nil
}

const (
	ErrWaitVarName = "ErrWaitVarName"
)
