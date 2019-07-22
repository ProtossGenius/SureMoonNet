package smn_data

const (
	ERR_UNKNOW_TYPE = "err: unknow type. in smn_data_file"
)

func iserr(err error) bool {
	return err != nil
}

func noerr(err error) bool {
	return err == nil
}
