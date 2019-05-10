package smn_err

func iserr(err error) bool {
	return err != nil
}

func noerr(err error) bool {
	return err == nil
}
