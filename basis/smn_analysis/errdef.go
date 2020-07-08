package smn_analysis

func iserr(err error) bool {
	return err != nil
}

func noerr(err error) bool {
	return err == nil
}

const (
	ErrNoMatchStateNode     = "ErrNoMatchStateNode: err list: [%s]"                               //没有满足的
	ErrTooMuchStateNodeLive = "ErrTooMuchStateNodeLive, input[%v] live nodes [%s], end nodes[%s]" //太多满足条件的
)
