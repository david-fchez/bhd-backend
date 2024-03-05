package utils

import (
	"bytes"
	"encoding/binary"
	"math"
	"net"
)

// SatoshiToMonetary converts from satoshi units to friendly
// view in bch decimal amount
func SatoshiToMonetary(satoshi int64) float64 {
	var res = float64(satoshi) / 100000000.0
	res = math.Round(res*10000) / 10000.0
	return res
}

// SatoshiToMonetaryPerUnit returns price per unit, so amount divided
// by the quantity
func SatoshiToMonetaryPerUnit(satoshi int64, quantity int64) float64 {
	var sat = int64(math.Round(float64(satoshi) / float64(quantity)))
	return SatoshiToMonetary(sat)
}

// SatoshiFromMonetary returns satoshi's from monetary format
func SatoshiFromMonetary(monetaryAmount float64) int64 {
	return int64(math.Round(monetaryAmount * 10000000.0))
}

// Int16ToByte converts 2 byte integer to litle endian
// byte representation
func Int16ToByte(val int16) []byte {
	var tempBuff = []byte{0x00, 0x00}
	binary.LittleEndian.PutUint16(tempBuff, uint16(val))
	return tempBuff
}

func UInt16ToByteBe(val uint16) []byte {
	var tempBuff = []byte{0x00, 0x00}
	binary.BigEndian.PutUint16(tempBuff, val)
	return tempBuff
}

func UInt16ToByte(val uint16) []byte {
	var tempBuff = []byte{0x00, 0x00}
	binary.LittleEndian.PutUint16(tempBuff, val)
	return tempBuff
}

// Int32ToByte converts 4 byte integer to litle endian
// byte representation
func Int32ToByte(val int32) []byte {
	var tempBuff = []byte{0x00, 0x00, 0x00, 0x00}
	binary.LittleEndian.PutUint32(tempBuff, uint32(val))
	return tempBuff
}

// UInt32ToByte converts 4 byte integer to litle endian
// byte representation
func UInt32ToByte(val uint32) []byte {
	var tempBuff = []byte{0x00, 0x00, 0x00, 0x00}
	binary.LittleEndian.PutUint32(tempBuff, val)
	return tempBuff
}

// Int64ToByte converts 8 byte integer to litle endian
// byte representation
func Int64ToByte(val int64) []byte {
	var tempBuff = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	binary.LittleEndian.PutUint64(tempBuff, uint64(val))
	return tempBuff
}

// UInt64ToByte converts 8 byte integer to litle endian
// byte representation
func UInt64ToByte(val uint64) []byte {
	var tempBuff = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	binary.LittleEndian.PutUint64(tempBuff, val)
	return tempBuff
}

// IntToByte converts 8 byte integer (golang default) to litle endian
// byte representation
func IntToByte(val int) []byte {
	return Int64ToByte(int64(val))
}

// ByteToInt16 returns int (2 byte) from little endian byte representation
func ByteToInt16(val []byte) int16 {
	return int16(binary.LittleEndian.Uint16(val))
}

// ByteToUInt16 returns int (2 byte) from little endian byte representation
func ByteToUInt16(val []byte) uint16 {
	return binary.LittleEndian.Uint16(val)
}

// ByteToInt32 returns int (4 byte) from little endian byte representation
func ByteToInt32(val []byte) int32 {
	return int32(binary.LittleEndian.Uint32(val))
}

// ByteToUInt32 returns int (4 byte) from little endian byte representation
func ByteToUInt32(val []byte) uint32 {
	return uint32(binary.LittleEndian.Uint32(val))
}

// ByteToUInt64 returns int (8 byte) from little endian byte representation
func ByteToUInt64(val []byte) uint64 {
	return binary.LittleEndian.Uint64(val)
}

// ByteToInt returns int (8 byte) from little endian byte representation
func ByteToInt(val []byte) int {
	return int(ByteToUInt64(val))
}

// ReadUint8 reads uint8 from bytes reader
func ReadUint8(reader *bytes.Reader) uint8 {
	val, err := reader.ReadByte()
	if err != nil {
		panic(err)
	}
	return uint8(val)
}

// ReadUint16 reads uint16 from bytes reader
func ReadUint16(reader *bytes.Reader) uint16 {
	var tempBuff = []byte{0x00, 0x00}
	_, err := reader.Read(tempBuff)
	if err != nil {
		panic(err)
	}
	return ByteToUInt16(tempBuff)
}

func ReadIp(reader *bytes.Reader) string {
	// if stats with 0x00000000000000000000FFFF then its ipv4
	var ip = make([]byte, 16)
	_, err := reader.Read(ip)
	if err != nil {
		panic(err)
	}
	return net.IP(ip).String()
}

// ReadUint16Be reads uint16 from bytes reader in big endian order
func ReadUint16Be(reader *bytes.Reader) uint16 {
	var tempBuff = []byte{0x00, 0x00}
	_, err := reader.Read(tempBuff)
	if err != nil {
		panic(err)
	}
	return binary.BigEndian.Uint16(tempBuff)
}

// ReadUint32 reads uint32 from bytes reader
func ReadUint32(reader *bytes.Reader) uint32 {
	var tempBuff = []byte{0x00, 0x00, 0x00, 0x00}
	_, err := reader.Read(tempBuff)
	if err != nil {
		panic(err)
	}
	return ByteToUInt32(tempBuff)
}

// ReadUint64 reads uint64 from bytes reader
func ReadUint64(reader *bytes.Reader) uint64 {
	var tempBuff = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	_, err := reader.Read(tempBuff)
	if err != nil {
		panic(err)
	}
	return ByteToUInt64(tempBuff)
}

func VarIntToByte(val uint64) []byte {
	var buff []byte
	if val < 0xFD {
		buff = []byte{uint8(val)}
	} else if val <= math.MaxUint16 {
		buff = append([]byte{0xFD}, UInt16ToByte(uint16(val))...)
	} else if val <= math.MaxUint32 {
		buff = append([]byte{0xFE}, UInt32ToByte(uint32(val))...)
	} else {
		buff = append([]byte{0xFF}, UInt64ToByte(uint64(val))...)
	}
	return buff
}

// VarStringToByte encodes the string with length at start
func VarStringToByte(str string) []byte {
	val := len(str)
	var buff []byte
	if val < 0xFD {
		buff = []byte{uint8(val)}
	} else if val <= math.MaxUint16 {
		buff = append([]byte{0xFD}, UInt16ToByte(uint16(val))...)
	} else if val <= math.MaxUint32 {
		buff = append([]byte{0xFE}, UInt32ToByte(uint32(val))...)
	} else {
		buff = append([]byte{0xFF}, UInt64ToByte(uint64(val))...)
	}
	return append(buff, []byte(str)...)
}

func ReadVarInt(reader *bytes.Reader) uint64 {
	firstByte, err := reader.ReadByte()
	if err != nil {
		panic(err)
	}
	if firstByte < 0xFD {
		return uint64(firstByte)
	}
	if firstByte == 0xFD {
		return uint64(ReadUint16(reader))
	}
	if firstByte == 0xFE {
		return uint64(ReadUint32(reader))
	}
	if firstByte == 0xFF {
		return uint64(ReadUint64(reader))
	}
	panic("why am I here?")
}

func ReadVarString(reader *bytes.Reader) string {
	// read the first byte, that will tell us the length
	var strLength = ReadVarInt(reader)
	var buff = make([]byte, strLength)
	if strLength > 0 {
		_, err := reader.Read(buff)
		if err != nil {
			panic(err)
		}
	}
	// trim buff
	buff = bytes.TrimRight(buff, string(rune(0)))
	return string(buff)
}
