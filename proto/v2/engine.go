package v2

import (
	"bytes"
	"compress/zlib"
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"gotdx/models"
	"gotdx/tdx"
	"io"
	"math/rand/v2"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type Client struct {
	ctx context.Context
	opt *tdx.Options

	done            chan struct{}
	heartbeatTicker *time.Ticker

	seqID uint32

	muConn sync.Mutex
	// mu start
	conn      net.Conn
	connected bool
	// mu end
	sendCh            chan *reqPkg
	muHandlerRedister sync.Mutex
	handlerRegister   map[uint32]*reqPkg
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
		opt:             tdx.ApplyOptions(opt),
		done:            make(chan struct{}),
		sendCh:          make(chan *reqPkg),
		handlerRegister: make(map[uint32]*reqPkg),
	}
	cli.heartbeatTicker = time.NewTicker(cli.opt.HeartbeatInterval)
	cli.init()
	cli.opt.MsgCallback(models.ProcessInfo{Msg: "init success."})
	return cli
}

// private
func (c *Client) init() error {
	c.seqID = rand.Uint32()
	// heartbeat ticker
	go func() {
		// TODO recover panic
		for {
			select {
			case <-c.heartbeatTicker.C:
				if !c.connected {
					continue
				}
				t0 := time.Now()
				err := c.heartbeat()
				if err != nil {
					c.muConn.Lock()
					c.connected = false
					c.muConn.Unlock()
					c.opt.MsgCallback(models.ProcessInfo{
						Msg: fmt.Sprintf("[%s] Detected connection broken, cost %d ms.", t0.Format("15:04:05"), time.Since(t0).Milliseconds()),
					})
				} else {
					c.opt.MsgCallback(models.ProcessInfo{
						Msg: fmt.Sprintf("[%s] 心跳成功, cost %d ms.", t0.Format("15:04:05"), time.Since(t0).Milliseconds()),
					})
				}
			case <-c.done:
				return
			}
		}
	}()
	return nil
}

func (c *Client) genSeqID() uint32 {
	return atomic.AddUint32(&c.seqID, 1)
}

func (c *Client) Log(msg string, args ...any) {
	c.opt.MsgCallback(models.ProcessInfo{Msg: fmt.Sprintf(msg, args...)})
}

// use generic type
func do[T Codec](c *Client, api T) error {
	if c == nil {
		return fmt.Errorf("client is nil")
	}
	c.opt.MsgCallback(models.ProcessInfo{Msg: "do start."})
	var err error
	reqHeader := ReqHeader{
		Zip:        0x0C,
		SeqID:      c.genSeqID(),
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
	c.sendCh <- &reqPkg{header: reqHeader, body: reqBuf.Bytes(), callback: callback}
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

func (c *Client) heartbeat() error {
	panic("implement me!")
}

// public
func (c *Client) Connect() error {
	c.muConn.Lock()
	defer c.muConn.Unlock()

	conn, err := net.Dial("tcp", c.opt.TCPAddress)
	if err != nil {
		return err
	}
	c.conn = conn

	go func() {
		for {
			c.opt.MsgCallback(models.ProcessInfo{Msg: "send start."})
			d := <-c.sendCh
			if d == nil {
				continue
			}
			c.Log("send: %s %d", d.header.Method, d.header.SeqID)
			c.muHandlerRedister.Lock()
			c.handlerRegister[d.header.SeqID] = d
			c.muHandlerRedister.Unlock()
			n, err := c.conn.Write(d.body)
			if err != nil {
				c.Log("write fail: %s", err.Error())
				break
			}
			c.Log("send: %d, size: %d", d.header.SeqID, n)
		}
		c.resetConn()
	}()

	go func() {
		for {
			var err error
			c.Log("read start.")

			// read header
			headerBuf := make([]byte, 16)
			_, err = io.ReadFull(c.conn, headerBuf)
			if err != nil {
				c.Log("read header fail: %v", err)
				break
			}
			var header RespHeader
			if err := binary.Read(bytes.NewBuffer(headerBuf), binary.LittleEndian, &header); err != nil {
				c.Log("parse header fail: %v", err)
				break
			}

			c.Log("read: %d, size: %d", header.SeqID, header.PkgDataSize)

			body := make([]byte, header.PkgDataSize)
			_, err = io.ReadFull(c.conn, body)
			if err != nil {
				c.Log("read body fail: %v", err)
				break
			}
			if header.PkgDataSize != header.RawDataSize {
				r, _ := zlib.NewReader(bytes.NewReader(body))
				body, err = io.ReadAll(r)
				if err != nil {
					c.Log("unzip fail: %v", err)
					continue
				}
			}
			// dispatch
			c.muHandlerRedister.Lock()
			reqPkg, ok := c.handlerRegister[header.SeqID]
			if !ok {
				c.muHandlerRedister.Unlock()
				c.Log("handler not found: %s %d", header.Method, header.SeqID)
				continue
			}
			delete(c.handlerRegister, header.SeqID)
			c.muHandlerRedister.Unlock()
			if reqPkg == nil || reqPkg.callback == nil {
				c.Log("callback is nil: %s %d", header.Method, header.SeqID)
				continue
			}
			c.Log("call callback: %s %d", header.Method, header.SeqID)
			reqPkg.callback <- &respPkg{header: header, body: body}
		}
		c.resetConn()
	}()
	err = c.Handshake()
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) resetConn() {
	c.muConn.Lock()
	defer c.muConn.Unlock()
	if c.conn != nil {
		_ = c.conn.Close()
	}
	c.connected = false
}

func (c *Client) Disconnect() error {
	panic("implement me!")
}

func (c *Client) Handshake() error {
	handshake := &Handshake{}
	handshake.SetDebug(c.ctx)
	err := do(c, handshake)
	if err != nil {
		return err
	}
	c.opt.MsgCallback(models.ProcessInfo{Msg: handshake.ContentHex})
	return nil
}

func (c *Client) Heartbeat() error {
	panic("implement me!")
}
