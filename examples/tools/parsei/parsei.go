package main

import (
	"encoding/hex"
	"fmt"
	"gotdx/proto"
	"os"
)

func main() {
	for i := 1; i < len(os.Args); i += 1 {
		c := 0
		bytes, err := hex.DecodeString(os.Args[i])
		if err != nil {
			fmt.Printf("%s, decode string err:%s\n", os.Args[i], err)
			continue
		}
		fmt.Printf("%s, %d\n", os.Args[i], proto.ParseInt(bytes, &c))
	}
}
