package smn_anlys_go_tif

func iserr(err error) bool {
	return err != nil
}

func noerr(err error) bool {
	return err == nil
}

const (
	ErrUnexpectedGoFunctionDefinition = "ErrUnexpectedGoFunctionDefinition: %s"
)
