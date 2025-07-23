package buffer_test

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"gatesvr/core/buffer"
	"gatesvr/utils/xrand"
	"testing"
)

type User struct {
	ID  int32
	Age int8
}

func TestNewBuffer(t *testing.T) {
	buff := &bytes.Buffer{}
	buff.Grow(2)

	binary.Write(buff, binary.BigEndian, int16(2))

	fmt.Println(buff.Bytes())

	writer := buffer.NewWriter(2)
	writer.WriteInt16s(binary.BigEndian, int16(2))

	fmt.Println(writer.Bytes())

	writer.Reset()
	writer.WriteInt16s(binary.BigEndian, int16(20))
	writer.WriteFloat32s(binary.BigEndian, 5.2)

	fmt.Println(writer.Bytes())

	data := writer.Bytes()

	reader := buffer.NewReader(data)
	v1, _ := reader.ReadInt16(binary.BigEndian)
	fmt.Println(v1)
	v2, _ := reader.ReadFloat32(binary.BigEndian)
	fmt.Println(v2)
}

func BenchmarkBuffer1(b *testing.B) {
	data := []byte(xrand.Letters(1024))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		buff := &bytes.Buffer{}
		buff.Grow(1024)
		binary.Write(buff, binary.BigEndian, data)
		buff.Reset()
	}
}

func BenchmarkBuffer2(b *testing.B) {
	writer := buffer.NewWriter(8)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		writer.WriteInt64s(binary.BigEndian, 2)
		writer.Reset()
	}
}

func BenchmarkNocopyBuffer_Malloc(b *testing.B) {
	data := []byte(xrand.Letters(1024))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		buf := buffer.NewNocopyBuffer()
		buf.Mount(data)
		buf.Release()
	}
}

func TestNewBuffer2(t *testing.T) {
	buff := buffer.NewNocopyBuffer()

	writer1 := buff.Malloc(8)
	writer1.WriteInt64s(binary.BigEndian, 2)

	writer2 := buff.Malloc(8)
	writer2.WriteInt64s(binary.BigEndian, 3)

	t.Log(buff.Len())
	t.Log(buff.Len())

	buff.Range(func(node *buffer.NocopyNode) bool {
		t.Log(node.Bytes())
		return true
	})

	buff.Release()

	fmt.Println(buff.Bytes())

}

func TestNocopyBuffer_Malloc(t *testing.T) {
	buff := buffer.NewNocopyBuffer()

	buff.Malloc(10)

	buff.Malloc(250)
}

func TestNocopyBuffer_Mount(t *testing.T) {
	buff1 := buffer.NewNocopyBuffer()

	writer1 := buff1.Malloc(8)
	writer1.WriteInt64s(binary.BigEndian, 1)

	writer2 := buff1.Malloc(8)
	writer2.WriteInt64s(binary.BigEndian, 2)

	buff2 := buffer.NewNocopyBuffer()

	writer3 := buff2.Malloc(8)
	writer3.WriteInt64s(binary.BigEndian, 3)

	writer4 := buff2.Malloc(8)
	writer4.WriteInt64s(binary.BigEndian, 4)

	buff1.Mount(buff2, buffer.Head)

	fmt.Println(buff1.Bytes())
}
