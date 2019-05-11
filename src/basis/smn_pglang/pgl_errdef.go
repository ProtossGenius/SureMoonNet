package smn_pglang

func iserr(err error) bool {
	return err != nil
}

func noerr(err error) bool {
	return err == nil
}
