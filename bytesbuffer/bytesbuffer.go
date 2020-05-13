package bytesbuffer

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
)

const (
	BigEndian = iota
	LittleEndian
)

type Buffer struct {
	packetBuffer bytes.Buffer
	enc          binary.ByteOrder
}

func NewBuffer(endian int) (Buffer, error) {
	var b Buffer
	if endian == BigEndian {
		b.enc = binary.BigEndian
	} else if endian == LittleEndian {
		b.enc = binary.LittleEndian
	} else {
		return b, fmt.Errorf("invalid endianness, must be big or little")
	}
	return b, nil
}

func (obj *Buffer) Grow(n int) {
	obj.packetBuffer.Grow(n)
}

func (obj *Buffer) Wrap(data []byte) {
	obj.packetBuffer.Write(data)
}

func (obj *Buffer) PutUint16(value uint16) {
	var buff = make([]byte, 2)
	obj.enc.PutUint16(buff, value)
	obj.packetBuffer.Write(buff)
}

func (obj *Buffer) PutUint32(value uint32) {
	var buff = make([]byte, 4)
	obj.enc.PutUint32(buff, value)
	obj.packetBuffer.Write(buff)
}

func (obj *Buffer) PutUint64(value uint64) {
	var buff = make([]byte, 8)
	obj.enc.PutUint64(buff, value)
	obj.packetBuffer.Write(buff)
}

func (obj *Buffer) GetShort() []byte {
	var tempBuff = obj.packetBuffer.Bytes()
	var shortValue = tempBuff[:2]
	var restValue = tempBuff[2:]
	var byteBuffer bytes.Buffer

	byteBuffer.Write(restValue)
	obj.packetBuffer = byteBuffer
	return shortValue
}

func (obj *Buffer) PutFloat32(value float32) {
	var bits = math.Float32bits(value)
	var buff = make([]byte, 4)
	obj.enc.PutUint32(buff, bits)
	obj.packetBuffer.Write(buff)
}

func (obj *Buffer) PutFloat64(value float64) {
	var bits = math.Float64bits(value)
	var buff = make([]byte, 8)
	obj.enc.PutUint64(buff, bits)
	obj.packetBuffer.Write(buff)
}

func (obj *Buffer) GetFloat() []byte {
	var tempBuff = obj.packetBuffer.Bytes()
	var floatValue = tempBuff[:4]
	var restValue = tempBuff[4:]
	var byteBuffer bytes.Buffer
	byteBuffer.Write(restValue)
	obj.packetBuffer = byteBuffer
	return floatValue
}

func (obj *Buffer) GetDouble() []byte {
	var tempBuff = obj.packetBuffer.Bytes()
	var doubleValue = tempBuff[:8]
	var restValue = tempBuff[8:]
	var byteBuffer bytes.Buffer
	byteBuffer.Write(restValue)
	obj.packetBuffer = byteBuffer
	return doubleValue
}

func (obj *Buffer) Put(value []byte) {
	obj.packetBuffer.Write(value)
}

func (obj *Buffer) PutByte(value byte) {
	var tempByte = []byte{value}
	obj.packetBuffer.Write(tempByte)
}

func (obj *Buffer) Get(size int) []byte {
	var tempBuff = obj.packetBuffer.Bytes()
	var value = tempBuff[:size]
	var restValue = tempBuff[size:]
	var byteBuffer bytes.Buffer
	byteBuffer.Write(restValue)
	obj.packetBuffer = byteBuffer
	return value
}

func (obj *Buffer) GetByte() []byte {
	var tempBuff = obj.packetBuffer.Bytes()
	var value = tempBuff[:1]
	var restValue = tempBuff[1:]
	var byteBuffer bytes.Buffer
	byteBuffer.Write(restValue)
	obj.packetBuffer = byteBuffer
	return value
}

func (obj *Buffer) Buffer() bytes.Buffer {
	return obj.packetBuffer
}

func (obj *Buffer) Bytes() []byte {
	return obj.packetBuffer.Bytes()
}

func (obj *Buffer) Size() int {
	return len(obj.packetBuffer.Bytes())
}

func (obj *Buffer) Flip() {
	var bytesArr = obj.packetBuffer.Bytes()
	for i, j := 0, len(bytesArr)-1; i < j; i, j = i+1, j-1 {
		bytesArr[i], bytesArr[j] = bytesArr[j], bytesArr[i]
	}
	var byteBuffer bytes.Buffer
	byteBuffer.Write(bytesArr)
	obj.packetBuffer = byteBuffer
}

func (obj *Buffer) Clear() {
	obj.packetBuffer = bytes.Buffer{}
}

func (obj *Buffer) Slice(start int, end int) error {
	var bytesArr = obj.packetBuffer.Bytes()
	if len(bytesArr) < (start + end) {
		return fmt.Errorf("buffer too small")
	}
	bytesArr = bytesArr[start:end]
	var byteBuffer bytes.Buffer
	byteBuffer.Write(bytesArr)
	obj.packetBuffer = byteBuffer
	return nil
}

func (obj *Buffer) Bytes2Str(data []byte) string {
	return string(data)
}

func (obj *Buffer) Str2Bytes(data string) []byte {
	return []byte(data)
}

func (obj *Buffer) Bytes2Short(data []byte) uint16 {
	return obj.enc.Uint16(data)
}

func (obj *Buffer) Bytes2Int(data []byte) uint32 {
	return obj.enc.Uint32(data)
}

func (obj *Buffer) Bytes2Long(data []byte) uint64 {
	return obj.enc.Uint64(data)
}
