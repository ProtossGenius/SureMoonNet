package smn_godll_head

func iserr(err error) bool {
	return err != nil
}

func noerr(err error) bool {
	return err == nil
}

const (
	ErrUnexpectedCppFunctionDefinition = "ErrUnexpectedCppFunctionDefinition: %s"
)
