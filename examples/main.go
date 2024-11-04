package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"gotdx/models"
	ee "gotdx/proto/v2"
	"gotdx/tdx"
)

func main() {
	var err error
	cli := ee.NewClient(tdx.DefaultOption.
		WithTCPAddress("110.41.147.114:7709").WithDebugMode().WithMsgCallback(func(pi models.ProcessInfo) {
		fmt.Printf("%s\n", pi.Msg)
	}))
	err = cli.Connect()
	if err != nil {
		fmt.Printf("error:%s", err)
		return
	}
	fmt.Printf("connected\n")
	time.Sleep(60 * time.Second)
}

func main1() {
	var err error

	// ip地址如果失效，请自行替换
	// cli := tdx.New(tdx.WithTCPAddress("124.71.187.122:7709"))
	cli := tdx.New(tdx.DefaultOption.
		WithTCPAddress("110.41.147.114:7709").WithDebugMode())
	_, err = cli.Connect()
	if err != nil {
		log.Fatalln(err)
	}
	defer cli.Disconnect()

	reply, err := cli.GetSecurityQuotes([]tdx.StockQuery{
		{Market: tdx.MarketSz, Code: "000100"},
		{Market: tdx.MarketSh, Code: "600000"},
		{Market: tdx.MarketSz, Code: "001979"},
		{Market: tdx.MarketSh, Code: "600048"},
		{Market: tdx.MarketSz, Code: "300748"},
	})
	if err != nil {
		log.Println(err)
		return
	}
	j, _ := json.MarshalIndent(reply, "", "  ")
	log.Println(string(j))
}
