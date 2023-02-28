package internal

import (
	"encoding/hex"
)

func EncodeHex(buf []byte, bytesUUID []byte) {
	hex.Encode(buf, bytesUUID[:4])
	buf[8] = '-'
	hex.Encode(buf[9:13], bytesUUID[4:6])
	buf[13] = '-'
	hex.Encode(buf[14:18], bytesUUID[6:8])
	buf[18] = '-'
	hex.Encode(buf[19:23], bytesUUID[8:10])
	buf[23] = '-'
	hex.Encode(buf[24:], bytesUUID[10:])
}
