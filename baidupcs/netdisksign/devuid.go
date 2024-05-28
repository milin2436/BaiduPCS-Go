package netdisksign

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"github.com/milin2436/BaiduPCS-Go/pcsutil/cachepool"
	"github.com/milin2436/BaiduPCS-Go/pcsutil/converter"
)

func DevUID(feature string) string {
	m := md5.New()
	m.Write(converter.ToBytes(feature))
	res := m.Sum(nil)
	resHex := cachepool.RawMallocByteSlice(34)
	hex.Encode(resHex[:32], res)
	resHex[32] = '|'
	resHex[33] = '0'
	return converter.ToString(bytes.ToUpper(resHex))
}
