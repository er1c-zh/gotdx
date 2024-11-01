package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"gotdx/proto"
	"os"
)

func main() {
	for i := 1; i < len(os.Args); i += 1 {
		hex, err := hex.DecodeString(os.Args[i])
		if err != nil {
			fmt.Printf("%s, decode err:%s\n", os.Args[i], err)
			continue
		}
		v := uint32(0)
		binary.Read(bytes.NewBuffer(hex), binary.LittleEndian, &v)
		fmt.Printf("%s: %f", os.Args[i], proto.ParseFloat(int32(v)))
	}
}
