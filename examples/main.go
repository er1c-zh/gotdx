package main

import (
	"fmt"
	"log"
	"os"

	"gotdx/tdx"
)

func main() {
	f, err := os.Create("output_sh.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	// ip地址如果失效，请自行替换
	cli := tdx.New(tdx.WithTCPAddress("124.71.187.122:7709"))
	_, err = cli.Connect()
	if err != nil {
		log.Fatalln(err)
	}
	defer cli.Disconnect()

	cursor := uint16(0)
	count := 0
	for {
		reply, err := cli.GetSecurityList(tdx.MarketSh, cursor)
		if err != nil {
			log.Println(err)
			return
		}
		count += int(reply.Count)

		for _, obj := range reply.List {
			fmt.Fprintf(f, "%s %032b %032b %s\n", obj.Code, obj.Reserved1, obj.Reserved2, obj.Name)
		}

		if len(reply.List) < 1000 {
			break
		}
		cursor += 1000
	}

	log.Printf("%d / 1000\n", count)
}
