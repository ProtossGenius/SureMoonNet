package smn_net

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

func ReadInt64(reader io.Reader) (int64, error) {
	var res int64
	err := binary.Read(reader, binary.BigEndian, &res)
	if iserr(err) {
		return 0, err
	}
	return res, nil
}

func ReadInt32(reader io.Reader) (int32, error) {
	var res int32
	err := binary.Read(reader, binary.BigEndian, &res)
	if iserr(err) {
		return 0, err
	}
	return res, nil
}

func WriteInt64(val int64, writer io.Writer) error {
	return binary.Write(writer, binary.BigEndian, val)
}

func WriteInt32(val int32, writer io.Writer) error {
	return binary.Write(writer, binary.BigEndian, val)
}

func ReadUint64(reader io.Reader) (uint64, error) {
	var res uint64
	err := binary.Read(reader, binary.BigEndian, &res)
	if iserr(err) {
		return 0, err
	}
	return res, nil
}

func ReadUint32(reader io.Reader) (uint32, error) {
	var res uint32
	err := binary.Read(reader, binary.BigEndian, &res)
	if iserr(err) {
		return 0, err
	}
	return res, nil
}

func WriteUint64(val uint64, writer io.Writer) error {
	return binary.Write(writer, binary.BigEndian, val)
}

func WriteUint32(val uint32, writer io.Writer) error {
	return binary.Write(writer, binary.BigEndian, val)
}

func ReadInt16(reader io.Reader) (int16, error) {
	var res int16
	err := binary.Read(reader, binary.BigEndian, &res)
	if iserr(err) {
		return 0, err
	}
	return res, nil
}

func ReadInt8(reader io.Reader) (int8, error) {
	var res int8
	err := binary.Read(reader, binary.BigEndian, &res)
	if iserr(err) {
		return 0, err
	}
	return res, nil
}

func WriteInt16(val int16, writer io.Writer) error {
	return binary.Write(writer, binary.BigEndian, val)
}

func WriteInt8(val int8, writer io.Writer) error {
	return binary.Write(writer, binary.BigEndian, val)
}

func ReadUint16(reader io.Reader) (uint16, error) {
	var res uint16
	err := binary.Read(reader, binary.BigEndian, &res)
	if iserr(err) {
		return 0, err
	}
	return res, nil
}

func ReadUint8(reader io.Reader) (uint8, error) {
	var res uint8
	err := binary.Read(reader, binary.BigEndian, &res)
	if iserr(err) {
		return 0, err
	}
	return res, nil
}

func WriteUint16(val uint16, writer io.Writer) error {
	return binary.Write(writer, binary.BigEndian, val)
}

func WriteUint8(val uint8, writer io.Writer) error {
	return binary.Write(writer, binary.BigEndian, val)
}

func ReadUint(reader io.Reader) (uint, error) {
	res, err := ReadUint64(reader)
	return uint(res), err
}

func WriteUint(val uint, writer io.Writer) error {
	return WriteUint64(uint64(val), writer)
}

func ReadInt(reader io.Reader) (int, error) {
	res, err := ReadInt64(reader)
	return int(res), err
}

func WriteInt(val int, writer io.Writer) error {
	return WriteInt64(int64(val), writer)
}

func ReadFloat64(reader io.Reader) (float64, error) {
	var res float64
	err := binary.Read(reader, binary.BigEndian, &res)
	if iserr(err) {
		return 0, err
	}
	return res, nil
}

func ReadFloat32(reader io.Reader) (float32, error) {
	var res float32
	err := binary.Read(reader, binary.BigEndian, &res)
	if iserr(err) {
		return 0, err
	}
	return res, nil
}

func WriteFloat64(val float64, writer io.Writer) error {
	return binary.Write(writer, binary.BigEndian, val)
}

func WriteFloat32(val float32, writer io.Writer) error {
	return binary.Write(writer, binary.BigEndian, val)
}

func ReadBytes(reader io.Reader) ([]byte, error) {
	len, err := ReadInt32(reader)
	if iserr(err) {
		return nil, err
	}
	bts := make([]byte, len)
	n, err := reader.Read(bts)
	if n < int(len) {
		return nil, fmt.Errorf(ErrNotGetEnoughLengthBytes, len, n)
	}
	if iserr(err) {
		return nil, err
	}
	return bts, nil
}

func ReadString(reader io.Reader) (string, error) {
	bts, err := ReadBytes(reader)
	return string(bts), err
}

func WriteBytes(bts []byte, writer io.Writer) (int, error) {
	blns := len(bts)
	buffer := bytes.NewBuffer(make([]byte, 0, blns+4))
	err := WriteInt32(int32(blns), buffer)
	if iserr(err) {
		return 0, err
	}
	buffer.Write(bts)

	return writer.Write(buffer.Bytes())
}

func WriteString(val string, writer io.Writer) (int, error) {
	return WriteBytes([]byte(val), writer)
}
