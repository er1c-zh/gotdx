package main

import (
	"encoding/hex"
	"fmt"
	"gotdx/tdx"
	"os"
	"time"
)

func main() {
	ts := time.Now().Unix()
	for i := 1; i < len(os.Args); i += 1 {
		bytes, err := hex.DecodeString(os.Args[i])
		if err != nil {
			fmt.Printf("%s, decode string err:%s\n", os.Args[i], err)
			continue
		}

		header, err := tdx.ParseReqHeader(bytes[:12])
		if err != nil {
			fmt.Printf("%s, parse header err:%s\n", os.Args[i], err)
			continue
		}

		{
			fmt.Printf("header: %04X\n", header.Method)
		}

		// save to file
		f, err := os.OpenFile(fmt.Sprintf("%d_%04X_%d.bin", ts, header.Method, i), os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer f.Close()
		_, err = f.Write(bytes)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
}
