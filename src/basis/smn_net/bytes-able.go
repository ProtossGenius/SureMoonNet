package smn_net

type BytesAbleItf interface {
	ToBytes() []byte
	FromBytes(bytes []byte) BytesAbleItf
}
