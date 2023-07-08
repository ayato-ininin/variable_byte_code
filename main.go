package main

import (
	"variableByteCode/variable_byte_decode"
	"variableByteCode/variable_byte_encode"
)

func main() {
	vByteEncode.EncodeCsv()
	// vByteDecode.DecodeCsv()
	vByteEncode.Check_encodeValue(uint64(1000))
	vByteDecode.Check_decodeValue([]byte{0xe8, 0x07})
}
