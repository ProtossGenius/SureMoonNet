package stream

type ConditionFunc func(inp []byte) ([]byte, error)

type ReadPipelineItf interface {
	RemainingSize() int
	Read(size int) ([]byte, error)
	ConditionRead(condition ConditionFunc) ([]byte, error)
	ByteBreakRead(condition byte) ([]byte, error)
}

type ReadContainerItf interface {
	RemainingSize() int
	GetReadPipeline() (ReadPipelineItf,error)
}

type BytePipeline struct {

}
