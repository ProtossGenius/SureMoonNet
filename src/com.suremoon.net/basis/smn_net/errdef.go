package smn_net

func iserr(err error) bool {
	return err != nil
}

func noerr(err error) bool {
	return err == nil
}

const (
	ErrNotGetEnoughLengthBytes = "ErrNotGetEnoughLengthBytes: needs:[%d], get[%d]"
)
