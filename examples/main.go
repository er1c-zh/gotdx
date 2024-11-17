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
	// testStockMeta()
	// testServerInfo()
	// testDownloadFile()
	test0547()
	// testServerInfo()
}

func test0547() {
	var err error
	cli := ee.NewClient(tdx.DefaultOption.
		WithDebugMode().
		WithTCPAddress("110.41.147.114:7709").
		WithDebugMode().
		WithMsgCallback(func(pi models.ProcessInfo) {
			fmt.Printf("%s\n", pi.Msg)
		}).WithMetaAddress("124.71.223.19:7727"))
	err = cli.Connect()
	if err != nil {
		fmt.Printf("error:%s", err)
		return
	}
	fmt.Printf("connected\n")

	// cli.TDXHandshake()
	// cli.Heartbeat()

	resp, err := cli.Realtime([]ee.StockQuery{
		{Market: tdx.MarketSh, Code: "999999"},
		{Market: tdx.MarketSz, Code: "399002"},
		{Market: tdx.MarketSz, Code: "300059"},
		{Market: tdx.MarketSz, Code: "300010"},
		{Market: tdx.MarketSh, Code: "999998"},
		{Market: tdx.MarketSh, Code: "999997"},
	})
	if err != nil {
		fmt.Printf("error:%s", err)
		return
	}
	fmt.Printf("%s\n", hex.Dump(resp.Data))
	j, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Printf("%s\n", j)
}

func testStockMeta() {
	var err error
	cli := ee.NewClient(tdx.DefaultOption.
		WithDebugMode().
		WithTCPAddress("110.41.147.114:7709").
		WithDebugMode().
		WithMsgCallback(func(pi models.ProcessInfo) {
			fmt.Printf("%s\n", pi.Msg)
		}).WithMetaAddress("124.71.223.19:7727"))
	err = cli.Connect()
	if err != nil {
		fmt.Printf("error:%s", err)
		return
	}
	fmt.Printf("connected\n")

	resp, err := cli.StockMeta(tdx.MarketSh, 0)
	if err != nil {
		fmt.Printf("error:%s", err)
		return
	}
	j, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Printf("%s\n", j)
}

func testDownloadFile() {
	var err error
	cli := ee.NewClient(tdx.DefaultOption.
		WithDebugMode().
		WithTCPAddress("110.41.147.114:7709").
		WithDebugMode().
		WithMsgCallback(func(pi models.ProcessInfo) {
			fmt.Printf("%s\n", pi.Msg)
		}).WithMetaAddress("124.71.223.19:7727"))
	err = cli.Connect()
	if err != nil {
		fmt.Printf("error:%s", err)
		return
	}
	fmt.Printf("connected\n")

	/*
		"block_gn.dat"
		"block_fg.dat"
		"block_zs.dat"
		"tdxhy.cfg"
		"spec/speckzzdata.txt"
		"spec/specetfdata.txt"
		"spec/speclofdata.txt"
		"spec/specgpext.txt"
		"tdxzsbase.cfg"
		"zhb.zip"
	*/

	for _, fileName := range []string{
		// "infoharbor_block.dat",
		// "infoharbor_ex.code",
		// "infoharbor_ex.name",
		// "block_gn.dat",
		// "block_fg.dat",
		// "block_zs.dat",
		// "tdxhy.cfg",
		// "spec/speckzzdata.txt",
		// "spec/specetfdata.txt",
		// "spec/speclofdata.txt",
		// "spec/specgpext.txt",
		// "spec/speczsevent.txt",
		// "spec/speczshot.txt",
		// "tdxzsbase.cfg",
		"tdxzsbase2.cfg",
		// "zhb.zip",
	} {

		data, err := cli.DownloadFile(fileName)
		if err != nil {
			fmt.Printf("error:%s", err)
			return
		}

		f, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			f.Close()
			continue
		}
		f.Write(data)
		f.Close()
	}

}

func main3() {
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
			fmt.Printf("%s %s\n", d.ID, d.Desc)
			buf.WriteString(fmt.Sprintf("%s %s %s\n", hex.EncodeToString(d.Reserved0), d.ID, d.Desc))
		}
		f.Write(buf.Bytes())
		if descMapResp.Count < 500 {
			break
		}
	}
}

func testServerInfo() {
	var err error
	cli := ee.NewClient(tdx.DefaultOption.
		WithDebugMode().
		WithTCPAddress("110.41.147.114:7709").WithDebugMode().WithMsgCallback(func(pi models.ProcessInfo) {
		fmt.Printf("%s\n", pi.Msg)
	}))
	err = cli.Connect()
	if err != nil {
		fmt.Printf("error:%s", err)
		return
	}
	fmt.Printf("connected\n")

	if true {
		resp, err := cli.TDXHandshake()
		if err != nil {
			fmt.Printf("error:%s", err)
			return
		}
		fmt.Printf("%s\n", resp)
	}

	if false {
		serverInfo, err := cli.ServerInfo()
		if err != nil {
			fmt.Printf("error:%s", err)
			return
		}

		fmt.Printf("%s\n", serverInfo.Resp.Name)
	}

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

func testOld() {
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
