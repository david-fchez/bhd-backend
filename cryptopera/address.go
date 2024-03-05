/*
	Misc functions to generate bech32 address and various
	converters. Some I wrote, some I've copied from
	various github sources
	GM 2020-10-26
 */
package cryptopera

import (
	"encoding/hex"
	"fmt"
	"strings"
)

const (
		charset = "qpzry9x8gf2tvdw0s3jn54khce6mua7l"
)

func bytes2Bits(data []byte) []byte {
	dst := make([]byte, 0)
	for _, v := range data {
		for i := 0; i < 8; i++ {
			move := uint(7 - i)
			dst = append(dst, byte((v>>move)&1))
		}
	}
	return dst
}


func convertTo5BitsFromStr(input string) ([]byte, error) {
	payload, err := hex.DecodeString(input)
	if err != nil {
		return nil,err
	}
	return convertTo5Bits(payload)
}

func toLower5BitFromString(input string) []byte {
	return toLower5Bit([]byte(input))
}

func toLower5Bit(input []byte) []byte {
	outArr := make([]byte, len(input))
	for i,v := range input {
		outArr[i] = v & 0x1F
	}
	return outArr
}


func convertTo5Bits(input []byte)  ([]byte, error) {
	bitArray := bytes2Bits(input)
	padLen := len(bitArray) % 5
	// pad with zeros to the right
	if padLen != 0 {
		for i := 0; i < (5 - padLen); i++ {
			bitArray = append(bitArray, 0)
		}
	}
	outArray := make([]byte, len(bitArray)/5)
	index := 0
	for i:=0; i < len(bitArray);i+=5 {
		value := bitArray[i]*16 + bitArray[i+1]*8 + bitArray[i+2]*4 + bitArray[i+3]*2 + bitArray[i+4]
		outArray[index] = value
		index++
	}
	return outArray, nil
}


func convertTo8Bits(input []byte)  ([]byte, error) {
	bitArray := bytes2Bits(input)
	//no padding in this case
	outArray := make([]byte, len(bitArray)/8)
	index := 0
	for i:=0; i < len(bitArray);i+=8 {
		value := bitArray[i]*128 + bitArray[i+1]*64 + bitArray[i+2]*32 + bitArray[i+3]*1 + bitArray[i+4]*8 + bitArray[i+5]*4 + bitArray[i+6]*2 + bitArray[i+7]
		outArray[index] = value
		index++
	}
	return outArray, nil
}


func calcPolymod(input []byte) int {
	var c int = 1
	for _,v := range input {
		c0 := c >> 35
		c = ((c & 0x07ffffffff) << 5) ^ int(v)
		if c0 & 0x01 != 0 {
			c ^= 0x98f2bc8e61
		}
		if c0 & 0x02 != 0 {
			c ^= 0x79b76d99e2
		}
		if c0 & 0x04 != 0 {
			c ^= 0xf33e5fb3c4
		}

		if c0 & 0x08 != 0 {
			c ^= 0xae2eabe2a8
		}
		if c0 & 0x10 != 0 {
			c ^= 0x1e4f43e470
		}
	}
	return c ^ 1
}


func toChars(data []byte) (string, error) {
	result := make([]byte, 0, len(data))
	for _, b := range data {
		if int(b) >= len(charset) {
			return "", fmt.Errorf("invalid data byte: %v", b)
		}
		result = append(result, charset[b])
	}
	return string(result), nil
}

// toBytes converts each character in the string 'chars' to the value of the
// index of the correspoding character in 'charset'.
func ToBytes(chars string) ([]byte, error) {
	decoded := make([]byte, 0, len(chars))
	for i := 0; i < len(chars); i++ {
		index := strings.IndexByte(charset, chars[i])
		if index < 0 {
			return nil, fmt.Errorf("invalid character not part of "+
				"charset: %v", chars[i])
		}
		decoded = append(decoded, byte(index))
	}
	return decoded, nil
}


func GetBech32Address(prefix string, hash []byte) (string, error){
	var address string = ""
	fiveBitArr, err := convertTo5Bits(append([]byte{0},hash...))
	if err != nil {
		return address,err
	}
	// calculate the checksum
	// now calculate the checksum
	prefixLower5bits := toLower5BitFromString(prefix) // []byte{ 2, 9, 20, 3, 15, 9, 14, 3, 1, 19, 8, 0}
	// to prefix append 0 for separator
	prefixLower5bits = append(prefixLower5bits,0)
	checkSumTempalate := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	checkSum := calcPolymod(append(append(prefixLower5bits,fiveBitArr...),checkSumTempalate...))
	// padd it
	checkSumStr := fmt.Sprintf("%X",checkSum)
	if len(checkSumStr) % 2 == 1 {
		checkSumStr = "0" + checkSumStr
	}
	ckhSumArr,err := hex.DecodeString(checkSumStr)
	if err != nil {
		return address,err
	}
	checkSum5bitArray,err := convertTo5Bits(ckhSumArr)
	if err != nil {
		return address,err
	}
	payloadB32,err := toChars(fiveBitArr)
	if err != nil {
		return address,err
	}
	checksumB32,err := toChars(checkSum5bitArray)
	if err != nil {
		return address,err
	}
	address = prefix + ":" + payloadB32 + checksumB32
	return address,nil
}