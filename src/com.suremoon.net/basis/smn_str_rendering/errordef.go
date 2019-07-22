package smn_str_rendering

const (
	ERR_INIT_NODATA    = "ERR_INIT_NODATA: %s"
	ERR_FILE_NOT_FOUND = "ERR_FILE_NOT_FOUND: %s"
)

func iserr(err error) bool {
	return err != nil
}

func noerr(err error) bool {
	return err == nil
}
