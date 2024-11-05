package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"gotdx/models"
	ee "gotdx/proto/v2"
	"gotdx/tdx"

	"bytes"
)

func main() {
	var err error
	cli := ee.NewClient(tdx.DefaultOption.
		WithTCPAddress("110.41.147.114:7709").
		WithDebugMode().
		WithMsgCallback(func(pi models.ProcessInfo) {
			fmt.Printf("%s\n", pi.Msg)
		}).WithMetaAddress("124.71.223.19:7727"),
	)
	conn, err := cli.NewMetaConnection()
	if err != nil {
		fmt.Printf("error:%s", err)
		return
	}
	err = cli.MetaShakehand(conn)
	if err != nil {
		fmt.Printf("error:%s", err)
		return
	}
	fmt.Printf("connected\n")

	offset := uint32(0)

	f, err := os.OpenFile("test_data/descmap.txt", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	for {
		descMapResp, err := cli.MetaDescMap(conn, offset)
		if err != nil {
			fmt.Printf("error:%s", err)
			return
		}

		buf := bytes.NewBuffer(nil)
		offset += uint32(descMapResp.Count)
		for _, d := range descMapResp.List {
			fmt.Printf("%s %s\n", d.IDInUtf8, d.DescInUtf8)
			buf.WriteString(fmt.Sprintf("%s %s %s\n", hex.EncodeToString(d.Reserved0), d.IDInUtf8, d.DescInUtf8))
		}
		f.Write(buf.Bytes())
		if descMapResp.Count < 500 {
			break
		}
	}

}

func main2() {
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

	if false {
		data, err := cli.List([]ee.StockQuery{
			{Market: tdx.MarketSh, Code: "603230"},
			{Market: tdx.MarketSh, Code: "601216"},
			{Market: tdx.MarketSz, Code: "000100"},
			{Market: tdx.MarketSz, Code: "300059"},
		})
		if err != nil {
			fmt.Printf("error:%s", err)
			return
		}

		// j, _ := json.MarshalIndent(data, "", "  ")
		// fmt.Printf("%s\n", j)
		for _, obj := range data.List {
			buf := bytes.NewBuffer(nil)
			for _, b := range obj.Reserved3 {
				buf.WriteString(fmt.Sprintf("%08b", b))
			}
			fmt.Printf("%s buf:%s\n", obj.Code, buf.String())
		}
	}

	if false {
		data, err := cli.Rank("region-desc-all")
		if err != nil {
			fmt.Printf("error:%s", err)
			return
		}

		for _, obj := range data.List {
			buf := bytes.NewBuffer(nil)
			for _, b := range obj.Reserved3 {
				buf.WriteString(fmt.Sprintf("%08b", b))
			}
			fmt.Printf("%s buf:%s\n", obj.Code, buf.String())
		}
	}

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
