package v2

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"gotdx/models"
	"gotdx/tdx"
)

type Client struct {
	ctx context.Context
	opt *tdx.Options

	dataConn *connRuntime

	metaConn *connRuntime
}

type reqPkg struct {
	body     []byte
	header   ReqHeader
	callback chan *respPkg
}
type respPkg struct {
	header RespHeader
	body   []byte
}

func NewClient(opt tdx.Option) *Client {
	cli := &Client{
		opt: tdx.ApplyOptions(opt),
	}
	cli.init()
	cli.Log("init success.")
	return cli
}

// private
func (c *Client) init() error {
	return nil
}

func (c *Client) Log(msg string, args ...any) {
	c.opt.MsgCallback(models.ProcessInfo{Msg: fmt.Sprintf(msg, args...)})
}

// use generic type
func do[T Codec](c *Client, conn *connRuntime, api T) error {
	if c == nil {
		return fmt.Errorf("client is nil")
	}
	if conn == nil {
		return fmt.Errorf("conn is nil or not connected")
	}
	c.Log("start do")
	var err error
	reqHeader := ReqHeader{
		Zip:        0x0C,
		SeqID:      conn.genSeqID(),
		PacketType: 0x01,
		PkgLen1:    0,
		PkgLen2:    0,
		Method:     0,
	}
	err = api.FillReqHeader(c.ctx, &reqHeader)
	if err != nil {
		return err
	}
	reqData, err := api.MarshalReqBody(c.ctx)
	if err != nil {
		return err
	}
	reqHeader.PkgLen1 = 2 + uint16(len(reqData))
	reqHeader.PkgLen2 = 2 + uint16(len(reqData))

	// send req
	reqBuf := bytes.NewBuffer(nil)
	err = binary.Write(reqBuf, binary.LittleEndian, reqHeader)
	if err != nil {
		return err
	}
	_, err = reqBuf.Write(reqData)
	if err != nil {
		return err
	}

	if c.opt.Debug && api.IsDebug(c.ctx) {
		c.Log("send %s: %s", reqHeader.Method, hex.EncodeToString(reqBuf.Bytes()))
	}

	callback := make(chan *respPkg)
	conn.sendCh <- &reqPkg{header: reqHeader, body: reqBuf.Bytes(), callback: callback}
	// wait resp
	respPkg := <-callback

	if c.opt.Debug && api.IsDebug(c.ctx) {
		c.Log("recv %s: %s", respPkg.header.Method, hex.EncodeToString(respPkg.body))
	}

	err = api.UnmarshalResp(c.ctx, respPkg.body)
	if err != nil {
		return err
	}
	return nil
}

// public
func (c *Client) Connect() error {
	if c.dataConn != nil && c.dataConn.isConnected() {
		return nil
	}
	if c.dataConn == nil {
		c.dataConn = newConnRuntime(c.ctx, connRuntimeOpt{
			heartbeatInterval: c.opt.HeartbeatInterval,
			log:               c.Log,
			heartbeatFunc: func() error {
				return c.Heartbeat()
			},
		})
	}
	return c.dataConn.connect(c.opt.TCPAddress)
}

func (c *Client) Disconnect() error {
	panic("implement me!")
}

func (c *Client) Handshake() error {
	handshake := &Handshake{}
	err := do(c, c.dataConn, handshake)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Heartbeat() error {
	c.Log("call heartbeat")
	heartbeat := &Heartbeat{}
	err := do(c, c.dataConn, heartbeat)
	if err != nil {
		return err
	}
	return nil
}
